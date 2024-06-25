package telemetry

import (
	"context"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/propagation"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.25.0"
	"go.opentelemetry.io/otel/trace"
)

type Telemetry struct {
	mexp  *otlpmetricgrpc.Exporter
	mp    *sdkmetric.MeterProvider
	Meter metric.Meter

	texp   *otlptrace.Exporter
	tp     *sdktrace.TracerProvider
	Tracer trace.Tracer
}

func newMeterProvider(ctx context.Context, resource *resource.Resource) (*otlpmetricgrpc.Exporter, *sdkmetric.MeterProvider) {
	exporter, err := otlpmetricgrpc.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry.newMeterProvider")
	}

	mp := sdkmetric.NewMeterProvider(
		sdkmetric.WithResource(resource),
		sdkmetric.WithReader(sdkmetric.NewPeriodicReader(exporter)),
	)
	otel.SetMeterProvider(mp)

	return exporter, mp
}

func newTraceProvider(ctx context.Context, resource *resource.Resource) (*otlptrace.Exporter, *sdktrace.TracerProvider) {
	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry.newTraceProvider")
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)

	return exporter, tp
}

func New(ctx context.Context, scope, appID string) *Telemetry {
	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(appID),
		),
	)
	if err != nil {
		log.Panic().Err(err).Msg("telemetry.New")
	}

	mexp, mp := newMeterProvider(ctx, res)
	texp, tp := newTraceProvider(ctx, res)

	err = runtime.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("telemetry.New")
	}

	return &Telemetry{
		mexp:  mexp,
		mp:    mp,
		Meter: mp.Meter(scope),

		texp:   texp,
		tp:     tp,
		Tracer: tp.Tracer(scope),
	}
}

func (t *Telemetry) Shutdown(ctx context.Context) {
	if err := t.mexp.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("telemetry.Shutdown: mexp")
	}

	if err := t.mp.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("telemetry.Shutdown: mp")
	}

	if err := t.texp.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("telemetry.Shutdown: texp")
	}

	if err := t.tp.Shutdown(ctx); err != nil {
		log.Error().Err(err).Msg("telemetry.Shutdown: tp")
	}
}
