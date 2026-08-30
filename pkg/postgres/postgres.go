package postgres

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// Config describes how to build the connection pool.
type Config struct {
	DSN             string
	MaxConns        int32
	MinConns        int32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration

	// Tracer, when set, is attached to every connection as pgx's QueryTracer.
	// pgxpool also traces pool acquire calls with it when it implements
	// pgxpool.AcquireTracer (otelpgx's tracer does).
	Tracer pgx.QueryTracer
}

// NewPool builds a pgx connection pool and verifies connectivity with a ping.
func NewPool(ctx context.Context, cfg Config) (*pgxpool.Pool, error) {
	poolCfg, err := pgxpool.ParseConfig(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("parsing postgres dsn: %w", err)
	}

	if cfg.MaxConns > 0 {
		poolCfg.MaxConns = cfg.MaxConns
	}

	if cfg.MinConns > 0 {
		poolCfg.MinConns = cfg.MinConns
	}

	if cfg.MaxConnLifetime > 0 {
		poolCfg.MaxConnLifetime = cfg.MaxConnLifetime
	}

	if cfg.MaxConnIdleTime > 0 {
		poolCfg.MaxConnIdleTime = cfg.MaxConnIdleTime
	}

	if cfg.Tracer != nil {
		poolCfg.ConnConfig.Tracer = cfg.Tracer
	}

	pool, err := pgxpool.NewWithConfig(ctx, poolCfg)
	if err != nil {
		return nil, fmt.Errorf("creating postgres pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		return nil, fmt.Errorf("pinging postgres: %w", err)
	}

	return pool, nil
}

// Pool wraps *pgxpool.Pool with samber/do lifecycle hooks so the container can
// close it on shutdown and include it in health checks.
type Pool struct {
	*pgxpool.Pool
}

// Shutdown closes the underlying pool. Implements samber/do's ShutdownerWithError.
func (p *Pool) Shutdown() error {
	p.Close()

	return nil
}

// HealthCheck pings the database. Implements samber/do's Healthchecker.
func (p *Pool) HealthCheck() error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	return p.Ping(ctx)
}
