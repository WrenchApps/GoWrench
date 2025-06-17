package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FuncGeneralHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncGeneralHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	if !wrenchContext.HasError {
		result := contexts.GetCalculatedValue(handler.ActionSettings.Func.Command, wrenchContext, bodyContext, handler.ActionSettings)
		bodyContext.SetBodyAction(handler.ActionSettings, []byte(result))
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *FuncGeneralHandler) SetNext(next Handler) {
	handler.Next = next
}
