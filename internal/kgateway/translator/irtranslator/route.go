package irtranslator

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"slices"

	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_route_v3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	envoy_type_matcher_v3 "github.com/envoyproxy/go-control-plane/envoy/type/matcher/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	gwv1 "sigs.k8s.io/gateway-api/apis/v1"

	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/translator/routeutils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/ir"
	"github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/policy"
	reportssdk "github.com/kgateway-dev/kgateway/v2/pkg/pluginsdk/reporter"
	"github.com/kgateway-dev/kgateway/v2/pkg/utils/regexutils"
)

const (
	webSocketUpgradeType = "websocket"
	// TODO(tim): determine if this is the correct reason and move to the api package.
	routeTranslationFailedReason = "RouteTranslationFailed"
)

type httpRouteConfigurationTranslator struct {
	gw               ir.GatewayIR
	listener         ir.ListenerIR
	fc               ir.FilterChainCommon
	attachedPolicies ir.AttachedPolicies

	routeConfigName          string
	reporter                 reportssdk.Reporter
	requireTlsOnVirtualHosts bool
	enableRouteReplacement   bool
	PluginPass               TranslationPassPlugins
	logger                   *slog.Logger
}

func (h *httpRouteConfigurationTranslator) ComputeRouteConfiguration(ctx context.Context, vhosts []*ir.VirtualHost) *envoy_config_route_v3.RouteConfiguration {
	var attachedPolicies ir.AttachedPolicies
	// the policies in order - first listener as they are more specific and thus higher priority.
	// then gateway policies.
	attachedPolicies.Append(h.attachedPolicies, h.gw.AttachedHttpPolicies)
	cfg := &envoy_config_route_v3.RouteConfiguration{
		Name: h.routeConfigName,
	}
	typedPerFilterConfigRoute := ir.TypedFilterConfigMap(map[string]proto.Message{})

	for _, gk := range attachedPolicies.ApplyOrderedGroupKinds() {
		pols := attachedPolicies.Policies[gk]
		pass := h.PluginPass[gk]
		if pass == nil {
			// TODO: user error - they attached a non http policy
			continue
		}
		reportPolicyAcceptanceStatus(h.reporter, h.listener.PolicyAncestorRef, pols...)
		for _, pol := range mergePolicies(pass, pols) {
			pass.ApplyRouteConfigPlugin(ctx, &ir.RouteConfigContext{
				FilterChainName:   h.fc.FilterChainName,
				TypedFilterConfig: typedPerFilterConfigRoute,
				Policy:            pol.PolicyIr,
			}, cfg)
		}
	}

	cfg.VirtualHosts = h.computeVirtualHosts(ctx, vhosts)
	cfg.TypedPerFilterConfig = typedPerFilterConfigRoute.ToAnyMap()

	// Gateway API spec requires that port values in HTTP Host headers be ignored when performing a match
	// See https://gateway-api.sigs.k8s.io/reference/spec/#gateway.networking.k8s.io/v1.HTTPRouteSpec - hostnames field
	cfg.IgnorePortInHostMatching = true

	//	if mostSpecificVal := h.parentListener.GetRouteOptions().GetMostSpecificHeaderMutationsWins(); mostSpecificVal != nil {
	//		cfg.MostSpecificHeaderMutationsWins = mostSpecificVal.GetValue()
	//	}

	return cfg
}

func (h *httpRouteConfigurationTranslator) computeVirtualHosts(ctx context.Context, virtualHosts []*ir.VirtualHost) []*envoy_config_route_v3.VirtualHost {
	var envoyVirtualHosts []*envoy_config_route_v3.VirtualHost
	for _, virtualHost := range virtualHosts {
		envoyVirtualHosts = append(envoyVirtualHosts, h.computeVirtualHost(ctx, virtualHost))
	}
	return envoyVirtualHosts
}

