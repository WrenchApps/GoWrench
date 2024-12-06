package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type HttpRequestClientMockHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}

func (httpRequestClientMockHandler *HttpRequestClientMockHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if wrenchContext.HasError == false {
		bodyContext.Body = httpRequestClientMockHandler.ActionSettings.Http.Mock.Body
		bodyContext.ContentType = httpRequestClientMockHandler.ActionSettings.Http.Mock.ContentType
		bodyContext.HttpStatusCode = httpRequestClientMockHandler.ActionSettings.Http.Mock.StatusCode
		bodyContext.Headers = httpRequestClientMockHandler.ActionSettings.Http.Mock.Headers
	}

	if httpRequestClientMockHandler.Next != nil {
		httpRequestClientMockHandler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (httpRequestClientMockHandler *HttpRequestClientMockHandler) SetNext(handler Handler) {
	httpRequestClientMockHandler.Next = handler
}
