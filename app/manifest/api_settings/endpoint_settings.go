package api_settings

import (
	"fmt"
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type EndpointSettings struct {
	Route           string           `yaml:"route"`
	Method          types.HttpMethod `yaml:"method"`
	ActionID        string           `yaml:"actionId"`
	FlowActionID    []string         `yaml:"flowActionId"`
	EnableAnonymous bool             `yaml:"enableAnonymous"`
	Roles           []string         `yaml:"roles"`
	Scopes          []string         `yaml:"scopes"`
	Claims          []string         `yaml:"claims"`
	IsProxy         bool             `yaml:"isProxy"`
	IdempId         string           `yaml:"idempId"`
	RateLimitId     string           `yaml:"rateLimitId"`
}

func (setting EndpointSettings) ShouldConfigureAuthorization(apiHasAuthorization bool) bool {
	return apiHasAuthorization && !setting.EnableAnonymous
}

func (setting EndpointSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Route) == 0 {
		var msg = fmt.Sprintf("api.endpoints[%s].route is required", setting.ActionID)
		result.AddError(msg)
	}

	if len(setting.ActionID) == 0 && len(setting.FlowActionID) == 0 {
		var msg = fmt.Sprintf("Should be informed an api.endpoints[%s].actionId or api.endpoints[%s].flowActionId", setting.Route, setting.ActionID)
		result.AddError(msg)
	}

	if !setting.IsProxy {

		if setting.Route[0] != '/' {
			var msg = fmt.Sprintf("api.endpoints[%s].route should start with '/' ex: /api/hello", setting.ActionID)
			result.AddError(msg)
		}

		if len(setting.Method) == 0 {
			var msg = fmt.Sprintf("api.endpoints[%s].method is required", setting.ActionID)
			result.AddError(msg)
		} else {
			if !(setting.Method == types.HttpMethodGet ||
				setting.Method == types.HttpMethodPost ||
				setting.Method == types.HttpMethodPut ||
				setting.Method == types.HttpMethodPatch ||
				setting.Method == types.HttpMethodDelete) {

				var msg = fmt.Sprintf("api.endpoints[%s].method should contain valid value (get, post, put, patch or delete)", setting.ActionID)
				result.AddError(msg)
			}
		}

	} else {
		if len(setting.Method) > 0 {
			var msg = fmt.Sprintf("api.endpoints[%s].method could not be defined when is proxy", setting.ActionID)
			result.AddError(msg)
		}
	}
	return result
}
