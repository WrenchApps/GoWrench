package handlers

import (
	"context"
	contexts "wrench/app/contexts"
)

type HttpFirstHandler struct {
	Next Handler
}

func (httpFirst *HttpFirstHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	if httpFirst.Next != nil {
		httpFirst.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (httpFirst *HttpFirstHandler) SetNext(next Handler) {
	httpFirst.Next = next
}
