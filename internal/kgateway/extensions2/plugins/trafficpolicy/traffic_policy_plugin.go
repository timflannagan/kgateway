package trafficpolicy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	envoy_core_v3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	envoy_config_listener_v3 "github.com/envoyproxy/go-control-plane/envoy/config/listener/v3"
	routev3 "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	exteniondynamicmodulev3 "github.com/envoyproxy/go-control-plane/envoy/extensions/dynamic_modules/v3"
	dynamicmodulesv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/dynamic_modules/v3"
	envoy_ext_authz_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_authz/v3"
	envoy_ext_proc_v3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ext_proc/v3"
	localratelimitv3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/local_ratelimit/v3"
	ratev3 "github.com/envoyproxy/go-control-plane/envoy/extensions/filters/http/ratelimit/v3"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/wrapperspb"
	skubeclient "istio.io/istio/pkg/config/schema/kubeclient"
	"istio.io/istio/pkg/kube/kclient"
	"istio.io/istio/pkg/kube/krt"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"

	// TODO(nfuden): remove once rustformations are able to be used in a production environment
	transformationpb "github.com/solo-io/envoy-gloo/go/config/filter/http/transformation/v2"

	apiannotations "github.com/kgateway-dev/kgateway/v2/api/annotations"
	"github.com/kgateway-dev/kgateway/v2/api/v1alpha1"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/common"
	extensionsplug "github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/plugin"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/extensions2/pluginutils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/ir"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/krtcollections"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/plugins"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/policy"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/reports"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/utils"
	"github.com/kgateway-dev/kgateway/v2/internal/kgateway/wellknown"
	"github.com/kgateway-dev/kgateway/v2/pkg/client/clientset/versioned"
	"github.com/kgateway-dev/kgateway/v2/pkg/logging"
	"github.com/kgateway-dev/kgateway/v2/pkg/validator"
)

const (
	transformationFilterNamePrefix              = "transformation"
	extAuthGlobalDisableFilterName              = "global_disable/ext_auth"
	extAuthGlobalDisableFilterMetadataNamespace = "dev.kgateway.disable_ext_auth"
	extAuthGlobalDisableKey                     = "extauth_disable"
	rustformationFilterNamePrefix               = "dynamic_modules/simple_mutations"
	metadataRouteTransformation                 = "transformation/helper"
	extauthFilterNamePrefix                     = "ext_auth"
	localRateLimitFilterNamePrefix              = "ratelimit/local"
	localRateLimitStatPrefix                    = "http_local_rate_limiter"
	rateLimitFilterNamePrefix                   = "ratelimit"
)

var logger = logging.New("plugin/trafficpolicy")

type trafficPolicySpecIr struct {
	transform                  *transformationpb.RouteTransformations
	rustformation              proto.Message
	rustformationStringToStash string

	ai              *aiPolicyIR
	extProc         *extProcIR
	extAuth         *extAuthIR
	localRateLimit  *localRateLimitIR
	globalRateLimit *globalRateLimitIR
}

type TrafficPolicy struct {
	ct   time.Time
	spec trafficPolicySpecIr
}

func (p *TrafficPolicy) Name() string { return "trafficpolicies" }

func (d *TrafficPolicy) CreationTime() time.Time { return d.ct }

func (d *TrafficPolicy) Equals(in any) bool {
	d2, ok := in.(*TrafficPolicy)
	if !ok {
		return false
	}
	if d.ct != d2.ct {
		return false
	}

	if !proto.Equal(d.spec.transform, d2.spec.transform) {
		return false
	}
	if !proto.Equal(d.spec.rustformation, d2.spec.rustformation) {
		return false
	}

	// delegate to each sub-policy's equality checks
	if !d.spec.ai.Equals(d2.spec.ai) {
		return false
	}
	if !d.spec.extAuth.Equals(d2.spec.extAuth) {
		return false
	}
	if !d.spec.extProc.Equals(d2.spec.extProc) {
		return false
	}
	if !d.spec.localRateLimit.Equals(d2.spec.localRateLimit) {
		return false
	}
	if !d.spec.globalRateLimit.Equals(d2.spec.globalRateLimit) {
		return false
	}
	return true
}

