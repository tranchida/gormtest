# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
go run .            # Run locally (listens on :8080)
go build -o bin/gormtest main.go   # Build binary
go vet ./...        # Lint
```

> **Note:** The `Makefile` and `Dockerfile` reference `cmd/gormtest/main.go` which does not exist — use the commands above directly. There are no tests (`*_test.go` files).

## Architecture

Single-file web app (`main.go`) using **Gin** for HTTP routing and **GORM** with a local SQLite file (`gorm.db`).

- `main.go` — server setup, routing, DB init with seed data, and all HTTP handlers
- `internal/models/` — GORM models: `Livre` (book) with a many-to-many join to `Recette` (recipe) via `livre_recette`
- `templates/` — Go HTML templates: `index.html` (base page), `tables.html` (defines named blocks `livresTable` and `recettesTable` rendered by HTMX partials)

The database is created automatically on startup via `AutoMigrate`. Seed data (one `Livre` with two `Recettes`) is inserted only when the `Recette` table is empty.

## Endpoints

| Route | Response |
|-------|----------|
| `GET /` | `index.html` |
| `GET /livres` | `livresTable` partial (Livre + preloaded Recettes) |
| `GET /recettes` | `recettesTable` partial |
| `GET /health` | `{"status":"ok"}` |
