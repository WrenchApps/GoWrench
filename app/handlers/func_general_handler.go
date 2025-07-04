package handlers

import (
	"context"
	"fmt"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FuncGeneralHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncGeneralHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError &&
		!wrenchContext.HasCache {
		ctxSpan, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		ctx = ctxSpan
		defer span.End()

		result := contexts.GetCalculatedValue(string(handler.ActionSettings.Func.Command), wrenchContext, bodyContext, handler.ActionSettings)
		bodyContext.SetBodyAction(handler.ActionSettings, []byte(fmt.Sprint(result)))
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *FuncGeneralHandler) SetNext(next Handler) {
	handler.Next = next
}
