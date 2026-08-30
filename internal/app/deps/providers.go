package deps

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/exaring/otelpgx"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/samber/do/v2"

	"github.com/Mort4lis/booking-tennis-court/internal/config"
	"github.com/Mort4lis/booking-tennis-court/internal/observability"
	repository "github.com/Mort4lis/booking-tennis-court/internal/repository/postgres"
	bookingsvc "github.com/Mort4lis/booking-tennis-court/internal/service/booking"
	"github.com/Mort4lis/booking-tennis-court/internal/transport/telegram"
	"github.com/Mort4lis/booking-tennis-court/pkg/postgres"
)

// registerConfig loads the configuration from the environment.
func registerConfig(i do.Injector) {
	do.Provide(i, func(do.Injector) (*config.Config, error) {
		return config.Load()
	})
}

// registerObservability register the logger, the Prometheus registry, the meter
// provider, the metrics server and the tracer provider. ctx bounds the OTLP
// exporter setup done by the tracer provider.
func registerObservability(ctx context.Context, i do.Injector) {
	do.Provide(i, func(i do.Injector) (*slog.Logger, error) {
		return observability.NewLogger(do.MustInvoke[*config.Config](i)), nil
	})

	do.Provide(i, func(do.Injector) (*prometheus.Registry, error) {
		return prometheus.NewRegistry(), nil
	})

	do.Provide(i, func(i do.Injector) (*observability.Metrics, error) {
		metrics, err := observability.NewMetrics(
			do.MustInvoke[*config.Config](i),
			do.MustInvoke[*prometheus.Registry](i),
		)
		if err != nil {
			return nil, fmt.Errorf("initializing metrics: %w", err)
		}

		return metrics, nil
	})

	do.Provide(i, func(i do.Injector) (*observability.MetricsServer, error) {
		return observability.NewMetricsServer(
			do.MustInvoke[*config.Config](i),
			do.MustInvoke[*prometheus.Registry](i),
			do.MustInvoke[*slog.Logger](i),
		), nil
	})

	do.Provide(i, func(i do.Injector) (*observability.TracerProvider, error) {
		tracer, err := observability.NewTracerProvider(ctx, do.MustInvoke[*config.Config](i))
		if err != nil {
			return nil, fmt.Errorf("initializing tracer provider: %w", err)
		}

		return tracer, nil
	})
}

// registerInfrastructure register the PostgreSQL pool, with pgx queries traced
// through otelpgx when tracing is enabled.
func registerInfrastructure(ctx context.Context, i do.Injector) {
	do.Provide(i, func(i do.Injector) (*postgres.Pool, error) {
		cfg := do.MustInvoke[*config.Config](i)

		// Force the tracer provider so otelpgx picks up the installed global one.
		do.MustInvoke[*observability.TracerProvider](i)

		pgCfg := postgres.Config{
			DSN:             cfg.Postgres.DSN,
			MaxConns:        cfg.Postgres.MaxConns,
			MinConns:        cfg.Postgres.MinConns,
			MaxConnLifetime: cfg.Postgres.MaxConnLifetime,
			MaxConnIdleTime: cfg.Postgres.MaxConnIdleTime,
		}
		if cfg.Tracing.Enabled {
			pgCfg.Tracer = otelpgx.NewTracer()
		}

		pool, err := postgres.NewPool(ctx, pgCfg)
		if err != nil {
			return nil, fmt.Errorf("initializing postgres pool: %w", err)
		}

		return &postgres.Pool{Pool: pool}, nil
	})
}

// registerDomain wires the repository and the application service.
func registerDomain(i do.Injector) {
	do.Provide(i, func(i do.Injector) (*repository.BookingRepository, error) {
		return repository.NewBookingRepository(do.MustInvoke[*postgres.Pool](i).Pool), nil
	})

	do.Provide(i, func(i do.Injector) (*bookingsvc.Service, error) {
		return bookingsvc.NewService(do.MustInvoke[*repository.BookingRepository](i)), nil
	})
}

// registerTransport register the Telegram bot.
func registerTransport(i do.Injector) {
	do.Provide(i, func(i do.Injector) (*telegram.Bot, error) {
		cfg := do.MustInvoke[*config.Config](i)

		tgBot, err := telegram.New(
			telegram.Config{Token: cfg.Telegram.Token, Debug: cfg.Telegram.Debug},
			do.MustInvoke[*bookingsvc.Service](i),
			do.MustInvoke[*observability.Metrics](i),
			do.MustInvoke[*slog.Logger](i),
		)
		if err != nil {
			return nil, fmt.Errorf("initializing telegram bot: %w", err)
		}

		return tgBot, nil
	})
}
