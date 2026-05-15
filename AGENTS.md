# AGENTS.md

## Build & Run

- **Entry point**: `main.go` at repository root. The Makefile and Dockerfile incorrectly reference `cmd/gormtest/main.go` (that path does not exist).
- **Run locally**: `go run .` (listens on `:8080` by default)
- **Build**: `go build -o bin/gormtest main.go`
- **No tests exist** — there are no `*_test.go` files in the repo.

## Makefile / Container Issues

- `make build` is broken: it builds `cmd/gormtest/main.go` which does not exist. Fix or use `go build` directly.
- `make run` expects a Podman image built by the broken `build` target.
- The Dockerfile also references `cmd/gormtest/main.go`; update the `COPY`/`RUN go build` lines if you use it.

## Tech Stack

- Go 1.24
- [Gin](https://gin-gonic.com/) web framework
- [GORM](https://gorm.io/) with SQLite driver (`gorm.db` file in repo root)
- HTML templates in `templates/*.html`, static files served from `./assets`

## Database

- SQLite file `gorm.db` is created automatically on startup (gitignored).
- Migrations run automatically via `AutoMigrate` in `main.go` for `models.Livre` and `models.Recette`.
- Seed data is inserted only if the `Recette` table is empty.

## Endpoints

| Route | Description |
|-------|-------------|
| `GET /` | Renders `index.html` |
| `GET /livres` | Renders `livresTable` with preloaded `Recettes` |
| `GET /recettes` | Renders `recettesTable` |
| `GET /health` | JSON `{ "status": "ok" }` |

## Project Layout

```
main.go                 # Application entry point and routing
internal/models/        # GORM models (Livre, Recette)
templates/              # HTML templates (index.html, tables.html)
assets/                 # Static files served at /assets
```

## Environment

- `POSTGRESQL_URL` env var is checked but **unused** — the app always uses SQLite regardless.
- The `.gitpod.yml` init task (`go get && go build ./... && go test ./...`) will fail because `go get` is deprecated for modules and there are no tests.
