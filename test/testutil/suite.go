package testutil

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/suite"
)

// Suite is the base for every integration test suite. Each concrete suite owns
// a database of its own — created and migrated by SetupDatabase, dropped by
// TearDownSuite — plus connections tests need.
type Suite struct {
	suite.Suite

	// Pool is a connection pool for building the real repository/service under test.
	Pool *pgxpool.Pool
	// DB is a connection for loading fixtures.
	DB *sql.DB

	baseDSN string
	dbName  string
}

// SetupDatabase creates a database named name on the server behind baseDSN,
// migrates it, and opens Pool and DB against it. Call it from the concrete suite's
// own SetupSuite, before building anything that needs a connection.
func (s *Suite) SetupDatabase(ctx context.Context, baseDSN, name string) {
	dsn, err := CreateDatabase(ctx, baseDSN, name)
	s.Require().NoError(err, "creating test database %q", name)

	s.Require().NoError(Migrate(dsn), "applying migrations to %q", name)

	pool, err := pgxpool.New(ctx, dsn)
	s.Require().NoError(err, "opening pgx pool")

	db, err := sql.Open("pgx", dsn)
	s.Require().NoError(err, "opening database/sql connection")

	s.baseDSN = baseDSN
	s.dbName = name

	s.Pool = pool
	s.DB = db
}

// LoadFixtures resets the suite's database to the state described by the YAML
// files in dir. Call it from the concrete suite's SetupTest so every test
// method starts from the same known state.
func (s *Suite) LoadFixtures(dir string) {
	s.Require().NoError(LoadFixtures(s.DB, dir), "loading fixtures from %q", dir)
}

// TearDownSuite closes connections and drops the suite's database.
func (s *Suite) TearDownSuite() {
	if s.DB != nil {
		_ = s.DB.Close()
	}

	if s.Pool != nil {
		s.Pool.Close()
	}

	if s.baseDSN != "" && s.dbName != "" {
		_ = DropDatabase(context.Background(), s.baseDSN, s.dbName)
	}
}
