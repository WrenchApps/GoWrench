package http_settings

import (
	"fmt"
	"wrench/app/manifest/validation"
)

type HttpSetting struct {
	Request  *HttpRequestSetting     `yaml:"request"`
	Response *HttpResponseSettings   `yaml:"response"`
	Mock     *HttpRequestMockSetting `yaml:"mock"`
}

func (setting HttpSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Request == nil && setting.Mock == nil {
		result.AddError("actions.http or actions.mock should be configured")
	}

	if setting.Request != nil {
		result.AppendValidable(setting.Request)
	}

	if setting.Response != nil {
		result.AppendValidable(setting.Response)
	}

	if setting.Mock != nil {
		result.AppendValidable(setting.Mock)
	}

	return result
}

func (setting HttpSetting) ValidTypeActionTypeHttpRequest(result *validation.ValidateResult) {

	if setting.Request == nil {
		var msg = fmt.Sprintf("actions.http.request should be nil")
		result.AddError(msg)
	}

	if setting.Mock != nil {
		var msg = fmt.Sprintf("actions.http.mock should be nil")
		result.AddError(msg)
	}
}

func (setting HttpSetting) ValidTypeActionTypeHttpRequestMock(result *validation.ValidateResult) {
	if setting.Request != nil {
		var msg = fmt.Sprintf("actions.http.request should be nil")
		result.AddError(msg)
	}

	if setting.Mock == nil {
		var msg = fmt.Sprintf("actions.http.mock is required")
		result.AddError(msg)
	}
}
