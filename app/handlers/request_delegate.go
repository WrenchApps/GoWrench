package handlers

import (
	"context"
	"net/http"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/api_settings"
)

type RequestDelegate struct {
	Endpoint *settings.EndpointSettings
}

func (request *RequestDelegate) HttpHandler(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()

	var chain = ChainStatic.GetStatic()
	var handler = chain.GetHandler(request.Endpoint.Route)

	bodyContext := new(contexts.BodyContext)
	wrenchContext := new(contexts.WrenchContext)
	wrenchContext.Endpoint = request.Endpoint
	wrenchContext.ResponseWriter = &w
	wrenchContext.Request = r
	handler.Do(ctx, wrenchContext, bodyContext)
}

func (request *RequestDelegate) SetEndpoint(endpoint *settings.EndpointSettings) {
	request.Endpoint = endpoint
}
