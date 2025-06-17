package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FuncStringConcateHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncStringConcateHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	if !wrenchContext.HasError {

		if len(handler.ActionSettings.Func.Concate) > 0 {

			var strincConcateResult = ""
			for _, item := range handler.ActionSettings.Func.Concate {
				strincConcateResult += contexts.GetCalculatedValue(item, wrenchContext, bodyContext, handler.ActionSettings)
			}

			bodyContext.SetBodyAction(handler.ActionSettings, []byte(strincConcateResult))
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *FuncStringConcateHandler) SetNext(next Handler) {
	handler.Next = next
}
