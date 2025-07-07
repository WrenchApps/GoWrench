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
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/connection_settings"
	"wrench/app/manifest/idemp_settings"
	"wrench/app/manifest_cross_funcs"
	"wrench/app/startup/connections"

	"github.com/go-redsync/redsync/v4"
	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type IdempHandler struct {
	Next             Handler
	EndpointSettings *api_settings.EndpointSettings
	IdempSettings    *idemp_settings.IdempSettings
	RedisSettings    *connection_settings.RedisConnectionSettings
}

type idempBodyContext struct {
	CurrentBodyByteArray []byte
	HttpStatusCode       int
	ContentType          string
	Headers              map[string]string
}

func (handler *IdempHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	var redisKeyData string
	var failed bool
	var mutex *redsync.Mutex
	uClient, _ := connections.GetRedisConnection(handler.IdempSettings.RedisConnectionId)

	spanDisplay := fmt.Sprintf("idemp.[%v]", handler.EndpointSettings.IdempId)
	ctxSpan, span := wrenchContext.GetSpan2(ctx, spanDisplay)
	ctx = ctxSpan
	defer span.End()
	start := time.Now()

	if !wrenchContext.HasError {

		keyValue := contexts.GetCalculatedValue(handler.IdempSettings.Key, wrenchContext, bodyContext, nil)
		valueArray := []byte(fmt.Sprint(keyValue))
		hashValue := cross_funcs.GetHash(handler.EndpointSettings.Route, sha256.New, valueArray)

		redisKeyLock := handler.getRedisKeyLock(handler.EndpointSettings.Route, hashValue)
		redisKeyData = handler.getRedisKeyData(handler.EndpointSettings.Route, hashValue)

		rd := cross_funcs.GetRedsyncInstance(handler.IdempSettings.RedisConnectionId)

		mutex = rd.NewMutex(redisKeyLock,
			redsync.WithTries(5),
			redsync.WithRetryDelay(500*time.Millisecond),
			redsync.WithExpiry(20*time.Second),
		)

		if err := mutex.Lock(); err != nil {
			msg := "the distributed lock block request"
			wrenchContext.SetHasError(span, msg, err)
			bodyContext.HttpStatusCode = 409
			bodyContext.CurrentBodyByteArray = []byte(msg)
		} else {

			val, err := uClient.Get(ctx, redisKeyData).Result()
			if err == redis.Nil {
				// do nothing yet
			} else if err != nil {
				wrenchContext.SetHasError(span, err.Error(), err)
				bodyContext.HttpStatusCode = 500
				failed = true
			} else {
				var idempBody idempBodyContext
				jsonErr := json.Unmarshal([]byte(val), &idempBody)

				if jsonErr != nil {
					wrenchContext.SetHasError(span, "idemp error to parse redis body", jsonErr)
					bodyContext.HttpStatusCode = 500
					failed = true
				} else {

					bodyContext.CurrentBodyByteArray = idempBody.CurrentBodyByteArray
					bodyContext.Headers = idempBody.Headers
					bodyContext.ContentType = idempBody.ContentType
					bodyContext.HttpStatusCode = idempBody.HttpStatusCode

					wrenchContext.SetHasCache()
				}
			}
		}
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}

	if !wrenchContext.HasError &&
		!wrenchContext.HasCache {

		idempBody := idempBodyContext{
			CurrentBodyByteArray: bodyContext.CurrentBodyByteArray,
			Headers:              bodyContext.Headers,
			ContentType:          bodyContext.ContentType,
			HttpStatusCode:       bodyContext.HttpStatusCode,
		}

		ttl := time.Duration(handler.IdempSettings.TtlInSeconds) * time.Second

		redisValue, _ := json.Marshal(idempBody)
		err := uClient.Set(ctx, redisKeyData, string(redisValue), ttl).Err()

		if err != nil {
			failed = true
		}
	}

	if mutex != nil {
		if ok, err := mutex.Unlock(); !ok || err != nil {
			app.LogError2(fmt.Sprintf("could not release lock, redis key %v", redisKeyData), err)
		}
	}

	handler.setTraceSpanAttributes(span, redisKeyData, handler.IdempSettings.Id, handler.IdempSettings.RedisConnectionId)
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

func (handler *IdempHandler) getRedisKeyLock(route string, hashValue string) string {
	service := manifest_cross_funcs.GetService()
	return fmt.Sprintf("%v:%v:%v:lock", service.Name, route, hashValue)
}

func (handler *IdempHandler) getRedisKeyData(route string, hashValue string) string {
	service := manifest_cross_funcs.GetService()
	return fmt.Sprintf("%v:%v:%v:data", service.Name, route, hashValue)
}
