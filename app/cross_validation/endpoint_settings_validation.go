package cross_validation

import (
	"fmt"
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/validation"
	"wrench/app/manifest_cross_funcs"
)

func endpointSettingsCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	endpoints := appSetting.Api.Endpoints

	if len(endpoints) > 0 {
		for _, endpoint := range endpoints {
			if len(endpoint.IdempId) > 0 {
				_, err := manifest_cross_funcs.GetIdempSettingById(endpoint.IdempId)

				if err != nil {
					result.AddError(fmt.Sprintf("api.endpoints[%v].idempId %v don't exist in idemps", endpoint.Route, endpoint.IdempId))
				}
			}
			if len(endpoint.Roles) > 0 ||
				len(endpoint.Scopes) > 0 ||
				len(endpoint.Claims) > 0 {
				authSetting := manifest_cross_funcs.GetAuthorizationSettings()

				if authSetting != nil && authSetting.Type == api_settings.HMACAuthorizationType {
					result.AddError(fmt.Sprintf("api.endpoints[%v] is using roles/scopes/claim which is not allowed for HMAC authorization", endpoint.Route))
				}
			}
		}
	}

	return result
}
