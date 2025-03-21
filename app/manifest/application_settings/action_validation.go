package application_settings

import (
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/validation"
)

func actionValidation(action action_settings.ActionSettings) validation.ValidateResult {
	var result validation.ValidateResult

	appSettings := ApplicationSettingsStatic

	endpoint, _ := appSettings.GetEndpointByActionId(action.Id)

	if endpoint != nil {
		if endpoint.IsProxy {
			if action.Http != nil && action.Http.Request != nil && len(action.Http.Request.Method) > 0 {
				result.AddError("Actions configured in proxy endpoints  shouldn't configure method")
			}
		}
	}

	return result
}
