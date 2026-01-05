package cross_validation

import (
	"fmt"
	"wrench/app/manifest/application_settings"
	credential "wrench/app/manifest/token_credential_settings"
	"wrench/app/manifest/validation"
	"wrench/app/startup/tls_load"
)

func tokenCredentialsCrossValidation(appSetting *application_settings.ApplicationSettings) validation.ValidateResult {
	var result validation.ValidateResult

	result.Append(tokenCredentialsIdDuplicated(appSetting.TokenCredentials))
	result.Append(tokenCredentialsTlsId(appSetting.TokenCredentials))

	return result
}

func tokenCredentialsIdDuplicated(settings []*credential.TokenCredentialSetting) validation.ValidateResult {

	var result validation.ValidateResult

	if len(settings) > 0 {
		hasIds := toHasIdSlice(settings)
		duplicateIds := duplicateIdsValid(hasIds)

		for _, id := range duplicateIds {
			result.AddError(fmt.Sprintf("tokenCredentials.id %v duplicated", id))
		}
	}

	return result
}

func tokenCredentialsTlsId(tokenCredentials []*credential.TokenCredentialSetting) validation.ValidateResult {
	var result validation.ValidateResult

	if len(tokenCredentials) > 0 {
		for _, tokenCredential := range tokenCredentials {
			if len(tokenCredential.TlsId) > 0 {
				_, err := tls_load.GetTlsConfigById(tokenCredential.TlsId)

				if err != nil {
					result.AddError(fmt.Sprintf("tokenCredentials[%v].tlsId %v don't exist in tls", tokenCredential.Id, tokenCredential.TlsId))
				}
			}
		}
	}

	return result
}
