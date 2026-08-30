package observability

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	otelprom "go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"

	"github.com/Mort4lis/booking-tennis-court/internal/config"
)

const meterName = "github.com/Mort4lis/booking-tennis-court"

// Metrics holds the application's OpenTelemetry instruments. Add new instruments
// here and record them from the relevant layer.
type Metrics struct {
	provider *sdkmetric.MeterProvider

	// UpdatesTotal counts processed Telegram updates, labelled by handler.
	UpdatesTotal metric.Int64Counter

	// HandlerErrors counts handler failures, labelled by handler.
	HandlerErrors metric.Int64Counter

	// HandlerDuration records handler execution time in seconds, labelled by handler.
	HandlerDuration metric.Float64Histogram
}

// NewMetrics configures an OTel meter provider backed by a Prometheus exporter
// (registered against reg) and pre-declares the app instruments.
func NewMetrics(cfg *config.Config, reg *prometheus.Registry) (*Metrics, error) {
	exporter, err := otelprom.New(otelprom.WithRegisterer(reg))
	if err != nil {
		return nil, fmt.Errorf("creating prometheus exporter: %w", err)
	}

	provider := sdkmetric.NewMeterProvider(
		sdkmetric.WithReader(exporter),
		sdkmetric.WithResource(resource.NewSchemaless(
			attribute.String("service.name", cfg.Tracing.ServiceName),
		)),
	)
	otel.SetMeterProvider(provider)

	meter := provider.Meter(meterName)

	updates, err := meter.Int64Counter(
		"telegram_updates_total",
		metric.WithDescription("Total number of processed Telegram updates."),
	)
	if err != nil {
		return nil, fmt.Errorf("creating updates counter: %w", err)
	}

	handlerErrors, err := meter.Int64Counter(
		"telegram_handler_errors_total",
		metric.WithDescription("Total number of Telegram handler errors."),
	)
	if err != nil {
		return nil, fmt.Errorf("creating handler errors counter: %w", err)
	}

	duration, err := meter.Float64Histogram(
		"telegram_handler_duration_seconds",
		metric.WithDescription("Telegram handler execution time in seconds."),
		metric.WithUnit("s"),
	)
	if err != nil {
		return nil, fmt.Errorf("creating handler duration histogram: %w", err)
	}

	return &Metrics{
		provider:        provider,
		UpdatesTotal:    updates,
		HandlerErrors:   handlerErrors,
		HandlerDuration: duration,
	}, nil
}

// Shutdown flushes and stops the meter provider. Implements samber/do's
// ShutdownerWithContextAndError.
func (m *Metrics) Shutdown(ctx context.Context) error {
	if err := m.provider.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutting down meter provider: %w", err)
	}

	return nil
}

// MetricsServer serves the Prometheus scrape endpoint on a dedicated HTTP port.
type MetricsServer struct {
	srv     *http.Server
	enabled bool
	log     *slog.Logger
}

// NewMetricsServer builds the HTTP server exposing a metrics url path from reg.
func NewMetricsServer(cfg *config.Config, reg *prometheus.Registry, log *slog.Logger) *MetricsServer {
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{}))
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	return &MetricsServer{
		srv: &http.Server{
			Addr:              cfg.Metrics.Addr,
			Handler:           mux,
			ReadHeaderTimeout: 5 * time.Second,
		},
		enabled: cfg.Metrics.Enabled,
		log:     log,
	}
}

// Start runs the server until ctx is cancelled (graceful shutdown) or
// ListenAndServe fails. It is a no-op when metrics are disabled.
func (s *MetricsServer) Start(ctx context.Context) error {
	if !s.enabled {
		s.log.InfoContext(ctx, "metrics server disabled")

		<-ctx.Done()

		return nil
	}

	errCh := make(chan error, 1)

	go func() {
		s.log.InfoContext(ctx, "metrics server listening", slog.String("addr", s.srv.Addr))

		if err := s.srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := s.srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("shutting down metrics server: %w", err)
		}

		return nil
	case err := <-errCh:
		return fmt.Errorf("metrics server: %w", err)
	}
}
