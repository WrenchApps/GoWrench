package application_settings

import (
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/validation"
)

func apiEndpointsValidation() validation.ValidateResult {

	var result validation.ValidateResult
	appSettings := ApplicationSettingsStatic

	if appSettings.Api != nil && len(appSettings.Api.Endpoints) > 0 {
		endpoints := appSettings.Api.Endpoints

		for _, endpoint := range endpoints {
			action, err := appSettings.GetActionById(endpoint.ActionID)
			if err != nil {
				result.AddError(err.Error())
			}

			if endpoint.IsProxy && action.Type != action_settings.ActionTypeHttpRequest {
				result.AddError("When endpoint is Proxy the action type should be httpRequest")
			}
		}
	}

	return result
}
