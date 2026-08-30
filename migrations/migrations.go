// Package migrations embeds the SQL migration files applied by `bot migrate`.
//
// Files follow golang-migrate's naming convention with a timestamp version:
//
//	<YYYYMMDDHHMMSS>_<name>.up.sql
//	<YYYYMMDDHHMMSS>_<name>.down.sql
//
// Create a new pair with:
//
//	migrate create -ext sql -dir migrations -format 20060102150405 <name>
package migrations

import "embed"

// FS holds every *.sql migration file in this directory.
//
//go:embed *.sql
var FS embed.FS
