package http_settings

import (
	"fmt"
	"strings"
	"wrench/app/manifest/types"
	"wrench/app/manifest/validation"
)

type HttpRequestSetting struct {
	Method            types.HttpMethod  `yaml:"method"`
	Url               string            `yaml:"url"`
	MapFixedHeaders   map[string]string `yaml:"mapFixedHeaders"`
	MapRequestHeaders []string          `yaml:"mapRequestHeaders"`
	TokenCredentialId string            `yaml:"tokenCredentialId"`
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

	if setting.MapFixedHeaders != nil {
		for _, mapHeader := range setting.MapFixedHeaders {
			mapSplitted := strings.Split(mapHeader, ":")
			if len(mapSplitted) != 2 {
				result.AddError("actions.http.request.mapFixedHeaders invalid")
			}
			if len(mapSplitted[0]) == 0 {
				result.AddError("actions.http.request.mapFixedHeaders header key is required")
			}
		}
	}

	if setting.MapRequestHeaders != nil {
		for _, mapHeader := range setting.MapRequestHeaders {
			mapSplitted := strings.Split(mapHeader, ":")
			if len(mapSplitted) > 2 {
				result.AddError("actions.http.request.mapRequestHeaders should contains only one splitter ':'")
			}

			if len(mapHeader) == 0 {
				result.AddError("actions.http.request.mapRequestHeaders itens can't contains empty values")
			}
		}
	}

	// TODO to validate the TokenCredentialId is causing "import cycle not allowed", this I will see later
	// if len(setting.TokenCredentialId) > 0 {
	// 	appSetting := application_settings.ApplicationSettingsStatic
	// 	credentials := appSetting.TokenCredentials

	// 	if credentials == nil {
	// 		result.AddError(fmt.Sprintf("actions.http.request.tokenCredentialId %v don't exist in tokenCredentials", setting.TokenCredentialId))
	// 	} else {
	// 		var hasTokenCredentialId bool = false
	// 		for _, credential := range credentials {
	// 			if setting.TokenCredentialId == credential.Id {
	// 				hasTokenCredentialId = true
	// 			}

	// 			if hasTokenCredentialId {
	// 				break
	// 			}
	// 		}

	// 		if hasTokenCredentialId == false {
	// 			result.AddError(fmt.Sprintf("actions.http.request.tokenCredentialId %v don't exist in tokenCredentials", setting.TokenCredentialId))
	// 		}
	// 	}
	// }

	return result
}
