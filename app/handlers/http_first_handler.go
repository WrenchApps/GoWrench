package handlers

import (
	"context"
	"io"
	contexts "wrench/app/contexts"
)

type HttpFirstHandler struct {
	Next Handler
}

func (httpFirst *HttpFirstHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	body, err := io.ReadAll(wrenchContext.Request.Body)
	if err != nil {
		bodyContext.ContentType = "text/plain"
		bodyContext.HttpStatusCode = 500
		bodyContext.SetBody([]byte("Failed to read request body"))
		wrenchContext.SetHasError2()
	}

	bodyContext.SetBody(body)
	bodyContext.ContentType = "application/json"
	bodyContext.HttpStatusCode = 200

	if httpFirst.Next != nil {
		httpFirst.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (httpFirst *HttpFirstHandler) SetNext(next Handler) {
	httpFirst.Next = next
}
