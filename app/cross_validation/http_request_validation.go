package cross_validation

import (
	"fmt"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/validation"
	"wrench/app/manifest_cross_funcs"
)

func httpRequestCrossValid(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	actionHttpRequest := getActionsByType(appSetting.Actions, action_settings.ActionTypeHttpRequest)

	if len(actionHttpRequest) > 0 {
		for _, action := range actionHttpRequest {
			// valid if exist tokenCredential
			if action.Http.Request != nil && len(action.Http.Request.TokenCredentialId) > 0 {
				_, err := manifest_cross_funcs.GetTokenCredentialSettingById(action.Http.Request.TokenCredentialId)

				if err != nil {
					result.AddError(fmt.Sprintf("actions.http.request.tokenCredentialId %v don't exist in tokenCredentials", action.Http.Request.TokenCredentialId))
				}
			}
		}
	}

	return result
}

func getActionsByType(actions []*action_settings.ActionSettings, actionType action_settings.ActionType) []*action_settings.ActionSettings {
	var actionsResult []*action_settings.ActionSettings

	if len(actions) > 0 {
		for _, action := range actions {
			if action.Type == actionType {
				actionsResult = append(actionsResult, action)
			}
		}
	}
	return actionsResult
}
