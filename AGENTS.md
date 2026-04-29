# AGENTS.md

## Module & Runtime

- Go module: `go_sample_code` (not matching the directory name)
- Go 1.25.0, no Makefile/Taskfile/CI — all commands are manual `go` invocations

## Commands

```bash
# Run server (requires DB + Redis)
go run cmd/server/main.go -c config/config.local.yaml

# Tests — single package, or all
go test ./internal/middleware/...
go test ./...

# Ent codegen after editing internal/ent/schema/
go install entgo.io/ent/cmd/ent@latest
ent generate ./internal/ent/schema

go fmt ./...   # no linter (no .golangci.yml)
```

## Architecture

Layered design wired via `uber-go/fx`. All layers (Handler, Service, Repo) are interface-based with constructor injection.

```
cmd/server/main.go        — fx app: providers + RegisterHooks (route registration)
internal/handler/{domain}/ — Fiber HTTP handlers
internal/service/{domain}/ — Business logic
internal/repo/{domain}/    — Ent ORM data access
internal/middleware/        — Fiber middleware chain
internal/errno/             — Typed error codes with HTTP status mapping
internal/ent/               — Ent generated code (DO NOT EDIT; edit schema/ then regenerate)
internal/ent/schema/        — Entity definitions + mixin/
internal/database/          — DB/Redis client setup, health checks
pkg/                        — jwx, logger, trace, metrics, rbac, validator
config/                     — YAML config templates
devops/                     — Docker Compose configs per subdirectory
```

## Key Facts

- **Routes**: Only `/api/health` is registered. User endpoints (`POST /api/users`, `GET /api/users/:id`, etc.) are **not wired** in `RegisterHooks` — add them when ready.
- **Config**: YAML with `yaml` struct tags. Config struct at `cmd/server/config.go`. If no config file exists, runs with hardcoded defaults (Postgres localhost:5432, Redis localhost:6379). Pass `-c <path>` to specify.
- **Error responses**: `errno.Decode(data, err)` → `{code, message, data}` JSON. Error code ranges: 0=OK, 1000-1099 common, 1100-1199 auth, 1200-1299 file, 1300+ database, 2000+ business.
- **Validator**: Wraps `go-playground/validator` via `pkg/validator`. Struct tags: `json` (body), `query` (query params), `params` (path params — custom parser in `validator.go` since Fiber lacks built-in).
- **Ent codegen**: No `go:generate`. Run `ent generate ./internal/ent/schema` manually. Only edit `internal/ent/schema/*.go` and `internal/ent/schema/mixin/`.
- **Middleware wiring** (order in `RegisterHooks`): Recovery → RateLimit → Trace → Metrics → Logger. Auth and RBAC middleware exist but are **not wired**.
- **Rate limit**: Enabled by default (20 req/s, burst 50, skips `/api/health`).

## Local Dev

```bash
docker compose -f devops/database/docker-compose.yml up -d
# Optional: full observability (Grafana, Prometheus, Tempo, Loki, OTel)
docker compose -f devops/grafana.v1/docker-compose.yml up -d
```

## Testing

- All tests are unit-level — no DB/Redis needed. Use `stretchr/testify` for assertions.
- Pattern: mock the layer below (e.g., handler tests mock the service interface), create a `fiber.App`, register routes, and call `app.Test(httptest.NewRequest(...))`.
- In-memory tracer: `trace.NewInMemoryProvider()` for handler/service tests.
- Rate limit tests use `mockLogger` (no-op Logger impl in the test package).
