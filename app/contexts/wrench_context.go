package contexts

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"wrench/app"
	settings "wrench/app/manifest/action_settings"
	api_settings "wrench/app/manifest/api_settings"

	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type WrenchContext struct {
	ResponseWriter *http.ResponseWriter
	Request        *http.Request
	HasError       bool
	HasCache       bool
	Endpoint       *api_settings.EndpointSettings
	Tracer         trace.Tracer
	Meter          metric.Meter
}

func (wrenchContext *WrenchContext) SetHasError(span trace.Span, msg string, err error) {
	span.RecordError(err)

	errorDescription := msg
	if err != nil {
		errorDescription = err.Error()
	}

	span.SetStatus(codes.Error, errorDescription)

	app.LogError(app.WrenchErrorLog{Message: msg, Error: err})
	wrenchContext.HasError = true
}

func (wrenchContext *WrenchContext) SetHasCache() {
	wrenchContext.HasCache = true
}

func (wrenchContext *WrenchContext) SetHasError2() {
	wrenchContext.HasError = true
}

func (wrenchContext *WrenchContext) GetSpan(ctx context.Context, action settings.ActionSettings) (context.Context, trace.Span) {
	traceSpanDisplay := fmt.Sprintf("actions[%v].[%v]", action.Id, action.Type)
	return wrenchContext.Tracer.Start(ctx, traceSpanDisplay)
}

func (wrenchContext *WrenchContext) GetSpan2(ctx context.Context, spanDisplay string) (context.Context, trace.Span) {
	return wrenchContext.Tracer.Start(ctx, spanDisplay)
}

func (wrenchContext *WrenchContext) GetContext(traceId string) context.Context {
	if len(traceId) > 0 {

		traceIdSpllited := strings.Split(traceId, "-")

		traceID, _ := trace.TraceIDFromHex(traceIdSpllited[1])
		spanID, _ := trace.SpanIDFromHex(traceIdSpllited[2])

		parent := trace.NewSpanContext(trace.SpanContextConfig{
			TraceID:    traceID,
			SpanID:     spanID,
			TraceFlags: trace.FlagsSampled,
			Remote:     true,
		})

		return trace.ContextWithSpanContext(context.Background(), parent)

	} else {
		return context.Background()
	}
}
