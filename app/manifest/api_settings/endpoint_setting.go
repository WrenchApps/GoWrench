package api_settings

import (
	"fmt"
	"wrench/app/manifest/validation"
)

type EndpointSetting struct {
	Route        string     `yaml:"route"`
	Method       HttpMethod `yaml:"method"`
	ActionID     string     `yaml:"actionId"`
	FlowActionID string     `yaml:"flowActionId"`
}

type HttpMethod string

const (
	HttpMethodGet    HttpMethod = "get"
	HttpMethodPost   HttpMethod = "post"
	HttpMethodPut    HttpMethod = "put"
	HttpMethodPatch  HttpMethod = "patch"
	HttpMethodDelete HttpMethod = "delete"
)

func (setting EndpointSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Route) == 0 {
		var msg = fmt.Sprintf("api.endpoints[%s].route is required", setting.ActionID)
		result.AddError(msg)
	}

	if len(setting.Method) == 0 {
		var msg = fmt.Sprintf("api.endpoints[%s].method is required", setting.ActionID)
		result.AddError(msg)
	} else {
		if (setting.Method == HttpMethodGet ||
			setting.Method == HttpMethodPost ||
			setting.Method == HttpMethodPut ||
			setting.Method == HttpMethodPatch ||
			setting.Method == HttpMethodDelete) == false {

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
