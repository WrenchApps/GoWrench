package handlers

import (
	"context"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/startup/connections"

	"github.com/nats-io/nats.go"
)

type NatsPublishHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *NatsPublishHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	ctx, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
	defer span.End()

	if !wrenchContext.HasError {
		settings := handler.ActionSettings

		natsConn := connections.GetNatsConnectionById(settings.Nats.ConnectionId)

		msg := &nats.Msg{
			Subject: settings.Nats.SubjectName,
			Data:    bodyContext.GetBody(settings),
			//Header:  nats.Header{},    // create mapper to add headers in message
		}

		var err error
		if settings.Nats.IsStream {
			js := connections.GetJetStreamByConnectionId(settings.Nats.ConnectionId)
			_, err = js.PublishMsg(msg)

		} else {
			err = natsConn.PublishMsg(msg)
		}

		if settings.ShouldPreserveBody() {
			bodyContext.SetBodyPreserved(settings.Id, []byte(""))
		} else {
			if err != nil {
				wrenchContext.SetHasError()
				bodyContext.HttpStatusCode = 500
				bodyContext.SetBody([]byte(err.Error()))
			} else {
				bodyContext.HttpStatusCode = 204
				bodyContext.SetBody([]byte(""))
			}
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *NatsPublishHandler) SetNext(next Handler) {
	handler.Next = next
}