func (p *TrafficPolicy) Validate() error {
	var errs []error

	// TODO: Transformation.
	// TODO: Rustformation.
	// TODO: Is it safe to index into the spec like this?

	if err := p.spec.ai.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := p.spec.extAuth.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := p.spec.extProc.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := p.spec.localRateLimit.Validate(); err != nil {
		errs = append(errs, err)
	}
	if err := p.spec.globalRateLimit.Validate(); err != nil {
		errs = append(errs, err)
	}
	if len(errs) > 0 {
		return fmt.Errorf("traffic policy %s has the following errors: %v", p.Name(), errs)
	}
	return nil
}

type providerWithFromListener struct {
	provider     *TrafficPolicyGatewayExtensionIR
	fromListener bool
}

type trafficPolicyPluginGwPass struct {
	reporter reports.Reporter

	setTransformationInChain bool // TODO(nfuden): make this multi stage
	// TODO(nfuden): dont abuse httplevel filter in favor of route level
	rustformationStash map[string]string
	listenerTransform  *transformationpb.RouteTransformations
	ir.UnimplementedProxyTranslationPass
	localRateLimitInChain *localratelimitv3.LocalRateLimit
	extAuthPerProvider    map[string]providerWithFromListener
	extProcPerProvider    map[string]providerWithFromListener
	rateLimitPerProvider  map[string]providerWithFromListener
}

func (p *trafficPolicyPluginGwPass) SupportsPolicyMerge() bool { return true }

var useRustformations bool

func registerTypes(ourCli versioned.Interface) {
	skubeclient.Register[*v1alpha1.TrafficPolicy](
		wellknown.TrafficPolicyGVR,
		wellknown.TrafficPolicyGVK,
		func(c skubeclient.ClientGetter, namespace string, o metav1.ListOptions) (runtime.Object, error) {
			return ourCli.GatewayV1alpha1().TrafficPolicies(namespace).List(context.Background(), o)
		},
		func(c skubeclient.ClientGetter, namespace string, o metav1.ListOptions) (watch.Interface, error) {
			return ourCli.GatewayV1alpha1().TrafficPolicies(namespace).Watch(context.Background(), o)
		},
	)
}

func NewPlugin(ctx context.Context, commoncol *common.CommonCollections, validator validator.Validator) extensionsplug.Plugin {
	registerTypes(commoncol.OurClient)

	useRustformations = commoncol.Settings.UseRustFormations // stash the state of the env setup for rustformation usage

	col := krt.WrapClient(kclient.NewFiltered[*v1alpha1.TrafficPolicy](
		commoncol.Client,
		kclient.Filter{ObjectFilter: commoncol.Client.ObjectFilter()},
	), commoncol.KrtOpts.ToOptions("TrafficPolicy")...)
	gk := wellknown.TrafficPolicyGVK.GroupKind()

	translator := NewTrafficPolicyBuilder(ctx, commoncol)

	// TrafficPolicy IR will have TypedConfig -> implement backendroute method to add prompt guard, etc.
	policyCol := krt.NewCollection(col, func(krtctx krt.HandlerContext, policyCR *v1alpha1.TrafficPolicy) *ir.PolicyWrapper {
		// first translate the spec into an IR representation
		// and then validate the IR representation will produce
		// valid xDS configuration.
		//
		// TODO: How often will this be run? Does it run for each TP update?
		// Yeah... I'm guessing this will be run for each TP update, which
		// kinda invalidates the per-policy-Validate logic, right?
		policyIR, errors := translator.Translate(krtctx, policyCR)
		if err := validator.Validate(ctx, nil); err != nil {
			errors = append(errors, err)
		}
		pol := &ir.PolicyWrapper{
			ObjectSource: ir.ObjectSource{
				Group:     gk.Group,
				Kind:      gk.Kind,
				Namespace: policyCR.Namespace,
				Name:      policyCR.Name,
			},
			Policy:     policyCR,
			PolicyIR:   policyIR,
			TargetRefs: pluginutils.TargetRefsToPolicyRefs(policyCR.Spec.TargetRefs, policyCR.Spec.TargetSelectors),
			Errors:     errors,
		}
		return pol
	})

	return extensionsplug.Plugin{
		ContributesPolicies: map[schema.GroupKind]extensionsplug.PolicyPlugin{
			wellknown.TrafficPolicyGVK.GroupKind(): {
				NewGatewayTranslationPass: NewGatewayTranslationPass,
				Policies:                  policyCol,
				MergePolicies:             mergePolicies,
				GetPolicyStatus:           getPolicyStatusFn(commoncol.CrudClient),
				PatchPolicyStatus:         patchPolicyStatusFn(commoncol.CrudClient),
			},
		},
		ExtraHasSynced: translator.gatewayExtensions.HasSynced,
	}
}

