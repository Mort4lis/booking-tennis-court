package testutil

import (
	"database/sql"
	"fmt"

	"github.com/go-testfixtures/testfixtures/v3"
)

// LoadFixtures truncates every table referenced by the YAML files in the dir and
// re-inserts their rows into db.
func LoadFixtures(db *sql.DB, dir string) error {
	loader, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect("postgres"),
		testfixtures.Directory(dir),
	)
	if err != nil {
		return fmt.Errorf("building fixtures loader for %q: %w", dir, err)
	}

	if err = loader.Load(); err != nil {
		return fmt.Errorf("loading fixtures from %q: %w", dir, err)
	}

	return nil
}
