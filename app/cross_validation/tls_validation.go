package cross_validation

import (
	"fmt"
	"wrench/app/manifest/application_settings"
	"wrench/app/manifest/tls_settings"
	"wrench/app/manifest/validation"
)

func tlsCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	if len(appSetting.Tls) > 0 {

		result.Append(tlsIdDuplicated(appSetting.Tls))
	}

	return result
}

func tlsIdDuplicated(settings []*tls_settings.TlsSettings) validation.ValidateResult {

	var result validation.ValidateResult

	hasIds := toHasIdSlice(settings)
	duplicateIds := duplicateIdsValid(hasIds)

	for _, id := range duplicateIds {
		result.AddError(fmt.Sprintf("tls.id %v duplicated", id))
	}

	return result
}