func (h *httpRouteConfigurationTranslator) computeVirtualHost(
	ctx context.Context,
	virtualHost *ir.VirtualHost,
) *envoy_config_route_v3.VirtualHost {
	sanitizedName := utils.SanitizeForEnvoy(ctx, virtualHost.Name, "virtual host")

	var (
		envoyRoutes []*envoy_config_route_v3.Route
		errs        []error
	)
	for i, route := range virtualHost.Rules {
		// TODO: not sure if we need listener parent ref here or the http parent ref
		var routeReport reportssdk.ParentRefReporter = &reports.ParentRefReport{}
		if route.Parent != nil {
			// route may be a fake one that we don't really report,
			// such as in the waypoint translator where we produce
			// synthetic routes if there none are attached to the Gateway/Service.
			routeReport = h.reporter.Route(route.Parent.SourceObject).ParentRef(&route.ParentRef)
		}
		generatedName := fmt.Sprintf("%s-route-%d", virtualHost.Name, i)
		computedRoute, err := h.envoyRoutes(ctx, routeReport, route, generatedName)
		if err != nil {
			errs = append(errs, err)
		}
		if computedRoute != nil {
			envoyRoutes = append(envoyRoutes, computedRoute)
		}
	}
	domains := []string{virtualHost.Hostname}
	if len(domains) == 0 || (len(domains) == 1 && domains[0] == "") {
		domains = []string{"*"}
	}
	var envoyRequireTls envoy_config_route_v3.VirtualHost_TlsRequirementType
	if h.requireTlsOnVirtualHosts {
		// TODO (ilackarms): support external-only TLS
		envoyRequireTls = envoy_config_route_v3.VirtualHost_ALL
	}
	if len(errs) == 0 {
		// Implement this correctly or remove this logic.
		logger.Info("route validation succeeded", "route", virtualHost.Name)
		// clear stale partially invalid condition. not sure how to implement this correctly.
		// h.reporter.Route(virtualHost.Name).SetCondition(reportssdk.RouteCondition{
		// 	Type:    gwv1.RouteConditionPartiallyInvalid,
		// 	Status:  metav1.ConditionFalse,
		// 	Reason:  gwv1.RouteReasonAccepted,
		// 	Message: "Route validation succeeded",
		// })
	}

	out := &envoy_config_route_v3.VirtualHost{
		Name:       sanitizedName,
		Domains:    domains,
		Routes:     envoyRoutes,
		RequireTls: envoyRequireTls,
	}

	typedPerFilterConfigRoute := ir.TypedFilterConfigMap(map[string]proto.Message{})
	// run the http plugins that are attached to the listener or gateway on the virtual host
	h.runVhostPlugins(ctx, virtualHost, out, typedPerFilterConfigRoute)
	out.TypedPerFilterConfig = typedPerFilterConfigRoute.ToAnyMap()

	return out
}

func (h *httpRouteConfigurationTranslator) envoyRoutes(
	ctx context.Context,
	routeReport reportssdk.ParentRefReporter,
	in ir.HttpRouteRuleMatchIR,
	generatedName string,
) (*envoy_config_route_v3.Route, error) {
	out := h.initRoutes(in, generatedName)

	typedPerFilterConfigRoute := ir.TypedFilterConfigMap(map[string]proto.Message{})
	// configure the route action
	if err := h.configureRouteAction(ctx, in, out, typedPerFilterConfigRoute); err != nil {
		return h.handleRouteError(err, in, out, routeReport)
	}
	// run route plugins that may set action or typed per filter config
	if err := h.runRoutePlugins(ctx, routeReport, in, out, typedPerFilterConfigRoute); err != nil {
		return h.handleRouteError(err, in, out, routeReport)
	}
	// validate envoy route for structural errors
	if err := validateEnvoyRoute(out); err != nil {
		return h.handleRouteError(err, in, out, routeReport)
	}
	// apply typed per filter config from translating route action and route plugins
	if err := h.applyTypedFilterConfig(out, typedPerFilterConfigRoute); err != nil {
		return h.handleRouteError(err, in, out, routeReport)
	}
	// check if action is required but missing. skip for delegates.
	if out.GetAction() == nil {
		if in.Delegates {
			return nil, nil
		}
		return h.handleRouteError(errors.New("no action specified"), in, out, routeReport)
	}

	return out, nil
}