func ResolveExtGrpcService(
	krtctx krt.HandlerContext,
	backends *krtcollections.BackendIndex,
	disableExtensionRefValidation bool,
	objectSource ir.ObjectSource,
	grpcService *v1alpha1.ExtGrpcService,
) (*envoy_core_v3.GrpcService, error) {
	var (
		clusterName string
		authority   string
	)
	if grpcService != nil {
		if grpcService.BackendRef == nil {
			return nil, errors.New("backend not provided")
		}
		backendRef := grpcService.BackendRef.BackendObjectReference

		var backend *ir.BackendObjectIR
		var err error
		if disableExtensionRefValidation {
			backend, err = backends.GetBackendFromRefWithoutRefGrantValidation(krtctx, objectSource, backendRef)
		} else {
			backend, err = backends.GetBackendFromRef(krtctx, objectSource, backendRef)
		}
		if err != nil {
			return nil, err
		}
		if backend != nil {
			clusterName = backend.ClusterName()
		}
		if grpcService.Authority != nil {
			authority = *grpcService.Authority
		}
	}
	if clusterName == "" {
		return nil, errors.New("backend not found")
	}
	envoyGrpcService := &envoy_core_v3.GrpcService{
		TargetSpecifier: &envoy_core_v3.GrpcService_EnvoyGrpc_{
			EnvoyGrpc: &envoy_core_v3.GrpcService_EnvoyGrpc{
				ClusterName: clusterName,
				Authority:   authority,
			},
		},
	}
	return envoyGrpcService, nil
}

func NewGatewayTranslationPass(ctx context.Context, tctx ir.GwTranslationCtx, reporter reports.Reporter) ir.ProxyTranslationPass {
	return &trafficPolicyPluginGwPass{
		reporter: reporter,
	}
}

func (p *trafficPolicyPluginGwPass) ApplyListenerPlugin(ctx context.Context, pCtx *ir.ListenerContext, out *envoy_config_listener_v3.Listener) {
	policy, ok := pCtx.Policy.(*TrafficPolicy)
	if !ok {
		return
	}

	if policy.spec.extAuth != nil && policy.spec.extAuth.provider != nil {
		if p.extAuthPerProvider == nil {
			p.extAuthPerProvider = make(map[string]providerWithFromListener)
		}
		k := policy.spec.extAuth.provider.ResourceName()
		p.extAuthPerProvider[k] = providerWithFromListener{
			provider:     policy.spec.extAuth.provider,
			fromListener: true,
		}
	}
	if policy.spec.extProc != nil && policy.spec.extProc.provider != nil {
		if p.extProcPerProvider == nil {
			p.extProcPerProvider = make(map[string]providerWithFromListener)
		}
		k := policy.spec.extAuth.provider.ResourceName()
		p.extAuthPerProvider[k] = providerWithFromListener{
			provider:     policy.spec.extAuth.provider,
			fromListener: true,
		}
	}
	if policy.spec.globalRateLimit != nil && policy.spec.globalRateLimit.provider != nil {
		if p.rateLimitPerProvider == nil {
			p.rateLimitPerProvider = make(map[string]providerWithFromListener)
		}
		k := policy.spec.globalRateLimit.provider.ResourceName()
		p.rateLimitPerProvider[k] = providerWithFromListener{
			provider:     policy.spec.globalRateLimit.provider,
			fromListener: true,
		}
	}

	if policy.spec.localRateLimit != nil {
		p.localRateLimitInChain = policy.spec.localRateLimit.localRateLimit
	}
	if policy.spec.transform != nil {
		p.setTransformationInChain = true
		p.listenerTransform = policy.spec.transform
		if p.rustformationStash == nil {
			p.rustformationStash = make(map[string]string)
		}
		// routes cannot have an empty string id so we special case here
		p.rustformationStash[""] = policy.spec.rustformationStringToStash
	}
}

