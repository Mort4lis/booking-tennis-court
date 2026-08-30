package deps

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/samber/do/v2"

	"github.com/Mort4lis/booking-tennis-court/internal/config"
	"github.com/Mort4lis/booking-tennis-court/internal/observability"
	"github.com/Mort4lis/booking-tennis-court/internal/transport/telegram"
	"github.com/Mort4lis/booking-tennis-court/pkg/postgres"
)

// shutdownTimeout bounds the whole teardown (pool close, meter and tracer
// flush).
const shutdownTimeout = 10 * time.Second

// Container owns the dependency-injection scope and exposes typed accessors for
// the services the runtime needs.
type Container struct {
	i do.Injector
}

// New registers every provider, then eagerly instantiates the long-lived
// services. On any failure it tears down whatever was already built and returns
// the error. The passed ctx bounds the initial connection attempts (e.g. the
// PostgreSQL ping); it is not retained for the service lifetime.
func New(ctx context.Context) (*Container, error) {
	i := do.New()

	registerConfig(i)
	registerObservability(ctx, i)
	registerInfrastructure(ctx, i)
	registerDomain(i)
	registerTransport(i)

	c := &Container{i: i}

	if err := c.instantiate(); err != nil {
		_ = c.Shutdown(ctx)

		return nil, err
	}

	return c, nil
}

// instantiate builds every long-lived service in dependency order. That order,
// reversed, is also the teardown order: samber/do closes the pool and flushes
// the meter provider and tracer (the bot and metrics server stop earlier, in Run).
func (c *Container) instantiate() error {
	steps := []struct {
		name string
		fn   func() error
	}{
		{"config", func() error { _, err := do.Invoke[*config.Config](c.i); return err }},
		{"logger", func() error { _, err := do.Invoke[*slog.Logger](c.i); return err }},
		{"tracer provider", func() error { _, err := do.Invoke[*observability.TracerProvider](c.i); return err }},
		{"metrics", func() error { _, err := do.Invoke[*observability.Metrics](c.i); return err }},
		{"postgres pool", func() error { _, err := do.Invoke[*postgres.Pool](c.i); return err }},
		{"metrics server", func() error { _, err := do.Invoke[*observability.MetricsServer](c.i); return err }},
		{"telegram bot", func() error { _, err := do.Invoke[*telegram.Bot](c.i); return err }},
	}

	for _, s := range steps {
		if err := s.fn(); err != nil {
			return fmt.Errorf("initializing %s: %w", s.name, err)
		}
	}

	return nil
}

// Shutdown tears the whole graph down, invoking each service's Shutdown hook in
// reverse instantiation order. The passed ctx is used only for its values: it is
// detached from cancellation (teardown must run even when ctx is already
// cancelled by the shutdown signal) and re-bounded by shutdownTimeout.
func (c *Container) Shutdown(ctx context.Context) *do.ShutdownReport {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), shutdownTimeout)
	defer cancel()

	return c.i.ShutdownWithContext(ctx)
}

// Config returns the loaded configuration.
func (c *Container) Config() *config.Config {
	return do.MustInvoke[*config.Config](c.i)
}

// Logger returns the application logger.
func (c *Container) Logger() *slog.Logger {
	return do.MustInvoke[*slog.Logger](c.i)
}

// Metrics returns the OpenTelemetry instruments handle.
func (c *Container) Metrics() *observability.Metrics {
	return do.MustInvoke[*observability.Metrics](c.i)
}

// MetricsServer returns the Prometheus scrape-endpoint server.
func (c *Container) MetricsServer() *observability.MetricsServer {
	return do.MustInvoke[*observability.MetricsServer](c.i)
}

// Bot returns the Telegram transport.
func (c *Container) Bot() *telegram.Bot {
	return do.MustInvoke[*telegram.Bot](c.i)
}
