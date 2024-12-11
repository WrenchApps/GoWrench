package token_credential

import (
	"fmt"
	"strings"
	"wrench/app/manifest/validation"
)

type TokenCredentialSetting struct {
	Id           string              `yaml:"id"`
	Type         TokenCredentialType `yaml:"type"`
	AuthEndpoint string              `yaml:"authEndpoint"`
	ClientId     string              `yaml:"clientId"`
	ClientSecret string              `yaml:"clientSecret"`
}

type TokenCredentialType string

const (
	TokenCredentialClientCredential TokenCredentialType = "clientCredentials"
)

func (setting TokenCredentialSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Id) == 0 {
		result.AddError("tokenCredentials.id is required")
	} else {
		if strings.ContainsAny(setting.Id, " ") {
			result.AddError("tokenCredentials.id can't contains space")
		}
	}

	if len(setting.AuthEndpoint) == 0 {
		result.AddError("tokenCredentials.authEndpoint is required")
	}

	if len(setting.ClientId) == 0 {
		result.AddError("tokenCredentials.clientId is required")
	}

	if len(setting.ClientSecret) == 0 {
		result.AddError("tokenCredentials.clientSecret is required")
	}

	if setting.Type != TokenCredentialClientCredential {
		result.AddError(fmt.Sprintf("tokenCredentials.type should be %s", TokenCredentialClientCredential))
	}

	return result
}
