package contexts

import (
	"context"
	"fmt"
	"net/http"
	settings "wrench/app/manifest/action_settings"
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

func (wrenchContext *WrenchContext) GetSpan(ctx context.Context, action settings.ActionSettings) (context.Context, trace.Span) {
	traceSpanDisplay := fmt.Sprintf("actions[%v].[%v]", action.Id, action.Type)
	return wrenchContext.Tracer.Start(ctx, traceSpanDisplay)
}

func (wrenchContext *WrenchContext) GetSpan2(ctx context.Context, spanDisplay string) (context.Context, trace.Span) {
	return wrenchContext.Tracer.Start(ctx, spanDisplay)
}
