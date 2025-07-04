package handlers

import (
	"context"
	"errors"
	"fmt"
	"time"
	"wrench/app"
	contexts "wrench/app/contexts"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/startup/connections"

	"github.com/segmentio/kafka-go"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type KafkaProducerHandler struct {
	Next           Handler
	ActionSettings *settings.ActionSettings
}

func (handler *KafkaProducerHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError {
		start := time.Now()

		ctx, span := wrenchContext.GetSpan(ctx, *handler.ActionSettings)
		defer span.End()
		settings := handler.ActionSettings

		writer, err := connections.GetKafkaWrite(settings.Kafka.ConnectionId, settings.Kafka.TopicName)

		if err != nil {
			setError("error to get kafka connection id", span, wrenchContext, bodyContext, settings)
		} else {

			value := bodyContext.GetBody(settings)

			var key []byte
			var keyValue string
			if len(settings.Kafka.MessageKey) > 0 {
				keyValue = fmt.Sprint(contexts.GetCalculatedValue(settings.Kafka.MessageKey, wrenchContext, bodyContext, settings))
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

			handler.setSpanAttributes(span, settings.Kafka.ConnectionId, settings.Kafka.TopicName, keyValue)

			duration := time.Since(start).Seconds() * 1000
			handler.metricRecord(ctx, duration, settings.Kafka.ConnectionId, settings.Kafka.TopicName)
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *KafkaProducerHandler) metricRecord(ctx context.Context, duration float64, connectionId string, topicName string) {
	app.KafkaProducerDuration.Record(ctx, duration,
		metric.WithAttributes(
			attribute.String("gowrench_connections_id", connectionId),
			attribute.String("kafka_producer_topic_name", topicName),
		),
	)
}

func (handler *KafkaProducerHandler) setSpanAttributes(span trace.Span, connectionId string, topicName string, key string) {
	span.SetAttributes(
		attribute.String("gowrench.connections.id", connectionId),
		attribute.String("kafka.producer.topic_name", topicName),
		attribute.String("kafka.producer.key", key),
	)
}

func setError(msg string, span trace.Span, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext, actionSettings *settings.ActionSettings) {
	bodyContext.HttpStatusCode = 500
	bodyContext.SetBodyAction(actionSettings, []byte(msg))
	bodyContext.ContentType = "text/plain"
	err := errors.New(msg)
	wrenchContext.SetHasError(span, msg, err)
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