func (p *trafficPolicyPluginGwPass) ApplyForRoute(ctx context.Context, pCtx *ir.RouteContext, outputRoute *routev3.Route) error {
	policy, ok := pCtx.Policy.(*TrafficPolicy)
	if !ok {
		return nil
	}

	if policy.spec.rustformation != nil {
		// TODO(nfuden): get back to this path once we have valid perroute
		// pCtx.TypedFilterConfig.AddTypedConfig(rustformationFilterNamePrefix, policy.spec.rustformation)

		// Hack around not having route level.
		// Note this is really really bad and rather fragile due to listener draining behaviors
		routeHash := strconv.Itoa(int(utils.HashProto(outputRoute)))
		if p.rustformationStash == nil {
			p.rustformationStash = make(map[string]string)
		}
		// encode the configuration that would be route level and stash the serialized version in a map
		p.rustformationStash[routeHash] = string(policy.spec.rustformationStringToStash)

		// augment the dynamic metadata so that we can do our route hack
		// set_dynamic_metadata filter DOES NOT have a route level configuration
		// set_filter_state can be used but the dynamic modules cannot access it on the current version of envoy
		// therefore use the old transformation just for rustformation
		reqm := &transformationpb.RouteTransformations_RouteTransformation_RequestMatch{
			RequestTransformation: &transformationpb.Transformation{
				TransformationType: &transformationpb.Transformation_TransformationTemplate{
					TransformationTemplate: &transformationpb.TransformationTemplate{
						ParseBodyBehavior: transformationpb.TransformationTemplate_DontParse, // Default is to try for JSON... Its kinda nice but failure is bad...
						DynamicMetadataValues: []*transformationpb.TransformationTemplate_DynamicMetadataValue{
							{
								MetadataNamespace: "kgateway",
								Key:               "route",
								Value: &transformationpb.InjaTemplate{
									Text: routeHash,
								},
							},
						},
					},
				},
			},
		}

		setmetaTransform := &transformationpb.RouteTransformations{
			Transformations: []*transformationpb.RouteTransformations_RouteTransformation{
				{
					Match: &transformationpb.RouteTransformations_RouteTransformation_RequestMatch_{
						RequestMatch: reqm,
					},
				},
			},
		}
		pCtx.TypedFilterConfig.AddTypedConfig(metadataRouteTransformation, setmetaTransform)

		p.setTransformationInChain = true
	}

	if policy.spec.ai != nil {
		var aiBackends []*v1alpha1.Backend
		// check if the backends selected by targetRef are all AI backends before applying the policy
		for _, backend := range pCtx.In.Backends {
			if backend.Backend.BackendObject == nil {
				// could be nil if not found or no ref grant
				continue
			}
			b, ok := backend.Backend.BackendObject.Obj.(*v1alpha1.Backend)
			if !ok {
				// AI policy cannot apply to kubernetes services
				// TODO(npolshak): Report this as a warning on status
				logger.Warn("AI Policy cannot apply to kubernetes services", "backend_name", backend.Backend.BackendObject.GetName())
				continue
			}
			if b.Spec.Type != v1alpha1.BackendTypeAI {
				// AI policy cannot apply to non-AI backends
				// TODO(npolshak): Report this as a warning on status
				logger.Warn("AI Policy cannot apply to non-AI backend", "backend_name", backend.Backend.BackendObject.GetName(), "backend_type", string(b.Spec.Type))
				continue
			}
			aiBackends = append(aiBackends, b)
		}
		if len(aiBackends) > 0 {
			// Apply the AI policy to the all AI backends
			p.processAITrafficPolicy(&pCtx.TypedFilterConfig, policy.spec.ai)
		}
	}
	p.handlePolicies(&pCtx.TypedFilterConfig, policy.spec)

	return nil
}