func (h *httpRouteConfigurationTranslator) configureRouteAction(
	ctx context.Context,
	in ir.HttpRouteRuleMatchIR,
	out *envoy_config_route_v3.Route,
	typedPerFilterConfigRoute ir.TypedFilterConfigMap,
) error {
	if len(in.Backends) == 1 {
		// if there's only one backend, we need to reuse typedPerFilterConfigRoute in both translateRouteAction and runRoutePlugins
		out.Action = h.translateRouteAction(ctx, in, out, typedPerFilterConfigRoute)
		return nil
	}
	if len(in.Backends) > 0 {
		// If there is more than one backend, we translate the backends as WeightedClusters and each weighted cluster
		// will have a TypedPerFilterConfig that overrides the parent route-level config.
		out.Action = h.translateRouteAction(ctx, in, out, nil)
		return nil
	}
	return nil
}

func (h *httpRouteConfigurationTranslator) applyTypedFilterConfig(
	out *envoy_config_route_v3.Route,
	typedPerFilterConfigRoute ir.TypedFilterConfigMap,
) error {
	typedPerFilterConfig := typedPerFilterConfigRoute.ToAnyMap()
	if out.GetTypedPerFilterConfig() == nil {
		out.TypedPerFilterConfig = typedPerFilterConfig
		return nil
	}
	for k, v := range typedPerFilterConfig {
		if _, exists := out.GetTypedPerFilterConfig()[k]; !exists {
			out.GetTypedPerFilterConfig()[k] = v
		}
	}
	return nil
}

func (h *httpRouteConfigurationTranslator) handleRouteError(
	err error,
	in ir.HttpRouteRuleMatchIR,
	out *envoy_config_route_v3.Route,
	routeReport reportssdk.ParentRefReporter,
) (*envoy_config_route_v3.Route, error) {
	if h.enableRouteReplacement {
		return h.replaceInvalidRoute(err, in, out, routeReport), err
	}
	h.logger.Debug("invalid route", "error", err)
	routeReport.SetCondition(reportssdk.RouteCondition{
		Type:    gwv1.RouteConditionPartiallyInvalid,
		Status:  metav1.ConditionTrue,
		Reason:  routeTranslationFailedReason,
		Message: fmt.Sprintf("Dropped Rule %d: %v", in.MatchIndex, err),
	})
	return nil, err
}

func (h *httpRouteConfigurationTranslator) replaceInvalidRoute(
	err error,
	in ir.HttpRouteRuleMatchIR,
	out *envoy_config_route_v3.Route,
	routeReport reportssdk.ParentRefReporter,
) *envoy_config_route_v3.Route {
	h.logger.Error("invalid route", "error", err)
	routeReport.SetCondition(reportssdk.RouteCondition{
		Type:    gwv1.RouteConditionPartiallyInvalid,
		Status:  metav1.ConditionTrue,
		Reason:  routeTranslationFailedReason,
		Message: fmt.Sprintf("Dropped Rule %d: %v", in.MatchIndex, err),
	})
	out.Action = &envoy_config_route_v3.Route_DirectResponse{
		DirectResponse: &envoy_config_route_v3.DirectResponseAction{
			Status: http.StatusInternalServerError,
			Body: &corev3.DataSource{
				Specifier: &corev3.DataSource_InlineString{
					InlineString: `Kgateway has detected an invalid route configuration`,
				},
			},
		},
	}
	return out
}

