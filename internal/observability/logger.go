package observability

import (
	"log/slog"
	"os"
	"strings"

	"github.com/Mort4lis/booking-tennis-court/internal/config"
)

// NewLogger builds the application logger from config and installs it as the
// slog default so package-level slog calls share the same handler.
func NewLogger(cfg *config.Config) *slog.Logger {
	opts := &slog.HandlerOptions{Level: parseLevel(cfg.Log.Level)}

	var handler slog.Handler
	if strings.EqualFold(cfg.Log.Format, "text") {
		handler = slog.NewTextHandler(os.Stdout, opts)
	} else {
		handler = slog.NewJSONHandler(os.Stdout, opts)
	}

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger
}

func parseLevel(s string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}
