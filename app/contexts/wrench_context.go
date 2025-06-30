package contexts

import (
	"net/http"
	api_settings "wrench/app/manifest/api_settings"

	"go.opentelemetry.io/otel/trace"
)

type WrenchContext struct {
	ResponseWriter *http.ResponseWriter
	Request        *http.Request
	HasError       bool
	Endpoint       *api_settings.EndpointSettings
	Tracer         trace.Tracer
}

func (wrenchContext *WrenchContext) SetHasError() {
	wrenchContext.HasError = true
}