func (h *httpRouteConfigurationTranslator) runVhostPlugins(ctx context.Context, virtualHost *ir.VirtualHost, out *envoy_config_route_v3.VirtualHost,
	typedPerFilterConfig ir.TypedFilterConfigMap,
) {
	for _, gk := range virtualHost.AttachedPolicies.ApplyOrderedGroupKinds() {
		pols := virtualHost.AttachedPolicies.Policies[gk]
		pass := h.PluginPass[gk]
		if pass == nil {
			// TODO: user error - they attached a non http policy
			continue
		}
		reportPolicyAcceptanceStatus(h.reporter, h.listener.PolicyAncestorRef, pols...)
		for _, pol := range mergePolicies(pass, pols) {
			pctx := &ir.VirtualHostContext{
				Policy:            pol.PolicyIr,
				TypedFilterConfig: typedPerFilterConfig,
				FilterChainName:   h.fc.FilterChainName,
			}
			pass.ApplyVhostPlugin(ctx, pctx, out)
			// TODO: check return value, if error returned, log error and report condition
		}
	}
}

func (h *httpRouteConfigurationTranslator) runRoutePlugins(
	ctx context.Context,
	routeReport reportssdk.ParentRefReporter,
	in ir.HttpRouteRuleMatchIR,
	out *envoy_config_route_v3.Route,
	typedPerFilterConfig ir.TypedFilterConfigMap,
) error {
	// all policies up to listener have been applied as vhost polices; we need to apply the httproute policies and below
	//
	// NOTE: AttachedPolicies must have policies in the ordered by hierarchy from root to leaf in the delegation chain where
	// each level has policies ordered by rule level policies before entire route level policies.

	var attachedPolicies ir.AttachedPolicies
	delegatingParent := in.DelegatingParent
	var hierarchicalPriority int
	for delegatingParent != nil {
		hierarchicalPriority++
		attachedPolicies.Prepend(hierarchicalPriority,
			delegatingParent.ExtensionRefs, delegatingParent.AttachedPolicies, delegatingParent.Parent.AttachedPolicies)
		delegatingParent = delegatingParent.DelegatingParent
	}

	// rule-level policies in priority order (high to low)
	attachedPolicies.Append(in.ExtensionRefs, in.AttachedPolicies)

	// route-level policy
	if in.Parent != nil {
		attachedPolicies.Append(in.Parent.AttachedPolicies)
	}

	for _, gk := range attachedPolicies.ApplyOrderedGroupKinds() {
		pols := attachedPolicies.Policies[gk]
		pass := h.PluginPass[gk]
		if pass == nil {
			logger.Debug("no pass found for policy", "policy", gk)
			continue
		}
		reportPolicyAcceptanceStatus(h.reporter, h.listener.PolicyAncestorRef, pols...)
		pctx := &ir.RouteContext{
			FilterChainName:   h.fc.FilterChainName,
			In:                in,
			TypedFilterConfig: typedPerFilterConfig,
		}
		for _, pol := range mergePolicies(pass, pols) {
			if hasTerminalError := h.processPolicyErrors(pol.Errors); hasTerminalError {
				return policy.NewTerminalError("policy reported terminal errors", errors.Join(pol.Errors...))
			}
			pctx.Policy = pol.PolicyIr
			pass.ApplyForRoute(ctx, pctx, out)
		}
		// TODO: check return value, if error returned, log error and report condition
	}
	return nil
}

// processPolicyErrors checks if there are any terminal errors in the policy errors.
// Returns true if there are terminal errors, false otherwise.
func (h *httpRouteConfigurationTranslator) processPolicyErrors(policyErrors []error) bool {
	if len(policyErrors) == 0 {
		return false
	}
	for _, err := range policyErrors {
		var policyErr *policy.PolicyError
		if !errors.As(err, &policyErr) {
			continue
		}
		if !policyErr.IsTerminal() {
			continue
		}
		return true
	}
	return false
}

func mergePolicies(pass *TranslationPass, policies []ir.PolicyAtt) []ir.PolicyAtt {
	if pass.MergePolicies != nil {
		merged := [1]ir.PolicyAtt{pass.MergePolicies(policies)}
		return merged[:]
	}

	return policies
}

