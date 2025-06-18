package token_credential_settings

import (
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type CustomAuthenticationConfigurations struct {
	AccessTokenPropertyName string `yaml:"accessTokenPropertyName"`
	TokenType               string `yaml:"tokenType"`
	HeaderName              string `yaml:"headerName"`
}

type CustomAuthentication struct {
	Method         types.HttpMethod                   `yaml:"method"`
	RequestBody    map[string]string                  `yaml:"requestBody"`
	RequestHeaders map[string]string                  `yaml:"requestHeaders"`
	Configs        CustomAuthenticationConfigurations `yaml:"configs"`
}

func (setting CustomAuthentication) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Configs.AccessTokenPropertyName) == 0 {
		result.AddError("tokenCredentials.custom_authentication.custom.configs.AccessTokenPropertyName is required")
	}

	if len(setting.Method) == 0 {
		result.AddError("tokenCredentials.custom_authentication.custom.method is required")
	} else {
		if (setting.Method == types.HttpMethodGet ||
			setting.Method == types.HttpMethodPost ||
			setting.Method == types.HttpMethodPut ||
			setting.Method == types.HttpMethodPatch ||
			setting.Method == types.HttpMethodDelete) == false {

			result.AddError("tokenCredentials.custom_authentication.custom.method should contain valid value (get, post, put, patch or delete)")
		}
	}

	return result
}
