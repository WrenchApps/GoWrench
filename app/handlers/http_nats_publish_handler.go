package handlers

import (
	"context"
	"time"
	"wrench/app"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/startup/connections"

	"github.com/nats-io/nats.go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

type NatsPublishHandler struct {
	ActionSettings *settings.ActionSettings
	Next           Handler
}

func (handler *NatsPublishHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError {
		start := time.Now()

		ctx, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		defer span.End()

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
				wrenchContext.SetHasError(span, err)
				bodyContext.HttpStatusCode = 500
				bodyContext.SetBody([]byte(err.Error()))
			} else {
				bodyContext.HttpStatusCode = 204
				bodyContext.SetBody([]byte(""))
			}
		}

		duration := time.Since(start).Seconds() * 1000
		handler.metricRecord(ctx, duration, settings.Nats.ConnectionId, settings.Nats.SubjectName)
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *NatsPublishHandler) metricRecord(ctx context.Context, duration float64, connectionId string, subjectName string) {
	app.NatsPublishDurtation.Record(ctx, duration,
		metric.WithAttributes(
			attribute.String("gowrench_connections_id", connectionId),
			attribute.String("nats_publish_subject_name", subjectName),
		),
	)
}

func (handler *NatsPublishHandler) SetNext(next Handler) {
	handler.Next = next
}
