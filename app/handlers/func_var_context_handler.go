package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
)

type FuncVarContextHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncVarContextHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	if !wrenchContext.HasError {

		varsConfigured := handler.ActionSettings.Func.Vars

		if len(varsConfigured) > 0 {
			varsResult := contexts.GetCalculatedMap(varsConfigured, wrenchContext, bodyContext, handler.ActionSettings)
			var result, _ = bodyContext.ConvertMapToByteArray(varsResult)

			bodyContext.SetBodyAction(handler.ActionSettings, result)
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}

}

func (handler *FuncVarContextHandler) SetNext(next Handler) {
	handler.Next = next
}
