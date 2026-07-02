# AGENTS.md

This file provides guidance to Codex (Codex.ai/code) when working with code in this repository.

## Язык общения

Всегда отвечать пользователю на русском языке.

## What this project is

`go-scaffold` is a CLI tool that generates complete, production-ready Go microservices following Clean Architecture. It walks an embedded filesystem of `.tmpl` files and renders them via `text/template` into an output directory.

## Commands

```bash
# Build the CLI tool
go build ./cmd/go-scaffold

# Run tests
go test ./...

# Run the tool locally
go run ./cmd/go-scaffold new <service-name> [--module <module-path>] [--output <dir>] [--redis] [--kafka-consumer] [--kafka-producer]
```

## Architecture

The tool itself is intentionally minimal — only `cobra` and `golang.org/x/mod` (module path validation) as dependencies.

- **`cmd/go-scaffold/main.go`** — CLI entry point with two commands: `new` and `version`
- **`internal/generator/generator.go`** — walks `templates.FS`, substitutes the `__service__` placeholder in file/dir names with `ServiceName`, strips `.tmpl` suffix, and executes each template
- **`internal/generator/vars.go`** — `Vars` struct, `Options` struct and `NewVars()`: derives `ServiceNameTitle` (title-case), `EntityName` (singular PascalCase, e.g. "payments" → "Payment"), `ModuleName`, `GoVersion`, plus opt-in feature flags (each has a CLI flag and an interactive prompt; `skippedPaths()` in generator.go excludes the related templates when a feature is off):
  - `WithRedis` (`--redis`) — cache adapter, `pkg/redis`
  - `WithKafkaConsumer` (`--kafka-consumer`) — `internal/controller/kafka_consumer`
  - `WithKafkaProducer` (`--kafka-producer`) — `internal/adapter/kafka_producer`, `internal/controller/worker` (outbox worker), `internal/domain/event.go`, outbox postgres methods + migration, `OutboxReadAndProduce` usecase
  - `WithKafka` — computed (`consumer || producer`): kafka in docker-compose, `pkg/logger/kafka.go`, kafka-go dependency
- **`templates/`** — embedded filesystem (`//go:embed`) of all `.tmpl` files; `embed.go` exports `FS`

## Template conventions

- The `__service__` placeholder in file and directory names is replaced with the lowercase service name at generation time (e.g. `api/http/__service___v1.yaml.tmpl` → `api/http/payments_v1.yaml`)
- Template variables: `{{ .ServiceName }}`, `{{ .ServiceNameTitle }}`, `{{ .EntityName }}`, `{{ .ModuleName }}`, `{{ .GoVersion }}`, `{{ .WithRedis }}`, `{{ .WithKafkaConsumer }}`, `{{ .WithKafkaProducer }}`, `{{ .WithKafka }}`
- The `singular()` helper simply strips a trailing `s` — service names ending in anything other than a regular English plural may produce unexpected entity names

## Generated service structure

Generated services follow Clean Architecture layers:
- `internal/domain/` — entities, validation, errors
- `internal/usecase/` — business logic, adapter interfaces (with `//go:generate mockery`)
- `internal/adapter/` — Postgres (pgx) implementation; with `--redis` also a Redis cache adapter (cache-aside for GET, invalidation on update/close); with `--kafka-producer` also a kafka-go producer adapter
- `internal/controller/http/` — Chi HTTP handlers; with `--kafka-consumer` also `kafka_consumer/` (consumer group, FetchMessage → handle → CommitMessages); with `--kafka-producer` also `worker/` (outbox worker: polls the `outbox` table in batches and produces to Kafka)
- `--kafka-producer` uses the transactional outbox pattern end-to-end: the Create usecase saves a `domain.Event` to the `outbox` table in the same transaction as the entity (`SaveOutboxKafka`); the worker drains it via `ReadOutboxKafka` (`FOR UPDATE SKIP LOCKED` + `DELETE ... RETURNING`) and `Produce`s in one transaction (at-least-once)
- `pkg/` — reusable packages: httpserver, router, logger, postgres, redis, transaction, render
- `api/http/` — OpenAPI spec; `gen/` — oapi-codegen output
- `migration/postgres/` — golang-migrate SQL files (+ `outbox` table with `--kafka-producer`)

## Adding a new template file

1. Create the `.tmpl` file under `templates/` using the template variables above
2. Use `__service__` in the path anywhere the service name should appear
3. No registration needed — the generator walks the entire embedded FS automatically
