package app

import (
	"context"
	"fmt"
	"log"
	"time"

	otelLog "go.opentelemetry.io/otel/log"
)

type WrenchErrorLog struct {
	Message string
	Error   error
}

func LogInfo(ctx context.Context, msg string) {
	var record otelLog.Record
	record.SetSeverity(otelLog.SeverityInfo)
	record.SetBody(otelLog.StringValue(msg))
	record.SetTimestamp(time.Now())

	log.Print(msg)
	Logger.Emit(ctx, record)
}

func LogWarning(ctx context.Context, msg string) {
	var record otelLog.Record
	record.SetSeverity(otelLog.SeverityWarn)
	record.SetBody(otelLog.StringValue(msg))
	record.SetTimestamp(time.Now())

	log.Print(msg)
	Logger.Emit(ctx, record)
}

func LogError(ctx context.Context, err WrenchErrorLog) {
	var record otelLog.Record
	record.SetSeverity(otelLog.SeverityError)
	record.SetBody(otelLog.StringValue(fmt.Sprint(err)))
	record.SetTimestamp(time.Now())

	log.Fatal(err)
	Logger.Emit(ctx, record)
}
