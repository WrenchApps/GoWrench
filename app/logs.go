package app

import (
	"fmt"
	"log"
	"time"

	otelLog "go.opentelemetry.io/otel/log"
)

type WrenchErrorLog struct {
	Message string
	Error   error
}

func LogInfo(msg string) {
	log.Print(msg)

	if Logger != nil {
		var record otelLog.Record
		record.SetSeverity(otelLog.SeverityInfo)
		record.SetBody(otelLog.StringValue(msg))
		record.SetTimestamp(time.Now())

		Logger.Emit(GetContext(), record)
	}
}

func LogWarning(msg string) {
	log.Print(msg)

	if Logger != nil {
		var record otelLog.Record
		record.SetSeverity(otelLog.SeverityWarn)
		record.SetBody(otelLog.StringValue(msg))
		record.SetTimestamp(time.Now())

		Logger.Emit(GetContext(), record)
	}
}

func LogError(err WrenchErrorLog) {
	log.Print(err)

	if Logger != nil {
		var record otelLog.Record
		record.SetSeverity(otelLog.SeverityError)
		record.SetBody(otelLog.StringValue(fmt.Sprint(err)))
		record.SetTimestamp(time.Now())

		Logger.Emit(GetContext(), record)
	}
}

func LogError2(msg string, err error) {
	LogError(WrenchErrorLog{Message: msg, Error: err})
}
