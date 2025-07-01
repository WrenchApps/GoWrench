package startup

import (
	"context"
	"log"
	"os"
	"wrench/app/manifest/application_settings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlplog/otlploghttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdklog "go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
)

func InitTracer() func(context.Context) error {
	ctx := context.Background()
	app := application_settings.ApplicationSettingsStatic
	otelSetting := app.Service.Otel

	if otelSetting != nil && otelSetting.Enable {

		res := resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(app.Service.Name),
			semconv.ServiceVersion(app.Service.Version),
		)

		exporter, err := otlptracehttp.New(ctx,
			otlptracehttp.WithEndpoint(otelSetting.CollectorUrl),
			otlptracehttp.WithInsecure(),
		)

		var stdoutExporter *stdouttrace.Exporter
		var stdoutErr error

		if otelSetting.TraceConsoleExport {
			stdoutExporter, stdoutErr = stdouttrace.New(
				stdouttrace.WithPrettyPrint(),
				stdouttrace.WithWriter(os.Stdout),
			)
		}

		if err != nil || stdoutErr != nil {
			log.Fatalf("error to create otlp exporter: %v %v", err, stdoutErr)
		} else {

		}

		var tp *trace.TracerProvider

		if stdoutExporter != nil {
			tp = trace.NewTracerProvider(
				trace.WithBatcher(exporter),
				trace.WithBatcher(stdoutExporter),
				trace.WithResource(res),
			)
		} else {
			tp = trace.NewTracerProvider(
				trace.WithBatcher(exporter),
				trace.WithResource(res),
			)
		}

		otel.SetTracerProvider(tp)
		return tp.Shutdown
	}
	return nil
}

func InitMeter() func(context.Context) error {
	ctx := context.Background()
	app := application_settings.ApplicationSettingsStatic
	otelSetting := app.Service.Otel

	if otelSetting != nil && otelSetting.Enable {
		exporter, err := otlpmetrichttp.New(ctx,
			otlpmetrichttp.WithEndpoint(otelSetting.CollectorUrl),
			otlpmetrichttp.WithInsecure(),
		)

		var stdoutExporter metric.Exporter
		var stdoutErr error

		if otelSetting.MetricConsoleExport {
			stdoutExporter, stdoutErr = stdoutmetric.New(stdoutmetric.WithWriter(os.Stdout))
		}

		if err != nil || stdoutErr != nil {
			log.Fatalf("failed to create metric exporter: %v %v", err, stdoutErr)
		} else {
			res := resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(app.Service.Name),
				semconv.ServiceVersion(app.Service.Version),
			)

			var provider *metric.MeterProvider

			if stdoutExporter != nil {
				provider = metric.NewMeterProvider(
					metric.WithResource(res),
					metric.WithReader(metric.NewPeriodicReader(exporter)),
					metric.WithReader(metric.NewPeriodicReader(stdoutExporter)),
				)
			} else {
				provider = metric.NewMeterProvider(
					metric.WithResource(res),
					metric.WithReader(metric.NewPeriodicReader(exporter)),
				)
			}

			otel.SetMeterProvider(provider)
			return provider.Shutdown
		}
		return nil
	}
	return nil
}

func InitLogProvider() *sdklog.LoggerProvider {
	ctx := context.Background()

	app := application_settings.ApplicationSettingsStatic
	otelSetting := app.Service.Otel

	res := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceName(app.Service.Name),
		semconv.ServiceVersion(app.Service.Version),
	)

	exp, err := otlploghttp.New(ctx,
		otlploghttp.WithEndpoint(otelSetting.CollectorUrl),
		otlploghttp.WithInsecure(),
	)

	var logExport *stdoutlog.Exporter
	var logErr error

	if otelSetting.LogConsoleExport {
		logExport, logErr = stdoutlog.New(stdoutlog.WithPrettyPrint())
	}

	if err != nil || logErr != nil {
		log.Fatalf("log exporter: %v %v", err, logErr)
	}

	if logExport != nil {
		return sdklog.NewLoggerProvider(
			sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
			sdklog.WithProcessor(sdklog.NewBatchProcessor(logExport)),
			sdklog.WithResource(res),
		)

	} else {
		return sdklog.NewLoggerProvider(
			sdklog.WithProcessor(sdklog.NewBatchProcessor(exp)),
			sdklog.WithResource(res),
		)
	}
}
