// Package app owns the application runtime lifecycle: it builds the dependency
// graph via internal/app/deps, starts the Telegram bot and the metrics server,
// and drives graceful shutdown.
package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os/signal"
	"syscall"

	"golang.org/x/sync/errgroup"

	"github.com/Mort4lis/booking-tennis-court/internal/app/deps"
)

// Run builds the application, starts the Telegram bot and the metrics server,
// and blocks until an interrupt signal triggers a graceful shutdown.
func Run(ctx context.Context) error {
	ctx, stop := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	container, err := deps.New(ctx)
	if err != nil {
		return fmt.Errorf("building application: %w", err)
	}

	logger := container.Logger()
	tgBot := container.Bot()
	metricsSrv := container.MetricsServer()

	group, groupCtx := errgroup.WithContext(ctx)

	group.Go(func() error { return metricsSrv.Start(groupCtx) })
	group.Go(func() error { return tgBot.Start(groupCtx) })

	// groupCtx is cancelled either by an interrupt signal (ctx.Err() != nil) or
	// because one of the components returned an error on its own.
	<-groupCtx.Done()

	if ctx.Err() != nil {
		logger.Info("shutdown signal received, stopping")
	} else {
		logger.Error("one of the component stopped unexpectedly, shutting down the rest")
	}

	// Wait for the transport and metrics servers to stop before tearing the
	// container down, so nothing is using the pool when it closes.
	runErr := group.Wait()

	// Container.Shutdown detaches from the cancelled ctx and applies its own timeout.
	if report := container.Shutdown(ctx); !report.Succeed {
		logger.Error("errors during shutdown", slog.String("report", report.Error()))
	}

	if runErr != nil && !errors.Is(runErr, context.Canceled) {
		return fmt.Errorf("runtime error: %w", runErr)
	}

	return nil
}
