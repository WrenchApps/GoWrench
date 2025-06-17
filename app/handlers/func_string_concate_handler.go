package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FuncStringConcatenateHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncStringConcatenateHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	if !wrenchContext.HasError {

		if len(handler.ActionSettings.Func.Concate) > 0 {

			var stringConcateResult = ""
			for _, item := range handler.ActionSettings.Func.Concate {
				stringConcateResult += contexts.GetCalculatedValue(item, wrenchContext, bodyContext, handler.ActionSettings)
			}

			bodyContext.SetBodyAction(handler.ActionSettings, []byte(stringConcateResult))
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *FuncStringConcatenateHandler) SetNext(next Handler) {
	handler.Next = next
}
