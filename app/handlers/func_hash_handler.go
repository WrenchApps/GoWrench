package handlers

import (
	"context"
	"fmt"
	contexts "wrench/app/contexts"
	"wrench/app/cross_funcs"
	settings "wrench/app/manifest/action_settings"
)

type FuncHashHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *FuncHashHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError &&
		!wrenchContext.HasCache {
		ctxSpan, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		ctx = ctxSpan
		defer span.End()

		key := contexts.GetCalculatedValue(handler.ActionSettings.Func.Hash.Key, wrenchContext, bodyContext, handler.ActionSettings)
		hashType := cross_funcs.GetHashFunc(handler.ActionSettings.Func.Hash.Alg)
		currentBody := bodyContext.GetBody(handler.ActionSettings)

		hashValue := cross_funcs.GetHash(fmt.Sprint(key), hashType, currentBody)
		bodyContext.SetBodyAction(handler.ActionSettings, []byte(hashValue))
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}

}

func (handler *FuncHashHandler) SetNext(next Handler) {
	handler.Next = next
}
