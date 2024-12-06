package api_settings

import (
	"fmt"
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type EndpointSettings struct {
	Route        string           `yaml:"route"`
	Method       types.HttpMethod `yaml:"method"`
	ActionID     string           `yaml:"actionId"`
	FlowActionID string           `yaml:"flowActionId"`
}

func (setting EndpointSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Route) == 0 {
		var msg = fmt.Sprintf("api.endpoints[%s].route is required", setting.ActionID)
		result.AddError(msg)
	}

	if setting.Route[0] != '/' {
		var msg = fmt.Sprintf("api.endpoints[%s].route should start with '/' ex: /api/hello", setting.ActionID)
		result.AddError(msg)
	}

	if len(setting.Method) == 0 {
		var msg = fmt.Sprintf("api.endpoints[%s].method is required", setting.ActionID)
		result.AddError(msg)
	} else {
		if (setting.Method == types.HttpMethodGet ||
			setting.Method == types.HttpMethodPost ||
			setting.Method == types.HttpMethodPut ||
			setting.Method == types.HttpMethodPatch ||
			setting.Method == types.HttpMethodDelete) == false {

			var msg = fmt.Sprintf("api.endpoints[%s].method should contain valid value (get, post, put, patch or delete)", setting.ActionID)
			result.AddError(msg)
		}
	}

	if len(setting.ActionID) == 0 && len(setting.FlowActionID) == 0 {
		var msg = fmt.Sprintf("Should be informed an api.endpoints[%s].actionId or api.endpoints[%s].flowActionId", setting.ActionID, setting.ActionID)
		result.AddError(msg)
	}
	return result
}
