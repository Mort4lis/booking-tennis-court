package testutil

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

// postgresImage pins the same PostgreSQL version used by the local dev stack (see
// docker-compose.yml), so integration tests run against the same major version.
const postgresImage = "postgres:17-alpine"

// Postgres is a container shared by every suite in the test binary.
type Postgres struct {
	container *postgres.PostgresContainer

	// DSN points at the container's maintenance database.
	DSN string
}

// StartPostgres boots a fresh Postgres container. Call Terminate once the test
// binary is done with it — typically from TestMain, after m.Run returns.
func StartPostgres(ctx context.Context) (*Postgres, error) {
	container, err := postgres.Run(ctx, postgresImage,
		postgres.WithDatabase("postgres"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		return nil, fmt.Errorf("starting postgres container: %w", err)
	}

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		_ = container.Terminate(ctx)

		return nil, fmt.Errorf("building admin dsn: %w", err)
	}

	return &Postgres{container: container, DSN: dsn}, nil
}

// Terminate stops and removes the container.
func (p *Postgres) Terminate(ctx context.Context) error {
	return p.container.Terminate(ctx)
}
