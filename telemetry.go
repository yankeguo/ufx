package ufx

import (
	"context"
	"errors"
	"github.com/yankeguo/rg"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	"go.uber.org/fx"
	"time"
)

func SetupTelemetry(lc fx.Lifecycle) (err error) {
	defer rg.Guard(&err)

	var (
		propagator propagation.TextMapPropagator

		traceExporter trace.SpanExporter
		traceProvider *trace.TracerProvider

		metricExporter metric.Exporter
		metricProvider *metric.MeterProvider

		loggerExporter log.Exporter
		loggerProvider *log.LoggerProvider
	)

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			shutdown := func(ctx context.Context, v any) error {
				if v == nil {
					return nil
				}
				if e, ok := v.(interface{ Shutdown(context.Context) error }); ok {
					return e.Shutdown(ctx)
				}
				return nil
			}
			return errors.Join(
				shutdown(ctx, loggerProvider),
				shutdown(ctx, loggerExporter),
				shutdown(ctx, metricProvider),
				shutdown(ctx, metricExporter),
				shutdown(ctx, traceProvider),
				shutdown(ctx, traceExporter),
			)
		},
	})

	ctx := context.Background()

	res := rg.Must(resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithProcess(),
		resource.WithOS(),
		resource.WithContainer(),
		resource.WithHost(),
	))

	propagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	traceExporter = rg.Must(otlptracehttp.New(ctx))
	traceProvider = trace.NewTracerProvider(
		trace.WithBatcher(traceExporter,
			trace.WithBatchTimeout(time.Second),
		),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(traceProvider)

	metricExporter = rg.Must(otlpmetrichttp.New(ctx))
	metricProvider = metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithInterval(3*time.Second)),
		),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(metricProvider)

	loggerExporter = rg.Must(stdoutlog.New())
	loggerProvider = log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(loggerExporter)),
		log.WithResource(res),
	)
	global.SetLoggerProvider(loggerProvider)
	return
}
