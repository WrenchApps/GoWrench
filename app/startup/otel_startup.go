package startup

import (
	"context"
	"log"
	"wrench/app/manifest/application_settings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer() func(context.Context) error {
	ctx := context.Background()

	app := application_settings.ApplicationSettingsStatic
	otelSetting := app.Service.Otel

	if otelSetting != nil && otelSetting.Enable {

		exporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(otelSetting.CollectorUrl),
			otlptracehttp.WithInsecure(),
		)
		if err != nil {
			log.Fatalf("error to create otlp http exporter: %v", err)
		}

		tp := trace.NewTracerProvider(
			trace.WithBatcher(exporter),
			trace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(app.Service.Name),
				semconv.ServiceVersion(app.Service.Version),
			)),
		)

		otel.SetTracerProvider(tp)
		return tp.Shutdown
	}
	return nil
}
