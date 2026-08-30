// Package config loads the application configuration from environment variables.
//
// Every field carries an explicit `envconfig` tag with its name; combined with
// the BTC_ prefix and the sub-struct name that yields e.g. BTC_TELEGRAM_TOKEN
// and BTC_POSTGRES_MAX_CONN_LIFETIME. See .env.example for the full list with
// defaults.
package config

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

const envPrefix = "BTC"

// Config is the root configuration tree.
type Config struct {
	Log      Log
	Telegram Telegram
	Postgres Postgres
	Metrics  Metrics
	Tracing  Tracing
}

// Log configures structured logging.
type Log struct {
	Level  string `envconfig:"LEVEL" default:"info"`
	Format string `envconfig:"FORMAT" default:"json"`
}

// Telegram configures the Telegram bot client.
type Telegram struct {
	Token string `envconfig:"TOKEN" required:"true"`
	Debug bool   `envconfig:"DEBUG" default:"false"`
}

// Postgres configures the PostgreSQL connection pool.
type Postgres struct {
	DSN             string        `envconfig:"DSN" required:"true"`
	MaxConns        int32         `envconfig:"MAX_CONNS" default:"10"`
	MinConns        int32         `envconfig:"MIN_CONNS" default:"2"`
	MaxConnLifetime time.Duration `envconfig:"MAX_CONN_LIFETIME" default:"1h"`
	MaxConnIdleTime time.Duration `envconfig:"MAX_CONN_IDLE_TIME" default:"30m"`
}

// Metrics configures the Prometheus scrape endpoint fed by the OTel meter provider.
type Metrics struct {
	Enabled bool   `envconfig:"ENABLED" default:"true"`
	Addr    string `envconfig:"ADDR" default:":9090"`
	Path    string `envconfig:"PATH" default:"/metrics"`
}

// Tracing configures OpenTelemetry tracing over OTLP/gRPC. Disabled by default;
// when enabled, OtlpEndpoint must point at an OTLP collector (the local Jaeger
// from docker-compose.yml listens on localhost:4317).
type Tracing struct {
	Enabled      bool   `envconfig:"ENABLED" default:"false"`
	OtlpEndpoint string `envconfig:"OTLP_ENDPOINT" default:""`
	ServiceName  string `envconfig:"SERVICE_NAME" default:"booking-tennis-court"`
}

// Load reads and validates the full configuration from the process environment.
// It is used by the long-running `serve` command.
func Load() (*Config, error) {
	var cfg Config

	if err := envconfig.Process(envPrefix, &cfg); err != nil {
		return nil, fmt.Errorf("processing environment config: %w", err)
	}

	return &cfg, nil
}

// LoadPostgres reads only the database configuration (BTC_POSTGRES_*). Admin
// commands such as `migrate` use it so they do not require the full environment
// (e.g. a Telegram token) just to touch the schema.
func LoadPostgres() (*Postgres, error) {
	var pg Postgres

	if err := envconfig.Process(envPrefix+"_POSTGRES", &pg); err != nil {
		return nil, fmt.Errorf("processing postgres config: %w", err)
	}

	return &pg, nil
}
