# Copilot instructions for `gormtest`

## Commands

### Build and run
- Build from the repository root with `go build ./...`.
- Run locally with `go run .` (the application entrypoint is the top-level `main.go`).
- Build the container setup with `docker compose build`.
- Run the containerized app with `docker compose up --build`.

### Tests
- Run the full suite with `go test ./...`.
- Run a single test with `go test ./path/to/package -run TestName`.
- Current targeted tests include `go test ./... -run TestInitDBIsIdempotent` and `go test ./... -run TestInitDBWithoutSeed`.

### Lint
- There is no dedicated lint command or lint target checked into this repository.

## Architecture

- The application is a single Gin server defined in the root `main.go`. Startup loads `.env` when present, resolves runtime config from environment variables, opens the SQLite database, loads HTML templates from `templates/*.html`, and registers the HTTP routes.
- Persistence is handled through GORM models in `internal/models/models.go`. `Livre` and `Recette` use GORM's embedded `gorm.Model`, and `Livre` owns a many-to-many relation to `Recette` through `livre_recette`.
- Startup always calls `initDB()` before serving traffic. That function runs `AutoMigrate`, and seeding is controlled by the `SEED_DB` environment variable.
- The UI is server-rendered HTML plus HTMX. `/` renders `templates/index.html`, which uses HTMX requests to load `/recettes` by default and swap either `/recettes` or `/livres` into the `#content` container. The partial templates live in `templates/tables.html`.
- Kubernetes deployment manifests live under `build/`. The deployment exposes `/health` for startup and liveness probes and injects `POSTGRESQL_URL` from a generated secret via `build/kustomization.yaml`.
- `compose.yaml` is the local container workflow. It builds from `Dockerfile`, reads defaults from `.env`, and mounts a named volume to `/app/data` so the SQLite database persists across container restarts.

## Repository-specific conventions

- Treat the repository root as the application package. The current Makefile, Dockerfile, and Compose setup all build from the root package, not a `cmd/gormtest` subdirectory.
- Runtime configuration comes from environment variables, with `.env` used for local defaults. The key variables are `APP_PORT`, `SQLITE_PATH`, and `SEED_DB`.
- Even though the deployment manifests and environment wiring mention `POSTGRESQL_URL`, the current application code still uses the SQLite driver. If you change database behavior, update both the runtime code and the deployment assets together.
- The `/livres` handler preloads `Recettes` before rendering, so relationship data is expected to be available in memory when templates need it.
- HTML responses are split between a full-page template (`index.html`) and named partial templates (`livresTable`, `recettesTable`) rendered by the data endpoints. Keep that split when adding new HTMX-driven views.
