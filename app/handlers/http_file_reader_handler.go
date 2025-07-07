package handlers

import (
	"context"
	"fmt"
	"os"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FileReaderHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}

func (handler *FileReaderHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError &&
		!wrenchContext.HasCache {
		ctxSpan, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		ctx = ctxSpan
		defer span.End()

		data, err := os.ReadFile(handler.ActionSettings.File.Path)

		if err != nil {
			msg := fmt.Sprintf("Couldn't read the file %v. Here's why: %v", handler.ActionSettings.File.Path, err)
			bodyContext.HttpStatusCode = 500
			bodyContext.SetBody([]byte(msg))
			bodyContext.ContentType = "text/plain"
			wrenchContext.SetHasError(span, msg, err)
		} else {
			bodyContext.SetBody([]byte(data))
			if handler.ActionSettings.File.Response != nil {
				bodyContext.ContentType = handler.ActionSettings.File.Response.ContentType
				bodyContext.HttpStatusCode = handler.ActionSettings.File.Response.StatusCode
				bodyContext.Headers = handler.ActionSettings.File.Response.Headers
			}
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handlerMock *FileReaderHandler) SetNext(handler Handler) {
	handlerMock.Next = handler
}
