package handlers

import (
	"context"
	"fmt"
	"strings"
	"time"
	"wrench/app"
	client "wrench/app/clients/http"
	"wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/startup/token_credentials"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type HttpRequestClientHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}

func (handler *HttpRequestClientHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError {

		start := time.Now()
		ctx, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		defer span.End()

		request := new(client.HttpClientRequestData)
		request.Body = bodyContext.GetBody(handler.ActionSettings)
		request.Method = handler.getMethod(wrenchContext)
		request.Url = handler.getUrl(wrenchContext)
		request.Insecure = handler.ActionSettings.Http.Request.Insecure
		request.SetHeaderTracestate(ctx)
		request.SetHeaders(contexts.GetCalculatedMap(handler.ActionSettings.Http.Request.Headers, wrenchContext, bodyContext, handler.ActionSettings))

		if len(handler.ActionSettings.Http.Request.TokenCredentialId) > 0 {
			tokenData := token_credentials.GetTokenCredentialById(handler.ActionSettings.Http.Request.TokenCredentialId)
			span.SetAttributes(attribute.String("gowrench.tokenCredentials.id", handler.ActionSettings.Http.Request.TokenCredentialId))
			if tokenData != nil {
				bearerToken := strings.Trim(fmt.Sprintf("%s %s", tokenData.TokenType, tokenData.AccessToken), " ")
				if len(tokenData.HeaderName) == 0 {
					request.SetHeader("Authorization", bearerToken)
				} else {
					request.SetHeader(tokenData.HeaderName, bearerToken)
				}
			}
		}

		response, err := client.HttpClientDo(ctx, request)

		if err != nil {
			wrenchContext.SetHasError(ctx, span, "error to call server client", err)
		} else {
			if response.StatusCode > 399 {
				wrenchContext.SetHasError(ctx, span, "server client return one error", err)
			}

			bodyContext.SetBodyAction(handler.ActionSettings, response.Body)

			bodyContext.HttpStatusCode = response.StatusCode
			if handler.ActionSettings.Http.Response != nil {
				bodyContext.SetHeaders(handler.ActionSettings.Http.Response.MapFixedHeaders)
				bodyContext.SetHeaders(mapHttpResponseHeaders(response, handler.ActionSettings.Http.Response.MapResponseHeaders))
			}
		}

		handler.setTraceSpanAttributes(span, response.StatusCode, request.Url, request.Method, request.Insecure)

		duration := time.Since(start).Seconds() * 1000
		handler.metricRecord(ctx, duration, response.StatusCode, request.Url, request.Method)
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *HttpRequestClientHandler) metricRecord(ctx context.Context, duration float64, statusCode int, url string, method string) {
	app.HttpClientDurantion.Record(ctx, duration,
		metric.WithAttributes(
			attribute.Int("http_client_status_code", statusCode),
			attribute.String("http_client_method", method),
			attribute.String("http_client_authority", handler.getAuhorityFromUrl(url)),
		),
	)
}

func (handler *HttpRequestClientHandler) getAuhorityFromUrl(url string) string {
	url = strings.ReplaceAll(url, "http://", "")
	url = strings.ReplaceAll(url, "https://", "")
	urlSplitted := strings.Split(url, "/")
	return urlSplitted[0]
}

func (handler *HttpRequestClientHandler) setTraceSpanAttributes(span trace.Span, statusCode int, url string, method string, insecure bool) {
	span.SetAttributes(
		attribute.Int("http.status_code", statusCode),
		attribute.String("http.url", url),
		attribute.String("http.method", method),
		attribute.Bool("http.insecure", insecure),
	)
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
