package handlers

import (
	"context"
	"fmt"
	"strconv"
	"wrench/app"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type HttpRequestClientMockHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}

func (handler *HttpRequestClientMockHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if (!wrenchContext.HasError || handler.ActionSettings.RunEvenIfFlowHasError) &&
		!wrenchContext.HasCache &&
		!wrenchContext.Unauthorized {

		ctx2, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		defer span.End()
		ctx = ctx2

		if !handler.ActionSettings.Http.Mock.MirrorBody {
			bodyContext.SetBodyAction(handler.ActionSettings, []byte(handler.ActionSettings.Http.Mock.Body))
		}

		if len(handler.ActionSettings.Http.Mock.ContentType) > 0 {
			bodyContext.ContentType = handler.ActionSettings.Http.Mock.ContentType
		}

		httpStatusCode, err := handler.getResponseStatusCode(wrenchContext, bodyContext)
		if err != nil {
			app.LogError(app.WrenchErrorLog{Message: "error to parse http status code", Error: err})
		}
		bodyContext.HttpStatusCode = httpStatusCode
		bodyContext.Headers = handler.ActionSettings.Http.Mock.Headers
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handlerMock *HttpRequestClientMockHandler) getResponseStatusCode(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) (int, error) {
	statusCode := contexts.GetCalculatedValue(handlerMock.ActionSettings.Http.Mock.StatusCode, wrenchContext, bodyContext, handlerMock.ActionSettings)
	statusCodeInt, err := strconv.Atoi(fmt.Sprintf("%v", statusCode))
	if err != nil {
		return 0, err
	}
	return statusCodeInt, nil
}

func (handlerMock *HttpRequestClientMockHandler) SetNext(handler Handler) {
	handlerMock.Next = handler
}
