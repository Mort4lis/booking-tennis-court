package postgres

import (
	"database/sql"
	"fmt"
	"io/fs"

	"github.com/golang-migrate/migrate/v4"
	pgxmigrate "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	_ "github.com/jackc/pgx/v5/stdlib" // registers the "pgx" database/sql driver
)

// NewMigrator builds a *migrate.Migrate over a migration source (an embed.FS
// of *.sql files, or any other fs.FS) and a database/sql connection to dsn
// using the pgx driver. It only builds the migrator — callers drive it
// themselves (Up, Down, Steps) so both a one-shot "apply everything" caller
// and a CLI with up/down/step subcommands can share this constructor.
//
// The returned close func releases both the source and the database handle
// and must be called once the migrator is no longer needed.
func NewMigrator(dsn string, src fs.FS) (*migrate.Migrate, func(), error) {
	source, err := iofs.New(src, ".")
	if err != nil {
		return nil, nil, fmt.Errorf("loading migrations: %w", err)
	}

	db, err := sql.Open("pgx", dsn)
	if err != nil {
		_ = source.Close()

		return nil, nil, fmt.Errorf("opening database: %w", err)
	}

	// Any failure past this point must release both handles; on success
	// m.Close (returned to the caller) takes over that responsibility.
	cleanup := func() {
		_ = source.Close()
		_ = db.Close()
	}

	driver, err := pgxmigrate.WithInstance(db, &pgxmigrate.Config{})
	if err != nil {
		cleanup()

		return nil, nil, fmt.Errorf("creating migrate driver: %w", err)
	}

	m, err := migrate.NewWithInstance("iofs", source, "postgres", driver)
	if err != nil {
		cleanup()

		return nil, nil, fmt.Errorf("creating migrator: %w", err)
	}

	// m.Close closes both the source and the database.
	return m, func() { _, _ = m.Close() }, nil
}
