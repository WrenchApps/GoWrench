package handlers

import (
	"fmt"
	"net/http"
	"wrench/app"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/api_settings"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RequestDelegate struct {
	Endpoint *settings.EndpointSettings
}

func (request *RequestDelegate) HttpHandler(w http.ResponseWriter, r *http.Request) {

	bodyContext := new(contexts.BodyContext)
	wrenchContext := new(contexts.WrenchContext)

	traceId := getHeader(r, "Tracestate")
	ctx := wrenchContext.GetContext(traceId)

	var chain = ChainStatic.GetStatic()
	var handler = chain.GetHandler(request.Endpoint.Route)

	wrenchContext.Tracer = app.Tracer
	wrenchContext.Endpoint = request.Endpoint
	wrenchContext.ResponseWriter = &w
	wrenchContext.Request = r

	traceDisplay := fmt.Sprintf("Api http %v %v", request.Endpoint.Method, request.Endpoint.Route)
	ctx, span := wrenchContext.GetSpan2(ctx, traceDisplay)

	handler.Do(ctx, wrenchContext, bodyContext)

	request.setSpanAttributes(span, request.Endpoint.Route, fmt.Sprint(request.Endpoint.Method), bodyContext.HttpStatusCode)
	defer span.End()
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
