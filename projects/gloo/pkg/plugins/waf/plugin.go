package waf

import (
	envoyroute "github.com/envoyproxy/go-control-plane/envoy/api/v2/route"
	"github.com/gogo/protobuf/proto"
	v1 "github.com/solo-io/gloo/projects/gloo/pkg/api/v1"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/pluginutils"
	"github.com/solo-io/gloo/projects/gloo/pkg/plugins/utils"
	"github.com/solo-io/go-utils/errors"
	. "github.com/solo-io/solo-projects/projects/gloo/pkg/api/external/envoy/waf"
	"github.com/solo-io/solo-projects/projects/gloo/pkg/api/v1/plugins/waf"
)

const (
	FilterName    = "io.solo.filters.http.modsecurity"
	ExtensionName = "waf"
)

type Plugin struct {
	listenerEnabled map[*v1.HttpListener]bool
}

var (
	_ plugins.Plugin            = new(Plugin)
	_ plugins.VirtualHostPlugin = new(Plugin)
	_ plugins.RoutePlugin       = new(Plugin)
	_ plugins.HttpFilterPlugin  = new(Plugin)

	// waf should happen before any code is run
	filterStage = plugins.DuringStage(plugins.WafStage)
)

var (
	ConvertingProtoError = func(err error, cfg string) error {
		if err == utils.NotFoundError {
			return nil
		}
		return errors.Wrapf(err, "Error converting proto to %s waf plugin", cfg)
	}
)

func NewPlugin() *Plugin {
	return &Plugin{
		listenerEnabled: make(map[*v1.HttpListener]bool),
	}
}

func (p *Plugin) Init(params plugins.InitParams) error {
	return nil
}

func (p *Plugin) addListener(listener *v1.HttpListener) {
	p.listenerEnabled[listener] = true
}

func (p *Plugin) listenerPresent(listener *v1.HttpListener) bool {
	val, ok := p.listenerEnabled[listener]
	if !ok {
		return false
	}
	return val
}

// Process virtual host plugin
func (p *Plugin) ProcessVirtualHost(params plugins.VirtualHostParams, in *v1.VirtualHost, out *envoyroute.VirtualHost) error {
	var wafConfig waf.VhostSettings
	err := utils.UnmarshalExtension(in.VirtualHostPlugins, ExtensionName, &wafConfig)
	if err != nil {
		return ConvertingProtoError(err, "virtual host")
	}

	// should never be nil
	p.addListener(params.Listener.GetHttpListener())

	perVhostCfg := &ModSecurityPerRoute{
		Disabled: wafConfig.Disabled,
	}

	if wafConfig.GetSettings() != nil {
		perVhostCfg.RuleSets = wafConfig.GetSettings().GetRuleSets()
		if coreRuleSet := getCoreRuleSet(wafConfig.GetSettings().GetCoreRuleSet()); coreRuleSet != nil {
			perVhostCfg.RuleSets = append(perVhostCfg.RuleSets, coreRuleSet)
		}
	}

	pluginutils.SetVhostPerFilterConfig(out, FilterName, perVhostCfg)

	return nil
}

// Process route plugin
func (p *Plugin) ProcessRoute(params plugins.RouteParams, in *v1.Route, out *envoyroute.Route) error {
	var wafConfig waf.RouteSettings
	err := utils.UnmarshalExtension(in.RoutePlugins, ExtensionName, &wafConfig)
	if err != nil {
		return ConvertingProtoError(err, "route")
	}

	p.addListener(params.Listener.GetHttpListener())

	perRouteCfg := &ModSecurityPerRoute{
		Disabled: wafConfig.Disabled,
	}

	if wafConfig.GetSettings() != nil {
		perRouteCfg.RuleSets = wafConfig.GetSettings().GetRuleSets()
		if coreRuleSet := getCoreRuleSet(wafConfig.GetSettings().GetCoreRuleSet()); coreRuleSet != nil {
			perRouteCfg.RuleSets = append(perRouteCfg.RuleSets, coreRuleSet)
		}
	}

	pluginutils.SetRoutePerFilterConfig(out, FilterName, perRouteCfg)
	return nil
}

// Http Filter to return the waf filter
func (p *Plugin) HttpFilters(params plugins.Params, listener *v1.HttpListener) ([]plugins.StagedHttpFilter, error) {
	var filters []plugins.StagedHttpFilter
	// If the list does not already have the listener than it is necessary to check for nil
	if !p.listenerPresent(listener) {
		if listener.GetListenerPlugins() == nil {
			return nil, nil
		}
	}

	var settings waf.Settings
	err := utils.UnmarshalExtension(listener.GetListenerPlugins(), ExtensionName, &settings)
	if err != nil {
		// important to check this because, on any other error this method should keep executing
		if err != utils.NotFoundError {
			return nil, ConvertingProtoError(err, "listener")
		}
	}

	modSecurityConfig := &ModSecurity{}

	if settings.GetCoreRuleSet() == nil && settings.GetRuleSets() == nil {
		modSecurityConfig.Disabled = true
	} else {
		modSecurityConfig.RuleSets = settings.GetRuleSets()
		modSecurityConfig.Disabled = settings.GetDisabled()

		if coreRuleSet := getCoreRuleSet(settings.GetCoreRuleSet()); coreRuleSet != nil {
			modSecurityConfig.RuleSets = append(modSecurityConfig.RuleSets, coreRuleSet)
		}
	}

	var cfg proto.Message = modSecurityConfig

	stagedFilter, err := plugins.NewStagedFilterWithConfig(FilterName, cfg, filterStage)
	if err != nil {
		return nil, err
	}
	filters = append(filters, stagedFilter)
	return filters, nil
}

func getCoreRuleSet(crs *waf.CoreRuleSet) *RuleSet {
	if crs == nil {
		return nil
	}
	coreRuleSet := &RuleSet{
		Files: getCoreRuleSetFiles(),
	}
	switch additionalSettings := crs.GetCustomSettingsType().(type) {
	case *waf.CoreRuleSet_CustomSettingsString:
		coreRuleSet.RuleStr = additionalSettings.CustomSettingsString
	case *waf.CoreRuleSet_CustomSettingsFile:
		coreRuleSet.Files = append([]string{additionalSettings.CustomSettingsFile}, coreRuleSet.Files...)
	}
	return coreRuleSet
}
