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

			if len(endpoint.ActionID) > 0 {
				action, err := appSettings.GetActionById(endpoint.ActionID)
				if err != nil {
					result.AddError(err.Error())
				}

				if endpoint.IsProxy && action.Type != action_settings.ActionTypeHttpRequest {
					result.AddError("When endpoint is Proxy the action type should be httpRequest")
				}
			}

			if len(endpoint.FlowActionID) > 0 {
				for _, actionId := range endpoint.FlowActionID {

					_, err := appSettings.GetActionById(actionId)
					if err != nil {
						result.AddError(err.Error())
					}
				}

				if endpoint.IsProxy {
					result.AddError("When endpoint is Proxy can be configured flowActionId should use actionId")
				}
			}
		}
	}

	return result
}
