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
	TlsId             string            `yaml:"tlsId"`
}

func (setting *HttpRequestSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Url) == 0 {
		result.AddError("actions.http.request.url is required")
	}

	return result
}
