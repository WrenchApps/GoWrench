package handlers

import (
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"net/http"
	"time"
	"wrench/app"
	contexts "wrench/app/contexts"
	"wrench/app/cross_funcs"
	"wrench/app/manifest/api_settings"
	"wrench/app/manifest/rate_limit_settings"
	"wrench/app/startup/connections"

	"github.com/go-redis/redis_rate/v10"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type RateLimitHandler struct {
	Next              Handler
	EndpointSettings  *api_settings.EndpointSettings
	RateLimitSettings *rate_limit_settings.RateLimitSettings
}

func (handler *RateLimitHandler) Do(ctx context.Context, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {

	if !wrenchContext.HasError &&
		!wrenchContext.HasCache {
		start := time.Now()
		rtSettings := handler.RateLimitSettings

		spanDisplay := fmt.Sprintf("rateLimit.%s", rtSettings.Id)
		ctx, span := wrenchContext.GetSpan2(ctx, spanDisplay)
		defer span.End()

		uClient, _ := connections.GetRedisConnection(rtSettings.RedisConnectionId)
		limiter := redis_rate.NewLimiter(uClient) // check if limiter can be static
		key := handler.getKey(wrenchContext, bodyContext)
		limit := redis_rate.Limit{
			Rate:   rtSettings.RequestsPerSecond,
			Burst:  rtSettings.BurstLimit,
			Period: time.Second,
		}
		res, err := limiter.Allow(ctx, key, limit)

		if err != nil {
			handler.setError(err, http.StatusInternalServerError, span, wrenchContext, bodyContext)
		} else {
			if res.Allowed == 0 {
				bodyContext.SetHeader("Retry-After", fmt.Sprintf("%d", res.RetryAfter/time.Second))
				handler.setError(errors.New("rate limit exceeded"), http.StatusTooManyRequests, span, wrenchContext, bodyContext)
			}
		}

		handler.setSpanAttributes(span, rtSettings.RedisConnectionId, key)
		duration := time.Since(start).Seconds() * 1000
		handler.metricRecord(ctx, duration)
	}

	if handler.Next != nil {
		handler.Next.Do(ctx, wrenchContext, bodyContext)
	}
}

func (handler *RateLimitHandler) getKey(wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) string {
	var keyTemp string

	rtSettings := handler.RateLimitSettings

	if rtSettings.RouteEnabled {
		keyTemp = handler.EndpointSettings.Route
	}

	if len(rtSettings.Keys) > 0 {
		for _, keyRef := range rtSettings.Keys {
			value := contexts.GetCalculatedValue(keyRef, wrenchContext, bodyContext, nil)
			keyTemp += fmt.Sprint(value)
		}
	}

	keyArray := []byte(fmt.Sprint(keyTemp))
	return cross_funcs.GetHash(handler.EndpointSettings.Route, sha256.New, keyArray)
}

func (handler *RateLimitHandler) metricRecord(ctx context.Context, duration float64) {
	app.KafkaProducerDuration.Record(ctx, duration)
}

func (handler *RateLimitHandler) setSpanAttributes(span trace.Span, redisConnectionId string, key string) {
	span.SetAttributes(
		attribute.String("gowrench.connections.redis.id", redisConnectionId),
		attribute.String("rate.limit.key", key),
	)
}

func (handler *RateLimitHandler) setError(err error, httpStatusCode int, span trace.Span, wrenchContext *contexts.WrenchContext, bodyContext *contexts.BodyContext) {
	bodyContext.HttpStatusCode = httpStatusCode
	bodyContext.ContentType = "text/plain"
	wrenchContext.SetHasError(span, err.Error(), err)
}

func (handler *RateLimitHandler) SetNext(next Handler) {
	handler.Next = next
}
