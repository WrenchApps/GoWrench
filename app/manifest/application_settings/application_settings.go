package application_settings

import (
	"errors"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/service_settings"
	credential "wrench/app/manifest/token_credential"
	"wrench/app/manifest/validation"
)

var ApplicationSettingsStatic *ApplicationSettings

type ApplicationSettings struct {
	Api              *api_settings.ApiSettings            `yaml:"api"`
	Actions          []action_settings.ActionSettings     `yaml:"actions"`
	Service          *service_settings.ServiceSettings    `yaml:"service"`
	TokenCredentials []*credential.TokenCredentialSetting `yaml:"tokenCredentials"`
}

func (settings ApplicationSettings) GetActionById(actionId string) (*action_settings.ActionSettings, error) {
	for _, action := range settings.Actions {
		if action.Id == actionId {
			return &action, nil
		}
	}

	return nil, errors.New("action not found")
}

func (settings ApplicationSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if settings.Api != nil {
		result.AppendValidable(settings.Api)
	}

	if settings.Actions != nil {
		for _, validable := range settings.Actions {
			result.AppendValidable(validable)
		}
	}

	if settings.Api != nil {
		result.AppendValidable(settings.Api)
	}

	if settings.TokenCredentials != nil {
		for _, validable := range settings.TokenCredentials {
			result.AppendValidable(validable)
		}
	}

	return result
}