func (h *httpRouteConfigurationTranslator) runBackendPolicies(ctx context.Context, in ir.HttpBackend, pCtx *ir.RouteBackendContext) error {
	for _, gk := range in.AttachedPolicies.ApplyOrderedGroupKinds() {
		pols := in.AttachedPolicies.Policies[gk]
		pass := h.PluginPass[gk]
		if pass == nil {
			// TODO: should never happen, log error and report condition
			continue
		}
		reportPolicyAcceptanceStatus(h.reporter, h.listener.PolicyAncestorRef, pols...)
		for _, pol := range mergePolicies(pass, pols) {
			// Policy on extension ref
			pass.ApplyForRouteBackend(ctx, pol.PolicyIr, pCtx)
			// TODO: check return value, if error returned, log error and report condition
		}
	}
	return nil
}

func (h *httpRouteConfigurationTranslator) runBackend(ctx context.Context, in ir.HttpBackend, pCtx *ir.RouteBackendContext, outRoute *envoy_config_route_v3.Route) error {
	if in.Backend.BackendObject != nil {
		backendPass := h.PluginPass[in.Backend.BackendObject.GetGroupKind()]
		if backendPass != nil {
			backendPass.ApplyForBackend(ctx, pCtx, in, outRoute)
		}
	}
	return nil
}

func (h *httpRouteConfigurationTranslator) translateRouteAction(
	ctx context.Context,
	in ir.HttpRouteRuleMatchIR,
	outRoute *envoy_config_route_v3.Route,
	parentTypedPerFilterConfig ir.TypedFilterConfigMap,
) *envoy_config_route_v3.Route_Route {
	var clusters []*envoy_config_route_v3.WeightedCluster_ClusterWeight

	for _, backend := range in.Backends {
		clusterName := backend.Backend.ClusterName

		// get backend for ref - we must do it to make sure we have permissions to access it.
		// also we need the service so we can translate its name correctly.
		cw := &envoy_config_route_v3.WeightedCluster_ClusterWeight{
			Name:   clusterName,
			Weight: wrapperspb.UInt32(backend.Backend.Weight),
		}

		typedPerFilterConfig := parentTypedPerFilterConfig
		if parentTypedPerFilterConfig == nil {
			typedPerFilterConfig = map[string]proto.Message{}
		}

		pCtx := ir.RouteBackendContext{
			FilterChainName:   h.fc.FilterChainName,
			Backend:           backend.Backend.BackendObject,
			TypedFilterConfig: typedPerFilterConfig,
		}

		// non attached policy translation
		err := h.runBackend(
			ctx,
			backend,
			&pCtx,
			outRoute,
		)
		if err != nil {
			// TODO: error on status
			h.logger.Error("error processing backends", "error", err)
		}

		err = h.runBackendPolicies(
			ctx,
			backend,
			&pCtx,
		)
		if err != nil {
			// TODO: error on status
			h.logger.Error("error processing backends with policies", "error", err)
		}

		// Translating weighted clusters needs the typed per filter config on each cluster
		cw.TypedPerFilterConfig = typedPerFilterConfig.ToAnyMap()
		clusters = append(clusters, cw)
	}

	// TODO: i think envoy nacks if all weights are 0, we should error on that.
	action := outRoute.GetRoute()
	if action == nil {
		action = &envoy_config_route_v3.RouteAction{
			ClusterNotFoundResponseCode: envoy_config_route_v3.RouteAction_INTERNAL_SERVER_ERROR,
		}
	}

	routeAction := &envoy_config_route_v3.Route_Route{
		Route: action,
	}
	switch len(clusters) {
	// case 0:
	// TODO: we should never get here
	case 1:
		// Only set the cluster name if unspecified since a plugin may have set it.
		if action.GetCluster() == "" {
			action.ClusterSpecifier = &envoy_config_route_v3.RouteAction_Cluster{
				Cluster: clusters[0].GetName(),
			}
		}
		// Skip setting the typed per filter config here, set it in the envoyRoutes() after runRoutePlugins runs

	default:
		// Only set weighted clusters if unspecified since a plugin may have set it.
		if action.GetWeightedClusters() == nil {
			action.ClusterSpecifier = &envoy_config_route_v3.RouteAction_WeightedClusters{
				WeightedClusters: &envoy_config_route_v3.WeightedCluster{
					Clusters: clusters,
				},
			}
		}
	}

	for _, backend := range in.Backends {
		if back := backend.Backend.BackendObject; back != nil && back.AppProtocol == ir.WebSocketAppProtocol {
			// add websocket upgrade if not already present
			if !slices.ContainsFunc(action.GetUpgradeConfigs(), func(uc *envoy_config_route_v3.RouteAction_UpgradeConfig) bool {
				return uc.GetUpgradeType() == webSocketUpgradeType
			}) {
				action.UpgradeConfigs = append(action.GetUpgradeConfigs(), &envoy_config_route_v3.RouteAction_UpgradeConfig{
					UpgradeType: webSocketUpgradeType,
				})
			}
		}
	}
	return routeAction
}

