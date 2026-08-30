//go:build integration

// Package test holds integration tests that exercise the service and
// repository layers together against a real PostgreSQL instance (the
// transport layer is out of scope).
//
// Run them with `make test-integration` (requires Docker); a plain `go test ./...` skips
// this package entirely because of the build tag above.
package test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/Mort4lis/booking-tennis-court/test/testutil"
)

// pg is the single Postgres container shared by every suite in this binary.
// Each suite creates its own database on it (see testutil.Suite.SetupDatabase)
// rather than getting its own container, so suites stay isolated from each
// other without paying container-startup costs.
var pg *testutil.Postgres

func TestMain(m *testing.M) {
	os.Exit(runTests(m))
}

// runTests wraps m.Run so the container gets torn down through a normal defer
// instead of being skipped by TestMain's os.Exit.
func runTests(m *testing.M) int {
	ctx := context.Background()

	var err error

	pg, err = testutil.StartPostgres(ctx)
	if err != nil {
		fmt.Fprintln(os.Stderr, "starting postgres container:", err)

		return 1
	}
	defer func() { _ = pg.Terminate(ctx) }()

	return m.Run()
}
