package cross_validation

import (
	"fmt"
	"wrench/app/manifest/action_settings"
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
	"wrench/app/manifest_cross_funcs"
	"wrench/app/startup/tls_load"
)

func httpRequestCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	actionHttpRequest := getActionsByType(appSetting.Actions, action_settings.ActionTypeHttpRequest)

	result.Append(actionHttpRequestTokenCredential(actionHttpRequest))
	result.Append(actionHttpRequestTlsId(actionHttpRequest))
	result.Append(actionHttpRequestMethod(actionHttpRequest, appSetting.Api))

	return result
}

func actionHttpRequestTokenCredential(actions []*action_settings.ActionSettings) validation.ValidateResult {
	var result validation.ValidateResult
	if len(actions) > 0 {
		for _, action := range actions {
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

func actionHttpRequestTlsId(actions []*action_settings.ActionSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if len(actions) > 0 {
		for _, action := range actions {
			if action.Http.Request != nil && len(action.Http.Request.TlsId) > 0 {
				_, err := tls_load.GetTlsConfigById(action.Http.Request.TlsId)

				if err != nil {
					result.AddError(fmt.Sprintf("actions[%v].tlsId %v don't exist in tls", action.Id, action.Http.Request.TlsId))
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

func actionHttpRequestMethod(actions []*action_settings.ActionSettings, api *api_settings.ApiSettings) validation.ValidateResult {
	var result validation.ValidateResult
	if len(actions) > 0 {
		for _, action := range actions {
			request := action.Http.Request

			endpoints, err := api.GetEndpointsByActionId(action.Id)

			if len(endpoints) == 0 {
				continue
			}

			isProxyEndpoint := false
			if err == nil {
				for _, endpoint := range endpoints {

					if endpoint.IsProxy {
						isProxyEndpoint = true
						break
					}
				}
			}

			if !isProxyEndpoint {
				if err != nil {
					result.AddError(fmt.Sprintf("actions.http.request.method is required"))
				} else {
					if len(request.Method) == 0 {
						result.AddError("actions.http.request.method is required")
					} else {
						if (request.Method == types.HttpMethodGet ||
							request.Method == types.HttpMethodPost ||
							request.Method == types.HttpMethodPut ||
							request.Method == types.HttpMethodPatch ||
							request.Method == types.HttpMethodDelete) == false {

							result.AddError("actions.http.request.method should contain valid value (get, post, put, patch or delete)")
						}
					}
				}
			}

		}
	}
	return result
}
