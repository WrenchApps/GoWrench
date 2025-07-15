package http_settings

import (
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type HttpRequestSetting struct {
	Method            types.HttpMethod  `yaml:"method"`
	Url               string            `yaml:"url"`
	Headers           map[string]string `yaml:"headers"`
	TokenCredentialId string            `yaml:"tokenCredentialId"`
	Insecure          bool              `yaml:"insecure"`
}

func (setting HttpRequestSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Method) == 0 {
		result.AddError("actions.http.request.method is required")
	} else {
		if (setting.Method == types.HttpMethodGet ||
			setting.Method == types.HttpMethodPost ||
			setting.Method == types.HttpMethodPut ||
			setting.Method == types.HttpMethodPatch ||
			setting.Method == types.HttpMethodDelete) == false {

			result.AddError("actions.http.request.method should contain valid value (get, post, put, patch or delete)")
		}
	}

	if len(setting.Url) == 0 {
		result.AddError("actions.http.request.url is required")
	}

	return result
}
