# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

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
go run ./cmd/go-scaffold new <service-name> [--module <module-path>] [--output <dir>]
```

## Architecture

The tool itself is intentionally minimal — only `cobra` as a dependency.

- **`cmd/go-scaffold/main.go`** — CLI entry point with two commands: `new` and `version`
- **`internal/generator/generator.go`** — walks `templates.FS`, substitutes `service` in file/dir names with `ServiceName`, strips `.tmpl` suffix, and executes each template
- **`internal/generator/vars.go`** — `Vars` struct and `NewVars()`: derives `ServiceNameTitle` (title-case), `EntityName` (singular PascalCase, e.g. "payments" → "Payment"), `ModuleName`, `GoVersion`
- **`templates/`** — embedded filesystem (`//go:embed`) of all `.tmpl` files; `embed.go` exports `FS`

## Template conventions

- File and directory names containing `service` are renamed to the lowercase service name at generation time (e.g. `internal/service/service.go.tmpl` → `internal/payments/payments.go`)
- Template variables: `{{ .ServiceName }}`, `{{ .ServiceNameTitle }}`, `{{ .EntityName }}`, `{{ .ModuleName }}`, `{{ .GoVersion }}`
- The `singular()` helper simply strips a trailing `s` — service names ending in anything other than a regular English plural may produce unexpected entity names

## Generated service structure

Generated services follow Clean Architecture layers:
- `internal/domain/` — entities, validation, errors
- `internal/usecase/` — business logic, adapter interfaces (with `//go:generate mockery`)
- `internal/adapter/` — Postgres (pgx) and Redis implementations
- `internal/controller/http/` — Chi HTTP handlers
- `pkg/` — reusable packages: httpserver, router, logger, postgres, redis, transaction, render
- `api/http/` — OpenAPI spec; `gen/` — oapi-codegen output
- `migration/postgres/` — golang-migrate SQL files

## Adding a new template file

1. Create the `.tmpl` file under `templates/` using the template variables above
2. Use `service` in the path anywhere the service name should appear
3. No registration needed — the generator walks the entire embedded FS automatically
