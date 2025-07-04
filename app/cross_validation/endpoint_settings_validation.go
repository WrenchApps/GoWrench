package cross_validation

import (
	"fmt"
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
		}
	}

	return result
}
