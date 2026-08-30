// Package postgres provides thin, reusable building blocks around pgx: a
// connection pool constructor (NewPool) and a golang-migrate migrator builder
// (NewMigrator). It deliberately does not import anything from internal/ so
// it stays usable by other projects.
package postgres
