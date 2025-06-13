package handlers

import (
	"context"
	"fmt"
	"strings"
	client "wrench/app/clients/http"
	"wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/startup/token_credentials"
)

type HttpRequestClientHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}

func (handler *HttpRequestClientHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError {

		request := new(client.HttpClientRequestData)
		request.Body = bodyContext.BodyByteArray
		request.Method = handler.getMethod(wrenchContext)
		request.Url = handler.getUrl(wrenchContext)
		request.Insecure = handler.ActionSettings.Http.Request.Insecure
		request.SetHeaders(mapHttpRequestHeaders(handler.ActionSettings.Http.Request.Headers, wrenchContext, bodyContext))

		if len(handler.ActionSettings.Http.Request.TokenCredentialId) > 0 {
			tokenData := token_credentials.GetTokenCredentialById(handler.ActionSettings.Http.Request.TokenCredentialId)
			if tokenData != nil {
				bearerToken := fmt.Sprintf("%s %s", tokenData.TokenType, tokenData.AccessToken)
				request.SetHeader("Authorization", bearerToken)
			}
		}

		response, err := client.HttpClientDo(ctx, request)

		if err != nil {
			wrenchContext.SetHasError()
		}

		if response.StatusCode > 399 {
			wrenchContext.SetHasError()
		}

		bodyContext.BodyByteArray = response.Body
		bodyContext.HttpStatusCode = response.StatusCode
		if handler.ActionSettings.Http.Response != nil {
			bodyContext.SetHeaders(handler.ActionSettings.Http.Response.MapFixedHeaders)
			bodyContext.SetHeaders(mapHttpResponseHeaders(response, handler.ActionSettings.Http.Response.MapResponseHeaders))
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *HttpRequestClientHandler) SetNext(next Handler) {
	handler.Next = next
}

func (handler *HttpRequestClientHandler) getMethod(wrenchContext *contexts.WrenchContext) string {

	if !wrenchContext.Endpoint.IsProxy {
		return string(handler.ActionSettings.Http.Request.Method)
	} else {
		return wrenchContext.Request.Method
	}
}

func (handler *HttpRequestClientHandler) getUrl(wrenchContext *contexts.WrenchContext) string {

	if !wrenchContext.Endpoint.IsProxy {
		return handler.ActionSettings.Http.Request.Url
	} else {
		prefix := wrenchContext.Endpoint.Route
		routeTriggered := wrenchContext.Request.RequestURI

		routeWithoutPrefix := strings.ReplaceAll(routeTriggered, prefix, "")
		return handler.ActionSettings.Http.Request.Url + routeWithoutPrefix
	}
}

func mapHttpRequestHeaders(mapRequestHeader map[string]string, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) map[string]string {
	if mapRequestHeader == nil {
		return nil
	}
	mapRequestHeaderResult := make(map[string]string)

	for key, value := range mapRequestHeader {
		var finalValue string
		if contexts.IsCalculatedValue(value) {
			finalValue = contexts.GetCalculatedValue(value, wrenchContext, bodyContext)
		} else {
			finalValue = value
		}

		mapRequestHeaderResult[key] = finalValue
	}

	return mapRequestHeaderResult
}

func mapHttpResponseHeaders(response *client.HttpClientResponseData, mapResponseHeader []string) map[string]string {

	if mapResponseHeader == nil {
		return nil
	}
	mapResponseHeaderResult := make(map[string]string)

	for _, mapHeader := range mapResponseHeader {
		mapSplitted := strings.Split(mapHeader, ":")
		sourceKey := mapSplitted[0]
		var destinationKey string
		if len(mapSplitted) > 1 {
			destinationKey = mapSplitted[1]
		}

		if len(destinationKey) == 0 {
			destinationKey = sourceKey
		}

		headerValue := response.HttpClientResponse.Header.Get(sourceKey)
		mapResponseHeaderResult[destinationKey] = headerValue
	}

	return mapResponseHeaderResult
}
