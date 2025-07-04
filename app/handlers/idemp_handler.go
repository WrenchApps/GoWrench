package handlers

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"
	"wrench/app"
	contexts "wrench/app/contexts"
	"wrench/app/cross_funcs"
	settings "wrench/app/manifest/action_settings"
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/idemp_settings"
	"wrench/app/manifest_cross_funcs"
	"wrench/app/startup/connections"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type IdempHandler struct {
	Next             Handler
	ActionSettings   *settings.ActionSettings
	EndpointSettings *api_settings.EndpointSettings
}

type idempBodyContext struct {
	CurrentBodyByteArray []byte
	HttpStatusCode       int
	ContentType          string
	Headers              map[string]string
}

func (handler *IdempHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	var redisKey string
	var redisClient *redis.Client
	var idemp *idemp_settings.IdempSettings
	var failed bool

	spanDisplay := fmt.Sprintf("idemp.[%v]", handler.EndpointSettings.IdempId)
	ctxSpan, span := wrenchContext.GetSpan2(ctx, spanDisplay)
	ctx = ctxSpan
	defer span.End()
	start := time.Now()

	if !wrenchContext.HasError {

		idempotency, err := manifest_cross_funcs.GetIdempSettingById(handler.EndpointSettings.IdempId)
		idemp = idempotency
		if err != nil {
			wrenchContext.SetHasError(span, err.Error(), err)
			failed = true
		} else {
			keyValue := contexts.GetCalculatedValue(idemp.Key, wrenchContext, bodyContext, handler.ActionSettings)
			redisClient, err = connections.GetRedisConnection(idemp.RedisConnectionId)

			if err != nil {
				wrenchContext.SetHasError(span, err.Error(), err)
				failed = true
			}

			valueArray := []byte(fmt.Sprint(keyValue))
			hashValue := cross_funcs.GetHash(handler.EndpointSettings.Route, sha256.New, valueArray)
			redisKey = handler.getRedisKey(handler.EndpointSettings.Route, hashValue)

			val, err := redisClient.Get(ctx, redisKey).Result()

			if err == redis.Nil {
				// do nothing yet
			} else if err != nil {
				wrenchContext.SetHasError(span, err.Error(), err)
				failed = true
			} else {
				var idempBody idempBodyContext
				jsonErr := json.Unmarshal([]byte(val), &idempBody)

				if jsonErr != nil {
					wrenchContext.SetHasError(span, jsonErr.Error(), jsonErr)
					failed = true
				}

				bodyContext.CurrentBodyByteArray = idempBody.CurrentBodyByteArray
				bodyContext.Headers = idempBody.Headers
				bodyContext.ContentType = idempBody.ContentType
				bodyContext.HttpStatusCode = idempBody.HttpStatusCode

				wrenchContext.SetHasCache()
			}
		}

	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}

	if !wrenchContext.HasError {

		idempBody := idempBodyContext{
			CurrentBodyByteArray: bodyContext.CurrentBodyByteArray,
			Headers:              bodyContext.Headers,
			ContentType:          bodyContext.ContentType,
			HttpStatusCode:       bodyContext.HttpStatusCode,
		}

		ttl := time.Duration(idemp.TtlInSeconds) * time.Second
		err := redisClient.Set(ctx, redisKey, idempBody, ttl).Err()
		if err != nil {
			failed = true
		}
	}

	handler.setTraceSpanAttributes(span, redisKey, idemp.Id, idemp.RedisConnectionId)
	duration := time.Since(start).Seconds() * 1000
	handler.metricRecord(ctx, duration, failed)
}

func (handler *IdempHandler) SetNext(next Handler) {
	handler.Next = next
}

func (handler *IdempHandler) setTraceSpanAttributes(span trace.Span, key string, idempId string, redisConnectionId string) {
	span.SetAttributes(
		attribute.String("idemp_key", key),
		attribute.String("idemp_id", idempId),
		attribute.String("idemp_redis_connection_id", redisConnectionId),
	)
}

func (handler *IdempHandler) metricRecord(ctx context.Context, duration float64, failed bool) {
	app.IdempDuration.Record(ctx, duration,
		metric.WithAttributes(
			attribute.Bool("failed", failed),
		),
	)
}

func (handler *IdempHandler) getRedisKey(route string, hashValue string) string {
	service := manifest_cross_funcs.GetService()
	return fmt.Sprintf("%v:%v:%v", service.Name, route, hashValue)
}
