package app

import (
	"context"

	"go.opentelemetry.io/otel/log"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"
	sdklog "go.opentelemetry.io/otel/sdk/log"
)

const ENV_PORT = "PORT"
const ENV_PATH_FILE_CONFIG string = "PATH_FILE_CONFIG"
const ENV_PATH_FOLDER_ENV_FILES string = "PATH_FOLDER_ENV_FILES"
const ENV_APP_ENV string = "APP_ENV"
const ENV_RUN_BASH_FILES_BEFORE_STARTUP string = "RUN_BASH_FILES_BEFORE_STARTUP"

var contextInitiated context.Context

var Tracer = otel.Tracer("trace")
var Meter = otel.Meter("meter")

var HttpServerDuration metric.Float64Histogram
var HttpClientDurantion metric.Float64Histogram
var KafkaProducerDurtation metric.Float64Histogram
var NatsPublishDurtation metric.Float64Histogram
var SnsPublishDurtation metric.Float64Histogram

var LoggerProvider *sdklog.LoggerProvider
var Logger log.Logger

func InitMetrics() {
	HttpServerDuration, _ = Meter.Float64Histogram("gowrench_http_server_duration_ms")
	HttpClientDurantion, _ = Meter.Float64Histogram("gowrench_http_client_duration_ms")
	KafkaProducerDurtation, _ = Meter.Float64Histogram("gowrench_kafka_producer_duration_ms")
	NatsPublishDurtation, _ = Meter.Float64Histogram("gowrench_nats_publish_duration_ms")
	SnsPublishDurtation, _ = Meter.Float64Histogram("gowrench_sns_publish_duration_ms")
}

func InitLogger(lp *sdklog.LoggerProvider) {
	LoggerProvider = lp
	Logger = LoggerProvider.Logger("logger")
}

func SetContext(ctx context.Context) {
	contextInitiated = ctx
}

func GetContext() context.Context {
	return contextInitiated
}
