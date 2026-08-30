// Package testutil holds shared infrastructure for integration tests: booting a
// disposable Postgres container, creating an isolated database per test suite,
// applying migrations to it, and loading YAML fixtures. It is a regular package
// (no _test.go suffix), so it can be imported from any suite file.
package testutil
