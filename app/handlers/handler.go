package handlers

import (
	"context"
	contexts "wrench/app/contexts"
)

type Handler interface {
	Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext)
	SetNext(Handler)
}
