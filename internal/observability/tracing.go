package observability

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace/noop"

	"github.com/Mort4lis/booking-tennis-court/internal/config"
)

// TracerProvider is the app's tracer-provider handle with a lifecycle hook.
type TracerProvider struct {
	shutdown func(context.Context) error
}

// NewTracerProvider installs the global text-map propagator and tracer provider.
//
// With BTC_TRACING_ENABLED=true it exports spans over OTLP/gRPC to
// BTC_TRACING_OTLP_ENDPOINT (run `make compose-up` for a local Jaeger, UI on
// http://localhost:16686). Otherwise a no-op provider is installed and spans are
// dropped.
func NewTracerProvider(ctx context.Context, cfg *config.Config) (*TracerProvider, error) {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	noopShutdown := func(context.Context) error { return nil }

	if !cfg.Tracing.Enabled {
		otel.SetTracerProvider(noop.NewTracerProvider())

		return &TracerProvider{shutdown: noopShutdown}, nil
	}

	if cfg.Tracing.OtlpEndpoint == "" {
		return nil, errors.New("tracing is enabled but BTC_TRACING_OTLP_ENDPOINT is empty")
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Tracing.OtlpEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("creating otlp trace exporter: %w", err)
	}

	provider := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(resource.NewSchemaless(
			attribute.String("service.name", cfg.Tracing.ServiceName),
		)),
	)
	otel.SetTracerProvider(provider)

	slog.Info("tracing enabled",
		slog.String("otlp_endpoint", cfg.Tracing.OtlpEndpoint),
		slog.String("service_name", cfg.Tracing.ServiceName),
	)

	return &TracerProvider{shutdown: provider.Shutdown}, nil
}

// Shutdown flushes pending spans. Implements samber/do's
// ShutdownerWithContextAndError.
func (t *TracerProvider) Shutdown(ctx context.Context) error {
	return t.shutdown(ctx)
}
