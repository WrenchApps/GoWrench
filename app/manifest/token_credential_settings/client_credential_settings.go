package token_credential_settings

import "wrench/app/manifest/validation"

type ClientCredentialSetting struct {
	ClientId     string `yaml:"clientId"`
	ClientSecret string `yaml:"clientSecret"`
	Audiance     string `yaml:"audiance"`
}

func (setting ClientCredentialSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.ClientId) == 0 {
		result.AddError("tokenCredentials.clientCredential.clientId is required")
	}

	if len(setting.ClientSecret) == 0 {
		result.AddError("tokenCredentials.clientCredential.clientSecret is required")
	}

	return result
}
