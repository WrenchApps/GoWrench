package http_settings

import (
	"fmt"

	//"wrench/app/cross_cutting"
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type HttpRequestSetting struct {
	Method            types.HttpMethod  `yaml:"method"`
	Url               string            `yaml:"url"`
	Headers           map[string]string `yaml:"headers"`
	TokenCredentialId string            `yaml:"tokenCredentialId"`
	Insecure          bool              `yaml:"insecure"`
}

func (setting HttpRequestSetting) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if len(setting.Method) == 0 {
		var msg = fmt.Sprintf("actions.http.request.method is required")
		result.AddError(msg)
	} else {
		if (setting.Method == types.HttpMethodGet ||
			setting.Method == types.HttpMethodPost ||
			setting.Method == types.HttpMethodPut ||
			setting.Method == types.HttpMethodPatch ||
			setting.Method == types.HttpMethodDelete) == false {

			var msg = fmt.Sprintf("actions.http.request.method should contain valid value (get, post, put, patch or delete)")
			result.AddError(msg)
		}
	}

	if len(setting.Url) == 0 {
		result.AddError("actions.http.request.url is required")
	}

	// if setting.Headers != nil {
	// 	for _, mapHeader := range setting.Headers {
	// 		mapSplitted := strings.Split(mapHeader, ":")
	// 		if len(mapSplitted) != 2 {
	// 			result.AddError("actions.http.request.headers invalid")
	// 		}
	// 		if len(mapSplitted[0]) == 0 {
	// 			result.AddError("actions.http.request.headers header key is required")
	// 		}
	// 	}
	// }

	// if len(setting.TokenCredentialId) > 0 {
	// 	tokenCredential := cross_cutting.GetTokenCredentialById(setting.TokenCredentialId)

	// 	if tokenCredential == nil {
	// 		result.AddError(fmt.Sprintf("actions.http.request.tokenCredentialId %v don't exist in tokenCredentials", setting.TokenCredentialId))
	// 	}
	// }

	return result
}
