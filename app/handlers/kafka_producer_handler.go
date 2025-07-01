package handlers

import (
	"context"
	"errors"
	"fmt"
	"log"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/startup/connections"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/trace"
)

type KafkaProducerHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}

func (handler *KafkaProducerHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	ctx, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
	defer span.End()

	if !wrenchContext.HasError {
		settings := handler.ActionSettings

		writer, err := connections.GetKafkaWrite(settings.Kafka.ConnectionId, settings.Kafka.TopicName)

		if err != nil {
			setError("error to get kafka connection id", span, wrenchContext, bodyContext, settings)
		} else {

			value := bodyContext.GetBody(settings)

			var key []byte
			if len(settings.Kafka.MessageKey) > 0 {
				keyValue := fmt.Sprint(contexts.GetCalculatedValue(settings.Kafka.MessageKey, wrenchContext, bodyContext, settings))
				key = []byte(keyValue)
			}

			headers := getKafkaMessageHeaders(settings.Kafka.Headers, wrenchContext, bodyContext, settings)

			err := writer.WriteMessages(context.Background(), kafka.Message{
				Key:     key,
				Value:   value,
				Headers: headers,
			})

			if err != nil {
				msg := fmt.Sprintf("error when will produce message to the topic %v error %v", writer.Topic, err)
				setError(msg, span, wrenchContext, bodyContext, settings)
			} else {

				bodyContext.HttpStatusCode = 200
				bodyContext.ContentType = "text/plain"
				bodyContext.SetBodyAction(settings, []byte(""))
			}
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func setError(msg string, span trace.Span, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, actionSettings *settings.ActionSettings) {
	log.Print(msg)
	bodyContext.HttpStatusCode = 500
	bodyContext.SetBodyAction(actionSettings, []byte(msg))
	bodyContext.ContentType = "text/plain"
	err := errors.New(msg)
	wrenchContext.SetHasError(span, err)
}

func getKafkaMessageHeaders(headersMap map[string]string, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, actionSettings *settings.ActionSettings) []kafka.Header {
	if len(headersMap) > 0 {
		headersCalculated := contexts.GetCalculatedMap(headersMap, wrenchContext, bodyContext, actionSettings)
		headersCalculatedLen := len(headersCalculated)
		if headersCalculatedLen > 0 {
			var headers []kafka.Header

			for key, value := range headersCalculated {
				header := kafka.Header{
					Key:   key,
					Value: []byte(fmt.Sprint(value)),
				}

				headers = append(headers, header)
			}

			return headers
		}
	}

	return nil
}

func (handler *KafkaProducerHandler) SetNext(next Handler) {
	handler.Next = next
}
