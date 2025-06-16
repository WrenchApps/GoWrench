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

		varsConfigured := handler.ActionSettings.Vars
		varsResult := contexts.GetCalculatedMap(varsConfigured, wrenchContext, bodyContext)

		if handler.ActionSettings.PreserveCurrentBody {
			bodyContext.SetMapObjectWithPreservedId(handler.ActionSettings.Id, varsResult)
		} else {
			bodyContext.SetMapObject(varsResult)
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}

}

func (handler *FuncVarContextHandler) SetNext(next Handler) {
	handler.Next = next
}
