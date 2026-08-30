package cli

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	"github.com/spf13/cobra"

	"github.com/Mort4lis/booking-tennis-court/internal/config"
	"github.com/Mort4lis/booking-tennis-court/migrations"
	"github.com/Mort4lis/booking-tennis-court/pkg/postgres"
)

func newMigrateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Manage database schema migrations (embedded in the binary)",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(newMigrateUpCmd(), newMigrateDownCmd())

	return cmd
}

func newMigrateUpCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "up",
		Short: "Apply all pending migrations",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			m, closeFn, err := newMigrator()
			if err != nil {
				return err
			}
			defer closeFn()

			if err = m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("applying migrations: %w", err)
			}

			cmd.Println("migrations applied")

			return nil
		},
	}
}

func newMigrateDownCmd() *cobra.Command {
	var (
		step int
		all  bool
	)

	cmd := &cobra.Command{
		Use:   "down",
		Short: "Roll back migrations (--step N or --all)",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, _ []string) error {
			if !all && step <= 0 {
				return errors.New("specify --step N (N > 0) or --all")
			}

			m, closeFn, err := newMigrator()
			if err != nil {
				return err
			}
			defer closeFn()

			rollback := func() error { return m.Steps(-step) }
			if all {
				rollback = m.Down
			}

			if err = rollback(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
				return fmt.Errorf("rolling back migrations: %w", err)
			}

			cmd.Println("rollback complete")

			return nil
		},
	}

	cmd.Flags().IntVar(&step, "step", 0, "number of migrations to roll back")
	cmd.Flags().BoolVar(&all, "all", false, "roll back every migration")
	cmd.MarkFlagsMutuallyExclusive("step", "all")

	return cmd
}

// newMigrator builds a migrate.Migrate over the embedded SQL files and a
// database/sql connection using the pgx driver. It only needs BTC_POSTGRES_*.
// The returned func closes the underlying source and database.
func newMigrator() (*migrate.Migrate, func(), error) {
	cfg, err := config.LoadPostgres()
	if err != nil {
		return nil, nil, err
	}

	return postgres.NewMigrator(cfg.DSN, migrations.FS)
}
