package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FuncHashHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncHashHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	if !wrenchContext.HasError {
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}

}

func (handler *FuncHashHandler) SetNext(next Handler) {
	handler.Next = next
}
