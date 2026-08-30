package telegram

import (
	"context"
	"log/slog"
	"time"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

// instrument wraps a handler with metrics recording and structured logging.
func (b *Bot) instrument(name string, next bot.HandlerFunc) bot.HandlerFunc {
	return func(ctx context.Context, api *bot.Bot, update *models.Update) {
		start := time.Now()
		attrs := metric.WithAttributes(attribute.String("handler", name))

		b.metrics.UpdatesTotal.Add(ctx, 1, attrs)

		next(ctx, api, update)

		elapsed := time.Since(start)
		b.metrics.HandlerDuration.Record(ctx, elapsed.Seconds(), attrs)

		b.log.InfoContext(ctx, "handled telegram update",
			slog.String("handler", name),
			slog.Duration("duration", elapsed),
		)
	}
}