func (p *trafficPolicyPluginGwPass) ApplyForRouteBackend(
	ctx context.Context,
	policy ir.PolicyIR,
	pCtx *ir.RouteBackendContext,
) error {
	rtPolicy, ok := policy.(*TrafficPolicy)
	if !ok {
		return nil
	}

	p.handlePolicies(&pCtx.TypedFilterConfig, rtPolicy.spec)

	if rtPolicy.spec.ai != nil && (rtPolicy.spec.ai.transformation != nil || rtPolicy.spec.ai.extProc != nil) {
		p.processAITrafficPolicy(&pCtx.TypedFilterConfig, rtPolicy.spec.ai)
	}

	return nil
}

func (p *trafficPolicyPluginGwPass) HttpFilters(ctx context.Context, fcc ir.FilterChainCommon) ([]plugins.StagedHttpFilter, error) {
	filters := []plugins.StagedHttpFilter{}

	// Add Ext_proc filters for listener
	for providerName, providerExtProc := range p.extProcPerProvider {
		extProcFilter := providerExtProc.provider.ExtProc
		if extProcFilter == nil {
			continue
		}

		// add the specific auth filter
		extProcName := extProcFilterName(providerName)
		stagedExtProcFilter := plugins.MustNewStagedFilter(extProcName,
			extProcFilter,
			plugins.AfterStage(plugins.WellKnownFilterStage(plugins.AuthZStage)))

		// handle the two enable attachement cases
		if !providerExtProc.fromListener {
			// handle the case where route level only should be fired
			stagedExtProcFilter.Filter.Disabled = true
		}

		filters = append(filters, stagedExtProcFilter)
	}

	// register classic transforms
	if p.setTransformationInChain && !useRustformations {
		// TODO(nfuden): support stages such as early
		transformationCfg := transformationpb.FilterTransformations{}
		if p.listenerTransform != nil {
			convertClassicRouteToListener(&transformationCfg, p.listenerTransform)
		}
		filters = append(filters, plugins.MustNewStagedFilter(transformationFilterNamePrefix,
			&transformationCfg,
			plugins.BeforeStage(plugins.AcceptedStage)))
	}
	if p.setTransformationInChain && useRustformations {
		// ---------------
		// | END CLASSIC |
		// ---------------
		// TODO(nfuden/yuvalk): how to do route level correctly probably contribute to dynamic module upstream
		// smash together configuration
		filterRouteHashConfig := map[string]string{}
		topLevel, ok := p.rustformationStash[""]

		if topLevel == "" {
			topLevel = "}"
		} else {
			// toplevel is already formatted and at this point its quicker to rip off the { than it is so unmarshal and all}
			topLevel = "," + topLevel[1:]
		}
		if ok {
			delete(p.rustformationStash, "")
		}
		for k, v := range p.rustformationStash {
			filterRouteHashConfig[k] = v
		}

		filterConfig, _ := json.Marshal(filterRouteHashConfig)
		msg, _ := utils.MessageToAny(&wrapperspb.StringValue{
			Value: fmt.Sprintf(`{"route_specific": %s%s`, string(filterConfig), topLevel),
		})
		rustCfg := dynamicmodulesv3.DynamicModuleFilter{
			DynamicModuleConfig: &exteniondynamicmodulev3.DynamicModuleConfig{
				Name: "rust_module",
			},
			FilterName: "http_simple_mutations",

			// currently we use stringvalue but we should look at using the json variant as supported in upstream
			FilterConfig: msg,
		}

		filters = append(filters, plugins.MustNewStagedFilter(rustformationFilterNamePrefix,
			&rustCfg,
			plugins.BeforeStage(plugins.AcceptedStage)))

		// filters = append(filters, plugins.MustNewStagedFilter(setFilterStateFilterName,
		// 	&set_filter_statev3.Config{}, plugins.AfterStage(plugins.FaultStage)))
		filters = append(filters, plugins.MustNewStagedFilter(metadataRouteTransformation,
			&transformationpb.FilterTransformations{},
			plugins.AfterStage(plugins.FaultStage)))
	}

	// register the transformation work once
	if len(p.extAuthPerProvider) != 0 {
		// register the filter that sets metadata so that it can have overrides on the route level
		f := plugins.MustNewStagedFilter(extAuthGlobalDisableFilterName,
			setMetadataConfig,
			plugins.BeforeStage(plugins.FaultStage))
		f.Filter.Disabled = true
		filters = append(filters, f)
	}

	// Add Ext_authz filter for listener
	for providerName, providerExtauth := range p.extAuthPerProvider {
		extAuthFilter := providerExtauth.provider.ExtAuth
		if extAuthFilter == nil {
			continue
		}

		// add the specific auth filter
		extauthName := extAuthFilterName(providerName)
		stagedExtAuthFilter := plugins.MustNewStagedFilter(extauthName,
			extAuthFilter,
			plugins.DuringStage(plugins.AuthZStage))

		// handle the two enable attachement cases
		if !providerExtauth.fromListener {
			// handle the case where route level only should be fired
			stagedExtAuthFilter.Filter.Disabled = true
		}

		filters = append(filters, stagedExtAuthFilter)
	}

	if p.localRateLimitInChain != nil {
		filters = append(filters, plugins.MustNewStagedFilter(localRateLimitFilterNamePrefix,
			p.localRateLimitInChain,
			plugins.BeforeStage(plugins.AcceptedStage)),
		)
	}

	// Add global rate limit filters from providers
	for providerName, providerRateLimit := range p.rateLimitPerProvider {
		rateLimitFilter := providerRateLimit.provider.RateLimit
		if rateLimitFilter == nil {
			continue
		}

		// add the specific rate limit filter with a unique name
		rateLimitName := getRateLimitFilterName(providerName)
		stagedRateLimitFilter := plugins.MustNewStagedFilter(rateLimitName,
			rateLimitFilter,
			plugins.DuringStage(plugins.RateLimitStage))

		// If this rate limit is not from a listener, disable it at the listener level
		// so it can be enabled selectively at the route level
		if !providerRateLimit.fromListener {
			stagedRateLimitFilter.Filter.Disabled = true
		}

		filters = append(filters, stagedRateLimitFilter)
	}

	if len(filters) == 0 {
		return nil, nil
	}
	return filters, nil
}

