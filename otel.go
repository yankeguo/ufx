package ufx

import (
	"context"
	"errors"
	"github.com/yankeguo/rg"
	"go.uber.org/fx"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/log"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/trace"
)

func SetupOTEL(lc fx.Lifecycle) (err error) {
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

	// Set up propagator.
	propagator = propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	)
	otel.SetTextMapPropagator(propagator)

	// Set up trace provider.
	traceExporter = rg.Must(stdouttrace.New(stdouttrace.WithPrettyPrint()))
	traceProvider = trace.NewTracerProvider(trace.WithBatcher(traceExporter, trace.WithBatchTimeout(time.Second)))
	otel.SetTracerProvider(traceProvider)

	// Set up meter provider.
	metricExporter = rg.Must(stdoutmetric.New())
	metricProvider = metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter,
			metric.WithInterval(3*time.Second))),
	)
	otel.SetMeterProvider(metricProvider)

	// Set up logger provider.
	loggerExporter = rg.Must(stdoutlog.New())
	loggerProvider = log.NewLoggerProvider(
		log.WithProcessor(log.NewBatchProcessor(loggerExporter)),
	)
	global.SetLoggerProvider(loggerProvider)
	return
}
