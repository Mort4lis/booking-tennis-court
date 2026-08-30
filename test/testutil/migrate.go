package testutil

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	"github.com/Mort4lis/booking-tennis-court/migrations"
	"github.com/Mort4lis/booking-tennis-court/pkg/postgres"
)

// Migrate applies every embedded migration (see migrations.FS) to the
// database at dsn, reusing postgres.NewMigrator.
func Migrate(dsn string) error {
	m, closeFn, err := postgres.NewMigrator(dsn, migrations.FS)
	if err != nil {
		return err
	}
	defer closeFn()

	if err := m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("applying migrations: %w", err)
	}

	return nil
}
