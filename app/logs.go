package app

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel/log"
)

func LogInfo(ctx context.Context, msg string) {
	var record log.Record
	record.SetSeverity(log.SeverityInfo)
	record.SetBody(log.StringValue(msg))
	record.SetTimestamp(time.Now())

	Logger.Emit(ctx, record)
}

func LogWarning(ctx context.Context, msg string) {
	var record log.Record
	record.SetSeverity(log.SeverityWarn)
	record.SetBody(log.StringValue(msg))
	record.SetTimestamp(time.Now())

	Logger.Emit(ctx, record)
}

func LogError(ctx context.Context, err error) {
	var record log.Record
	record.SetSeverity(log.SeverityError)
	record.SetBody(log.StringValue(fmt.Sprint(err)))
	record.SetTimestamp(time.Now())

	Logger.Emit(ctx, record)
}
