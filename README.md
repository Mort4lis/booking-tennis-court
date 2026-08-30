# booking-tennis-court

Production-ready Telegram bot that books a tennis court.

## Stack

| Concern  | Choice                                                                       |
|----------|------------------------------------------------------------------------------|
| CLI      | `spf13/cobra`                                                                |
| Config   | `kelseyhightower/envconfig` (prefix `BTC_`)                                  |
| DI       | `samber/do/v2`                                                               |
| Telegram | `go-telegram/bot` (long polling)                                             |
| Database | `jackc/pgx/v5` + `golang-migrate` (embedded SQL)                             |
| Logging  | `log/slog` (JSON)                                                            |
| Metrics  | `go.opentelemetry.io/otel` SDK → Prometheus `/metrics` endpoint              |
| Tracing  | `go.opentelemetry.io/otel` SDK → OTLP/gRPC; pgx traced via `exaring/otelpgx` |
| Tests    | `stretchr/testify`, `testcontainers-go`, `go-testfixtures` (integration)     |

## Architecture & Layout

It ships the full skeleton - Clean Architecture layers, Dependency Injection, env-based configuration, structured
logging, OpenTelemetry metrics, database migrations, linting, Docker and a Makefile.

```
cmd/bot/                 entrypoint
internal/
  cli/                   serve, migrate, version commands
  config/                environment configuration
  app/                   DI composition root + runtime lifecycle
  transport/telegram/    long-polling bot, thin handlers
  service/booking/       application use-cases
  repository/booking/    persistence layer
  entity/                domain types
  dto/                   transport contracts (Data Transfer Objects)
  observability/         slog, metrics, tracing
pkg/postgres/            reusable pgx pool constructor
migrations/              embedded, timestamp-named SQL migrations
test/                    integration tests (build tag: integration); testutil/ + testdata/ fixtures
```

## Getting started

```bash
cp .env.example .env          # then set BTC_TELEGRAM_TOKEN
make compose-up               # start PostgreSQL + Jaeger
make migrate                  # apply migrations (needs BTC_POSTGRES_DSN)
make run                      # build + run the bot (long polling)
```

`curl localhost:9090/metrics` once running. Set `BTC_TRACING_ENABLED=true` to export spans (pgx queries included) to the
local Jaeger — UI on `http://localhost:16686`.

## Configuration

All variables are prefixed `BTC_` and grouped by sub-struct — see [`.env.example`](.env.example). 

`serve` needs `BTC_TELEGRAM_*` and `BTC_POSTGRES_*`; `migrate` needs only `BTC_POSTGRES_*`.

## Make targets

| Target                        | Description                                           |
|-------------------------------|-------------------------------------------------------|
| `build`                       | compile to `bin/bot`                                  |
| `run`                         | build and run `bot serve`                             |
| `lint`                        | `golangci-lint run ./...`                             |
| `test`                        | `test-unit` then `test-integration`                   |
| `test-unit`                   | unit tests with the race detector                     |
| `test-integration`            | integration tests (`-tags=integration`, needs Docker) |
| `migrate`                     | apply all pending migrations                          |
| `migrate-down-step`           | roll back one migration                               |
| `migrate-down`                | roll back all migrations                              |
| `docker-build`                | build the container image                             |
| `compose-up` / `compose-down` | local PostgreSQL                                      |

To get more information about a target run `make help`.

## CI and Docker Hub

Pull requests and pushes to `master` run lint, unit tests, and integration tests in order. A failed job skips the jobs
that depend on it. On `master`, passing all checks also builds and publishes a `linux/amd64` image tagged with the full
commit SHA.

In the GitHub repository settings under **Secrets and variables** → **Actions**, configure these values.

| Name | Type | Value |
|------|------|-------|
| `DOCKERHUB_IMAGE` | Variable | Docker Hub repository in `namespace/repository` format |
| `DOCKERHUB_USERNAME` | Variable | Docker Hub username |
| `DOCKERHUB_PASSWORD` | Secret | Password for the Docker Hub account with write access to the target repository |

Docker Hub requires a personal access token for CLI sign-in when the account uses two-factor authentication or enforced
SSO. In that case, store a token with write access in `DOCKERHUB_PASSWORD`.

## Migrations

Embedded into the binary via `go:embed`; applied with the `migrate` subcommand. Create a new pair with:

```bash
migrate create -ext sql -dir migrations -format 20060102150405 add_something
```

## Where to add code

- Command parsing / new commands → `internal/transport/telegram/handlers.go`
- Use-cases → `internal/service/booking`
- Queries → `internal/repository/booking`
- New env vars → `internal/config/config.go`
- New metrics → `internal/observability/metrics.go` (`Metrics` struct)
