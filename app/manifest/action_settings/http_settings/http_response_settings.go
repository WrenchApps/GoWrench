package http_settings

import (
	"strings"
	"wrench/app/manifest/validation"
)

type HttpResponseSettings struct {
	MapFixedHeaders    map[string]string `yaml:"mapFixedHeaders"`
	MapResponseHeaders []string          `yaml:"mapResponseHeaders"`
}

func (setting HttpResponseSettings) Valid() validation.ValidateResult {
	var result validation.ValidateResult

	if setting.MapFixedHeaders != nil {
		for _, mapHeader := range setting.MapFixedHeaders {
			mapSplitted := strings.Split(mapHeader, ":")
			if len(mapSplitted) != 2 {
				result.AddError("actions.http.response.mapFixedHeaders invalid")
			}
			if len(mapSplitted[0]) == 0 {
				result.AddError("actions.http.response.mapFixedHeaders header key is required")
			}
		}
	}

	if setting.MapResponseHeaders != nil {
		for _, mapHeader := range setting.MapResponseHeaders {
			mapSplitted := strings.Split(mapHeader, ":")
			if len(mapSplitted) > 2 {
				result.AddError("actions.http.response.mapResponseHeaders should contains only one splitter ':'")
			}

			if len(mapHeader) == 0 {
				result.AddError("actions.http.response.mapResponseHeaders itens can't contains empty values")
			}
		}
	}

	return result
}