func (p *trafficPolicyPluginGwPass) handleTransformation(typedFilterConfig *ir.TypedFilterConfigMap, transform *transformationpb.RouteTransformations) {
	if transform == nil {
		return
	}
	typedFilterConfig.AddTypedConfig(transformationFilterNamePrefix, transform)
	p.setTransformationInChain = true
}

func (p *trafficPolicyPluginGwPass) handleLocalRateLimit(typedFilterConfig *ir.TypedFilterConfigMap, localRateLimit *localRateLimitIR) {
	if localRateLimit == nil {
		return
	}
	if localRateLimit.localRateLimit != nil {
		typedFilterConfig.AddTypedConfig(localRateLimitFilterNamePrefix, localRateLimit.localRateLimit)
	}

	// Add a filter to the chain. When having a rate limit for a route we need to also have a
	// globally disabled rate limit filter in the chain otherwise it will be ignored.
	// If there is also rate limit for the listener, it will not override this one.
	if p.localRateLimitInChain == nil {
		p.localRateLimitInChain = &localratelimitv3.LocalRateLimit{
			StatPrefix: localRateLimitStatPrefix,
		}
	}
}

func (p *trafficPolicyPluginGwPass) handlePolicies(typedFilterConfig *ir.TypedFilterConfigMap, spec trafficPolicySpecIr) {
	p.handleTransformation(typedFilterConfig, spec.transform)
	// Apply ExtAuthz configuration if present
	// ExtAuth does not allow for most information such as destination
	// to be set at the route level so we need to smuggle info upwards.
	p.handleExtAuth(typedFilterConfig, spec.extAuth)
	p.handleExtProc(typedFilterConfig, spec.extProc)
	// Apply rate limit configuration if present
	p.handleRateLimit(typedFilterConfig, spec.globalRateLimit)
	p.handleLocalRateLimit(typedFilterConfig, spec.localRateLimit)
}

