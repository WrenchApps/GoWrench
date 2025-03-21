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

func (handler *HttpRequestClientMockHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError {
		if !handler.ActionSettings.Http.Mock.MirrorBody {
			bodyContext.BodyByteArray = []byte(handler.ActionSettings.Http.Mock.Body)
		}
		bodyContext.ContentType = handler.ActionSettings.Http.Mock.ContentType
		bodyContext.HttpStatusCode = handler.ActionSettings.Http.Mock.StatusCode
		bodyContext.Headers = handler.ActionSettings.Http.Mock.Headers
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handlerMock *HttpRequestClientMockHandler) SetNext(handler Handler) {
	handlerMock.Next = handler
}
