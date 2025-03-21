package contexts

import (
	"net/http"
	api_settings "wrench/app/manifest/api_settings"
)

type WrenchContext struct {
	ResponseWriter *http.ResponseWriter
	Request        *http.Request
	HasError       bool
	Endpoint       *api_settings.EndpointSettings
}

func (wrenchContext *WrenchContext) SetHasError() {
	wrenchContext.HasError = true
}