func validateEnvoyRoute(r *envoy_config_route_v3.Route) error {
	var errs []error
	match := r.GetMatch()
	route := r.GetRoute()
	re := r.GetRedirect()
	validatePath(match.GetPath(), &errs)
	validatePath(match.GetPrefix(), &errs)
	validatePath(match.GetPathSeparatedPrefix(), &errs)
	validatePath(re.GetPathRedirect(), &errs)
	validatePath(re.GetHostRedirect(), &errs)
	validatePath(re.GetSchemeRedirect(), &errs)
	validatePrefixRewrite(route.GetPrefixRewrite(), &errs)
	validatePrefixRewrite(re.GetPrefixRewrite(), &errs)
	if len(errs) == 0 {
		return nil
	}
	return fmt.Errorf("error %s: %w", r.GetName(), errors.Join(errs...))
}

// creates Envoy routes for each matcher provided on our Gateway route
func (h *httpRouteConfigurationTranslator) initRoutes(
	in ir.HttpRouteRuleMatchIR,
	generatedName string,
) *envoy_config_route_v3.Route {
	//	if len(in.Matches) == 0 {
	//		return []*envoy_config_route_v3.Route{
	//			{
	//				Match: &envoy_config_route_v3.RouteMatch{
	//					PathSpecifier: &envoy_config_route_v3.RouteMatch_Prefix{Prefix: "/"},
	//				},
	//			},
	//		}
	//	}

	out := &envoy_config_route_v3.Route{
		Match: translateGlooMatcher(in.Match),
	}
	name := in.Name
	if name != "" {
		out.Name = fmt.Sprintf("%s-%s-matcher-%d", generatedName, name, in.MatchIndex)
	} else {
		out.Name = fmt.Sprintf("%s-matcher-%d", generatedName, in.MatchIndex)
	}

	return out
}

func translateGlooMatcher(matcher gwv1.HTTPRouteMatch) *envoy_config_route_v3.RouteMatch {
	match := &envoy_config_route_v3.RouteMatch{
		Headers:         envoyHeaderMatcher(matcher.Headers),
		QueryParameters: envoyQueryMatcher(matcher.QueryParams),
	}
	if matcher.Method != nil {
		match.Headers = append(match.GetHeaders(), &envoy_config_route_v3.HeaderMatcher{
			Name: ":method",
			HeaderMatchSpecifier: &envoy_config_route_v3.HeaderMatcher_StringMatch{
				StringMatch: &envoy_type_matcher_v3.StringMatcher{
					MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
						Exact: string(*matcher.Method),
					},
				},
			},
		})
	}

	setEnvoyPathMatcher(matcher, match)
	return match
}

var separatedPathRegex = regexp.MustCompile("^[^?#]+[^?#/]$")

func isValidPathSparated(path string) bool {
	// see envoy docs:
	//	Expect the value to not contain "?" or "#" and not to end in "/"
	return separatedPathRegex.MatchString(path)
}