// handleRateLimit adds rate limit configurations to routes
func (p *trafficPolicyPluginGwPass) handleRateLimit(typedFilterConfig *ir.TypedFilterConfigMap, rateLimit *globalRateLimitIR) {
	if rateLimit == nil {
		return
	}
	if rateLimit.rateLimitActions == nil {
		return
	}

	providerName := rateLimit.provider.ResourceName()

	// Initialize the map if it doesn't exist yet
	if p.rateLimitPerProvider == nil {
		p.rateLimitPerProvider = make(map[string]providerWithFromListener)
	}
	if _, ok := p.rateLimitPerProvider[providerName]; !ok {
		p.rateLimitPerProvider[providerName] = providerWithFromListener{
			provider: rateLimit.provider,
		}
	}

	// Configure rate limit per route - enabling it for this specific route
	rateLimitPerRoute := &ratev3.RateLimitPerRoute{
		RateLimits: rateLimit.rateLimitActions,
	}
	typedFilterConfig.AddTypedConfig(getRateLimitFilterName(providerName), rateLimitPerRoute)
}

func (p *trafficPolicyPluginGwPass) handleExtAuth(pCtxTypedFilterConfig *ir.TypedFilterConfigMap, extAuth *extAuthIR) {
	if extAuth == nil {
		return
	}

	// Handle the enablement state
	if extAuth.enablement == v1alpha1.ExtAuthDisableAll {
		// Disable the filter under all providers via the metadata match
		// we have to use the metadata as we dont know what other configurations may have extauth
		pCtxTypedFilterConfig.AddTypedConfig(extAuthGlobalDisableFilterName, enableFilterPerRoute)
	} else {
		providerName := extAuth.provider.ResourceName()
		if extAuth.extauthPerRoute != nil {
			pCtxTypedFilterConfig.AddTypedConfig(extAuthFilterName(providerName),
				extAuth.extauthPerRoute,
			)
		} else if !p.extAuthPerProvider[providerName].fromListener {
			// if you are on a route and not trying to disable it then we need to override the top level disable on the filter chain
			pCtxTypedFilterConfig.AddTypedConfig(extAuthFilterName(providerName),
				&envoy_ext_authz_v3.ExtAuthzPerRoute{
					Override: &envoy_ext_authz_v3.ExtAuthzPerRoute_CheckSettings{
						CheckSettings: &envoy_ext_authz_v3.CheckSettings{},
					},
				},
			)
		}
		if p.extAuthPerProvider == nil {
			p.extAuthPerProvider = make(map[string]providerWithFromListener)
		}
		if _, ok := p.extAuthPerProvider[providerName]; !ok {
			p.extAuthPerProvider[providerName] = providerWithFromListener{
				provider: extAuth.provider,
			}
		}
	}
}

func (p *trafficPolicyPluginGwPass) handleExtProc(pCtxTypedFilterConfig *ir.TypedFilterConfigMap, extProc *extProcIR) {
	if extProc == nil || extProc.provider == nil {
		return
	}
	providerName := extProc.provider.ResourceName()
	// Handle the enablement state
	if extProc.extProcPerRoute != nil {
		pCtxTypedFilterConfig.AddTypedConfig(extProcFilterName(providerName),
			extProc.extProcPerRoute,
		)
	} else if !p.extProcPerProvider[providerName].fromListener {
		// if you are on a route and not trying to disable it then we need to override the top level disable on the filter chain
		pCtxTypedFilterConfig.AddTypedConfig(extProcFilterName(providerName),
			&envoy_ext_proc_v3.ExtProcPerRoute{Override: &envoy_ext_proc_v3.ExtProcPerRoute_Overrides{Overrides: &envoy_ext_proc_v3.ExtProcOverrides{}}},
		)
	}

	if p.extProcPerProvider == nil {
		p.extProcPerProvider = make(map[string]providerWithFromListener)
	}
	if _, ok := p.extProcPerProvider[providerName]; !ok {
		p.extProcPerProvider[providerName] = providerWithFromListener{
			provider: extProc.provider,
		}
	}
}

