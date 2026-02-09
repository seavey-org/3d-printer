# Agent Guide (3d-printer)

This repo is a small Go (Gin) backend + Vue 3 (Vite) frontend, shipped as a single Docker image.

Repo map:
- `backend/`: Go API server (also serves frontend dist in prod)
- `frontend/`: Vue 3 SPA (Pinia, Tailwind)
- `deployment/`: nginx + systemd units/scripts for the production server
- `Dockerfile`, `docker-compose.yml`: build and runtime

## Build / Lint / Test

CI source of truth: `.github/workflows/ci-cd.yml`.

### Backend (Go)

From `backend/`:

```sh
# run all tests (matches CI behavior)
go test -race -coverprofile=coverage.out ./...

# run a single package
go test ./internal/handlers

# run a single test by name (exact match)
go test ./internal/handlers -run '^TestListTimelapses$'

# run a single test with verbose output
go test ./internal/api -run '^TestCORSHeaders$' -v
```

Lint (CI installs `golangci-lint` v2.8.0):

```sh
golangci-lint run --timeout=5m
```

Formatting/imports (enforced via golangci-lint):
- `gofmt` and `goimports` are enabled in `backend/.golangci.yml`.

Run the server locally:

```sh
PORT=8080 TIMELAPSE_DIR=./videos STREAM_M3U8_PATH=./live/stream.m3u8 go run ./cmd/server
```

### Frontend (Vue / Vite)

From `frontend/`:

```sh
npm ci

# dev server
npm run dev

# lint + autofix
npm run lint
npm run lint:fix

# type checking (strict)
npm run type-check

# production build
npm run build
```

Notes:
- There is no dedicated frontend unit test runner configured right now (no Vitest/Jest). CI relies on ESLint + `vue-tsc` + build.
- Vite dev proxy targets the Go server on `http://localhost:8080` for `/api`, `/live`, `/videos` (see `frontend/vite.config.ts`).

### Docker (prod-like)

```sh
docker build -t 3d-printer-app:local .
IMAGE_TAG=local docker compose up -d
```

`docker-compose.yml` bind-mounts production host paths under `/var/www/...`; local runs may need edits or stub folders.

## Code Style Guidelines

### General
- Prefer small, obvious changes over clever refactors.
- No secrets in git. Never commit `.env`/credential files.
- Keep behavior consistent with production deployment in `deployment/`.

### Go (backend)

Imports/formatting:
- Always run `gofmt` (and prefer `goimports`).
- Import groups: stdlib, third-party, then local (`github.com/codyseavey/3d-printer/...`). `goimports` handles this.

Types/naming:
- Exported types/functions use `PascalCase`, unexported use `camelCase`.
- Constructor pattern: `NewXxx(...) *Xxx` (see `handlers.NewTimelapseHandler`).
- JSON structs use explicit tags and stable field names (see `backend/internal/models/timelapse.go`).

Error handling/logging:
- Check errors immediately, return early.
- Log server-side details with context, return a simple JSON error to clients (current pattern: `gin.H{"error": "..."}`).
- Avoid panics in request paths; reserve `log.Fatalf` for unrecoverable startup failures.

HTTP/API conventions:
- Routes are registered in `backend/internal/api/routes.go`.
- Prefer explicit status codes, and keep the response shape predictable.
- CORS config lives in the router setup; keep origins configurable via `CORS_ALLOWED_ORIGINS`.

Testing:
- Use `httptest` + Gin test mode (`gin.SetMode(gin.TestMode)`).
- Prefer real filesystem state via `t.TempDir()` over mocks (tests already do this).
- Keep tests deterministic, no network calls.

### Vue/TypeScript (frontend)

Formatting/linting:
- ESLint is the formatter of record here.
- Style rules: single quotes, no semicolons (see `frontend/eslint.config.js`).
- Prefer `import type { ... }` for type-only imports.

Types:
- `tsconfig` is strict and checks unused locals/params; keep types accurate.
- Avoid `any`. If you must, isolate it and document why in code.

Vue conventions:
- Use `<script setup lang="ts">`.
- Type `defineProps` and `defineEmits`.
- Clean up side effects (timers/listeners) in `onBeforeUnmount`.

State/data fetching:
- Centralize remote calls in `frontend/src/services/api.ts`.
- Use Pinia stores for shared state (`frontend/src/stores/*`).
- For async actions: set loading flags, handle errors, keep UI resilient.

Routing:
- Routes live in `frontend/src/router.ts` and use history mode.
- Ensure server-side SPA fallback is preserved (backend `NoRoute` serves `index.html` when `FRONTEND_DIST_PATH` is set).

Styling/accessibility:
- Tailwind is used for styling; keep classes readable.
- Add `aria-*` labels for icon-only buttons and meaningful roles for loading/error states.

## Cursor / Copilot Rules

- No `.cursor/rules/`, `.cursorrules`, or `.github/copilot-instructions.md` found in this repo at the time of writing.
