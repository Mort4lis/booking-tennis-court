package testutil

import (
	"context"
	"fmt"
	"net/url"

	"github.com/jackc/pgx/v5"
)

// CreateDatabase creates a new, empty database named name on the server behind
// baseDSN and returns a DSN pointing at it.
func CreateDatabase(ctx context.Context, baseDSN, name string) (string, error) {
	conn, err := pgx.Connect(ctx, baseDSN)
	if err != nil {
		return "", fmt.Errorf("connecting to database: %w", err)
	}
	defer func() { _ = conn.Close(ctx) }()

	// CREATE DATABASE doesn't support query parameters; name is always a value we
	// generate ourselves (never user input), so quoting the identifier is enough.
	ident := pgx.Identifier{name}.Sanitize()
	if _, err = conn.Exec(ctx, "CREATE DATABASE "+ident); err != nil {
		return "", fmt.Errorf("creating database %q: %w", name, err)
	}

	return rewriteDSNPath(baseDSN, name)
}

// DropDatabase drops the database named name from the server behind baseDSN.
// It forcibly disconnects any remaining clients first (WITH FORCE, PostgreSQL
// 13+), so a lingering pool connection can't block teardown.
func DropDatabase(ctx context.Context, baseDSN, name string) error {
	conn, err := pgx.Connect(ctx, baseDSN)
	if err != nil {
		return fmt.Errorf("connecting to admin database: %w", err)
	}
	defer func() { _ = conn.Close(ctx) }()

	ident := pgx.Identifier{name}.Sanitize()
	if _, err = conn.Exec(ctx, "DROP DATABASE IF EXISTS "+ident+" WITH (FORCE)"); err != nil {
		return fmt.Errorf("dropping database %q: %w", name, err)
	}

	return nil
}

// rewriteDSNPath rewrites dsn's path to point at the given database name.
func rewriteDSNPath(dsn, name string) (string, error) {
	u, err := url.Parse(dsn)
	if err != nil {
		return "", fmt.Errorf("parsing dsn: %w", err)
	}

	u.Path = "/" + name

	return u.String(), nil
}