func setEnvoyPathMatcher(match gwv1.HTTPRouteMatch, out *envoy_config_route_v3.RouteMatch) {
	pathType, pathValue := routeutils.ParsePath(match.Path)
	switch pathType {
	case gwv1.PathMatchPathPrefix:
		if !isValidPathSparated(pathValue) {
			out.PathSpecifier = &envoy_config_route_v3.RouteMatch_Prefix{
				Prefix: pathValue,
			}
		} else {
			out.PathSpecifier = &envoy_config_route_v3.RouteMatch_PathSeparatedPrefix{
				PathSeparatedPrefix: pathValue,
			}
		}
	case gwv1.PathMatchExact:
		out.PathSpecifier = &envoy_config_route_v3.RouteMatch_Path{
			Path: pathValue,
		}
	case gwv1.PathMatchRegularExpression:
		out.PathSpecifier = &envoy_config_route_v3.RouteMatch_SafeRegex{
			SafeRegex: regexutils.NewRegexWithProgramSize(pathValue, nil),
		}
	}
}

func envoyHeaderMatcher(in []gwv1.HTTPHeaderMatch) []*envoy_config_route_v3.HeaderMatcher {
	var out []*envoy_config_route_v3.HeaderMatcher
	for _, matcher := range in {
		envoyMatch := &envoy_config_route_v3.HeaderMatcher{
			Name: string(matcher.Name),
		}
		regex := false
		if matcher.Type != nil && *matcher.Type == gwv1.HeaderMatchRegularExpression {
			regex = true
		}

		// TODO: not sure if we should do PresentMatch according to the spec.
		if matcher.Value == "" {
			envoyMatch.HeaderMatchSpecifier = &envoy_config_route_v3.HeaderMatcher_PresentMatch{
				PresentMatch: true,
			}
		} else {
			if regex {
				envoyMatch.HeaderMatchSpecifier = &envoy_config_route_v3.HeaderMatcher_StringMatch{
					StringMatch: &envoy_type_matcher_v3.StringMatcher{
						MatchPattern: &envoy_type_matcher_v3.StringMatcher_SafeRegex{
							SafeRegex: regexutils.NewRegexWithProgramSize(matcher.Value, nil),
						},
					},
				}
			} else {
				envoyMatch.HeaderMatchSpecifier = &envoy_config_route_v3.HeaderMatcher_StringMatch{
					StringMatch: &envoy_type_matcher_v3.StringMatcher{
						MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
							Exact: matcher.Value,
						},
					},
				}
			}
		}
		out = append(out, envoyMatch)
	}
	return out
}

func envoyQueryMatcher(in []gwv1.HTTPQueryParamMatch) []*envoy_config_route_v3.QueryParameterMatcher {
	var out []*envoy_config_route_v3.QueryParameterMatcher
	for _, matcher := range in {
		envoyMatch := &envoy_config_route_v3.QueryParameterMatcher{
			Name: string(matcher.Name),
		}
		regex := false
		if matcher.Type != nil && *matcher.Type == gwv1.QueryParamMatchRegularExpression {
			regex = true
		}

		// TODO: not sure if we should do PresentMatch according to the spec.
		if matcher.Value == "" {
			envoyMatch.QueryParameterMatchSpecifier = &envoy_config_route_v3.QueryParameterMatcher_PresentMatch{
				PresentMatch: true,
			}
		} else {
			if regex {
				envoyMatch.QueryParameterMatchSpecifier = &envoy_config_route_v3.QueryParameterMatcher_StringMatch{
					StringMatch: &envoy_type_matcher_v3.StringMatcher{
						MatchPattern: &envoy_type_matcher_v3.StringMatcher_SafeRegex{
							SafeRegex: regexutils.NewRegexWithProgramSize(matcher.Value, nil),
						},
					},
				}
			} else {
				envoyMatch.QueryParameterMatchSpecifier = &envoy_config_route_v3.QueryParameterMatcher_StringMatch{
					StringMatch: &envoy_type_matcher_v3.StringMatcher{
						MatchPattern: &envoy_type_matcher_v3.StringMatcher_Exact{
							Exact: matcher.Value,
						},
					},
				}
			}
		}
		out = append(out, envoyMatch)
	}
	return out
}
