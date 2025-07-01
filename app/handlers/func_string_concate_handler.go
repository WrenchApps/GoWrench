package handlers

import (
	"context"
	"fmt"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FuncStringConcatenateHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncStringConcatenateHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	ctx, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
	defer span.End()

	if !wrenchContext.HasError {

		if len(handler.ActionSettings.Func.Concatenate) > 0 {

			var stringConcateResult = ""
			for _, item := range handler.ActionSettings.Func.Concatenate {
				stringConcateResult += fmt.Sprint(contexts.GetCalculatedValue(item, wrenchContext, bodyContext, handler.ActionSettings))
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
