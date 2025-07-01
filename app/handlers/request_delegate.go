package handlers

import (
	"context"
	"fmt"
	"net/http"
	"time"
	"wrench/app"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/api_settings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type RequestDelegate struct {
	Endpoint *settings.EndpointSettings
}

func (request *RequestDelegate) HttpHandler(w http.ResponseWriter, r *http.Request) {

	start := time.Now()
	bodyContext := new(contexts.BodyContext)
	wrenchContext := new(contexts.WrenchContext)

	traceId := getHeader(r, "Tracestate")
	ctx := wrenchContext.GetContext(traceId)

	var chain = ChainStatic.GetStatic()
	var handler = chain.GetHandler(request.Endpoint.Route)

	wrenchContext.Tracer = app.Tracer
	wrenchContext.Meter = app.Meter
	wrenchContext.Endpoint = request.Endpoint
	wrenchContext.ResponseWriter = &w
	wrenchContext.Request = r

	traceDisplay := fmt.Sprintf("Api http %v %v", request.Endpoint.Method, request.Endpoint.Route)
	ctx, span := wrenchContext.GetSpan2(ctx, traceDisplay)
	defer span.End()

	handler.Do(ctx, wrenchContext, bodyContext)

	request.setSpanAttributes(span, request.Endpoint.Route, fmt.Sprint(request.Endpoint.Method), bodyContext.HttpStatusCode)
	duration := time.Since(start).Seconds() * 1000
	request.metricRecord(ctx, duration, request.Endpoint.Route, fmt.Sprint(request.Endpoint.Method), bodyContext.HttpStatusCode)
}

func (handler *RequestDelegate) metricRecord(ctx context.Context, duration float64, route string, method string, statusCode int) {
	app.HttpServerDuration.Record(ctx, duration,
		metric.WithAttributes(
			attribute.String("http_server_method", route),
			attribute.String("http_server_method", method),
			attribute.Int("http_server_status_code", statusCode),
		),
	)
}

func (handler *RequestDelegate) setSpanAttributes(span trace.Span, route string, method string, statusCode int) {
	span.SetAttributes(
		attribute.String("server.route", route),
		attribute.String("server.method", method),
		attribute.Int("server.status_code", statusCode),
	)
}

func (request *RequestDelegate) SetEndpoint(endpoint *settings.EndpointSettings) {
	request.Endpoint = endpoint
}

func getHeader(r *http.Request, headerName string) string {
	traceIdArray := r.Header[headerName]
	traceId := ""

	if len(traceIdArray) > 0 {
		traceId = traceIdArray[0]
	}

	return traceId
}
