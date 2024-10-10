package application_settings

import (
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/service_settings"
	"wrench/app/manifest/validation"
)

type ApplicationSetting struct {
	Api     *api_settings.ApiSetting         `yaml:"api"`
	Actions []action_settings.ActionSetting  `yaml:"actions"`
	Service *service_settings.ServiceSetting `yaml:"service"`
}

func (setting ApplicationSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.Api != nil {
		result.AppendValidable(setting.Api)
	}

	if setting.Actions != nil {
		for _, validable := range setting.Actions {
			result.AppendValidable(validable)
		}
	}

	if setting.Api != nil {
		result.AppendValidable(setting.Api)
	}

	return result
}
