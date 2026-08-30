package telegram

import (
	"context"
	"log/slog"

	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"

	"github.com/Mort4lis/booking-tennis-court/internal/dto"
)

const (
	startMessage = "Привет! Я помогу забронировать теннисный корт.\n" +
		"Доступные команды: /book, /help"

	helpMessage = "Команды:\n" +
		"/book — забронировать корт\n" +
		"/help — эта справка"
)

func (b *Bot) handleStart(ctx context.Context, api *bot.Bot, update *models.Update) {
	b.reply(ctx, api, update, startMessage)
}

func (b *Bot) handleHelp(ctx context.Context, api *bot.Bot, update *models.Update) {
	b.reply(ctx, api, update, helpMessage)
}

// handleBook runs the full transport -> service -> repository path. The service
// and repository are stubs, so it returns the "not implemented" error for now.
// Parsing the court and slot from the message is the extension point.
func (b *Bot) handleBook(ctx context.Context, api *bot.Bot, update *models.Update) {
	if update.Message == nil {
		return
	}

	_, err := b.svc.CreateBooking(ctx, dto.CreateBookingRequest{TelegramID: update.Message.Chat.ID})
	if err != nil {
		b.handlerError(ctx, api, update, "book", err)

		return
	}

	b.reply(ctx, api, update, "Бронь создана.")
}

func (b *Bot) handleDefault(ctx context.Context, api *bot.Bot, update *models.Update) {
	b.reply(ctx, api, update, "Не понимаю команду. Наберите /help.")
}

func (b *Bot) reply(ctx context.Context, api *bot.Bot, update *models.Update, text string) {
	if update.Message == nil {
		return
	}

	if _, err := api.SendMessage(ctx, &bot.SendMessageParams{
		ChatID: update.Message.Chat.ID,
		Text:   text,
	}); err != nil {
		b.log.ErrorContext(ctx, "sending telegram message", slog.Any("error", err))
		b.metrics.HandlerErrors.Add(ctx, 1)
	}
}

func (b *Bot) handlerError(ctx context.Context, api *bot.Bot, update *models.Update, handler string, err error) {
	b.log.ErrorContext(ctx, "handler failed",
		slog.String("handler", handler),
		slog.Any("error", err),
	)
	b.metrics.HandlerErrors.Add(ctx, 1, metric.WithAttributes(attribute.String("handler", handler)))
	b.reply(ctx, api, update, "Произошла ошибка, попробуйте позже.")
}