// mergePolicies merges the given policy ordered from high to low priority (both hierarchically
// and within the same hierarchy) based on the constraints defined per PolicyAtt.
//
// It iterates policies in reverse order (low to high) to ensure higher priority policies can
// always use an OverridableMerge strategy to override lower priority ones. Iterating policies
// in the given priority order (high to low) requires more complex merging for delegated chains
// because policies anywhere in the chain may enable policy overrides for their children but we
// still need to ensure these children cannot override any policies set by their ancestors that
// are not marked as overridable, i.e., (r1,p1)-delegate->(r2,p2)-delegate->(r3,p3) where
// r=route p=policy needs to ensure p3 does not override p1 (assuming p1 does not enable overrides)
// even if p2 allows overrides. This is easier to guarantee by using an OverridableMerge strategy
// by merging in higher priority policies with different HierarchicalPriority.
func mergePolicies(policies []ir.PolicyAtt) ir.PolicyAtt {
	var out ir.PolicyAtt
	if len(policies) == 0 {
		return out
	}
	_, ok := policies[0].PolicyIr.(*TrafficPolicy)
	// ignore unknown types
	if !ok {
		return out
	}

	// base policy to merge into has an empty PolicyIr so it can always be merged into
	out = ir.PolicyAtt{
		GroupKind:    policies[0].GroupKind,
		PolicyRef:    policies[0].PolicyRef,
		MergeOrigins: map[string]*ir.AttachedPolicyRef{},
		PolicyIr:     &TrafficPolicy{},
	}
	merged := out.PolicyIr.(*TrafficPolicy)

	for i := len(policies) - 1; i >= 0; i-- {
		mergeOpts := policy.MergeOptions{
			Strategy: policy.OverridableMerge,
		}
		// If merging a policy lower in the hierarchy with a policy higher in the hierarchy AND
		// the policy higher in the hierarchy enables policy overrides, use an AugmentedMerge strategy
		// to preserve existing fields set by lower levels.
		// NOTE: the HierarchicalPriority check is necessary to prevent enabling override behavior among
		// policies in the same hierarchy, e.g., ExtensionRef vs TargetRef policy attached to the same route, as
		// DelegationInheritedPolicyPriorityPreferChild strictly applies to parent->child policy inheritance and is not applicable
		// outside delegated policy inheritance.
		if out.HierarchicalPriority < policies[i].HierarchicalPriority && policies[i].DelegationInheritedPolicyPriority == apiannotations.DelegationInheritedPolicyPriorityPreferChild {
			mergeOpts.Strategy = policy.AugmentedMerge
		}

		p2 := policies[i].PolicyIr.(*TrafficPolicy)
		p2Ref := policies[i].PolicyRef

		if policy.IsMergeable(merged.spec.ai, p2.spec.ai, mergeOpts) {
			merged.spec.ai = p2.spec.ai
			out.MergeOrigins["ai"] = p2Ref
		}
		if policy.IsMergeable(merged.spec.extProc, p2.spec.extProc, mergeOpts) {
			merged.spec.extProc = p2.spec.extProc
			out.MergeOrigins["extProc"] = p2Ref
		}
		if policy.IsMergeable(merged.spec.transform, p2.spec.transform, mergeOpts) {
			merged.spec.transform = p2.spec.transform
			out.MergeOrigins["transformation"] = p2Ref
		}
		if policy.IsMergeable(merged.spec.rustformation, p2.spec.rustformation, mergeOpts) {
			merged.spec.rustformation = p2.spec.rustformation
			merged.spec.rustformationStringToStash = p2.spec.rustformationStringToStash
			out.MergeOrigins["rustformation"] = p2Ref
		}
		if policy.IsMergeable(merged.spec.extAuth, p2.spec.extAuth, mergeOpts) {
			merged.spec.extAuth = p2.spec.extAuth
			out.MergeOrigins["extAuth"] = p2Ref
		}
		if policy.IsMergeable(merged.spec.localRateLimit, p2.spec.localRateLimit, mergeOpts) {
			merged.spec.localRateLimit = p2.spec.localRateLimit
			out.MergeOrigins["rateLimit"] = p2Ref
		}
		// Handle global rate limit merging
		if policy.IsMergeable(merged.spec.globalRateLimit, p2.spec.globalRateLimit, mergeOpts) {
			merged.spec.globalRateLimit = p2.spec.globalRateLimit
			out.MergeOrigins["rateLimit"] = p2Ref
		}

		out.HierarchicalPriority = policies[i].HierarchicalPriority
	}

	return out
}
