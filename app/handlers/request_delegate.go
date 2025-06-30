package handlers

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"wrench/app"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/api_settings"

	"go.opentelemetry.io/otel/trace"
)

type RequestDelegate struct {
	Endpoint *settings.EndpointSettings
}

func (request *RequestDelegate) HttpHandler(w http.ResponseWriter, r *http.Request) {
	traceIdArray := r.Header["Tracestate"]
	traceId := ""

	if len(traceIdArray) > 0 {
		traceId = traceIdArray[0]
	}

	var ctx context.Context
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

		ctx = trace.ContextWithSpanContext(context.Background(), parent)

	} else {
		ctx = context.Background()
	}

	var chain = ChainStatic.GetStatic()
	var handler = chain.GetHandler(request.Endpoint.Route)

	bodyContext := new(contexts.BodyContext)
	wrenchContext := new(contexts.WrenchContext)
	wrenchContext.Endpoint = request.Endpoint
	wrenchContext.ResponseWriter = &w
	wrenchContext.Request = r
	traceDisplay := fmt.Sprintf("Api http %v %v", request.Endpoint.Method, request.Endpoint.Route)

	wrenchContext.Tracer = app.Tracer

	ctx, span := wrenchContext.Tracer.Start(ctx, traceDisplay)

	handler.Do(ctx, wrenchContext, bodyContext)

	defer span.End()
}

func (request *RequestDelegate) SetEndpoint(endpoint *settings.EndpointSettings) {
	request.Endpoint = endpoint
}
