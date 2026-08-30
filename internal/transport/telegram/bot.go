package telegram

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/go-telegram/bot"

	"github.com/Mort4lis/booking-tennis-court/internal/dto"
	"github.com/Mort4lis/booking-tennis-court/internal/observability"
)

// Config is the transport-local configuration, mapped from the app config so the
// transport layer does not depend on the whole config tree.
type Config struct {
	Token string
	Debug bool
}

// BookingService is the port the transport depends on.
type BookingService interface {
	CreateBooking(ctx context.Context, req dto.CreateBookingRequest) (dto.BookingResponse, error)
}

// Bot wraps the Telegram client and its handler registrations.
type Bot struct {
	api     *bot.Bot
	svc     BookingService
	metrics *observability.Metrics
	log     *slog.Logger
}

// New builds the Telegram bot and registers the command handlers.
func New(cfg Config, svc BookingService, metrics *observability.Metrics, log *slog.Logger) (*Bot, error) {
	b := &Bot{svc: svc, metrics: metrics, log: log}

	opts := []bot.Option{
		bot.WithDefaultHandler(b.instrument("default", b.handleDefault)),
	}
	if cfg.Debug {
		opts = append(opts, bot.WithDebug())
	}

	api, err := bot.New(cfg.Token, opts...)
	if err != nil {
		return nil, fmt.Errorf("creating telegram bot: %w", err)
	}

	b.api = api

	b.register("/start", b.handleStart)
	b.register("/help", b.handleHelp)
	b.register("/book", b.handleBook)

	return b, nil
}

func (b *Bot) register(command string, handler bot.HandlerFunc) {
	b.api.RegisterHandler(
		bot.HandlerTypeMessageText,
		command,
		bot.MatchTypePrefix,
		b.instrument(command, handler),
	)
}

// Start begins long polling and blocks until ctx is cancelled.
func (b *Bot) Start(ctx context.Context) error {
	b.log.InfoContext(ctx, "starting telegram bot (long polling)")

	b.api.Start(ctx)

	return ctx.Err()
}
