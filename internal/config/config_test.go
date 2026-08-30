package config_test

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/Mort4lis/booking-tennis-court/internal/config"
)

func TestLoad(t *testing.T) {
	t.Setenv("BTC_TELEGRAM_TOKEN", "test-token")
	t.Setenv("BTC_POSTGRES_DSN", "postgres://localhost:5432/test")
	t.Setenv("BTC_POSTGRES_MAX_CONN_LIFETIME", "2h")
	t.Setenv("BTC_LOG_LEVEL", "debug")
	t.Setenv("BTC_METRICS_ADDR", ":9999")

	cfg, err := config.Load()
	require.NoError(t, err)

	assert.Equal(t, "test-token", cfg.Telegram.Token)
	assert.Equal(t, "postgres://localhost:5432/test", cfg.Postgres.DSN)
	assert.Equal(t, 2*time.Hour, cfg.Postgres.MaxConnLifetime)
	assert.Equal(t, "debug", cfg.Log.Level)
	assert.Equal(t, ":9999", cfg.Metrics.Addr)

	// Defaults.
	assert.Equal(t, "json", cfg.Log.Format)
	assert.Equal(t, int32(10), cfg.Postgres.MaxConns)
	assert.True(t, cfg.Metrics.Enabled)
	assert.False(t, cfg.Tracing.Enabled)
}

// envconfig treats a set-but-empty var as present, so the required keys must be
// genuinely removed from the environment here — which t.Setenv cannot do, and
// which makes the test unsafe to run in parallel.
//
//nolint:paralleltest,usetesting // mutates the real process environment by unsetting vars
func TestLoad_MissingRequired(t *testing.T) {
	for _, key := range []string{"BTC_TELEGRAM_TOKEN", "BTC_POSTGRES_DSN"} {
		if orig, ok := os.LookupEnv(key); ok {
			t.Cleanup(func() { _ = os.Setenv(key, orig) })
		}

		require.NoError(t, os.Unsetenv(key))
	}

	_, err := config.Load()
	require.Error(t, err)
}
