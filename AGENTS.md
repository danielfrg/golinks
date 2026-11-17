# Agent Instructions

Golinks is a self-hosted short-link service.
It builds one Go binary that serves redirects, a
Huma HTTP API, an MCP endpoint, and an embedded Svelte 5 frontend backed by SQLite.

Read `PRODUCT.md` for user intent and `DESIGN.md` before changing the frontend.

The `README.md` and `justfile` own setup and usage.
They are written for humans: avoid big restructures, but update them when changes are user facing.

## Architecture

- Go for the API in /cmd and /internal.
- TypeScript and Svelte 5 for the frontend in /frontend.
- Should build into a single binary with an embedded frontend.

## Behavior Invariants

- Keep SQLite rollback journaling; it supports databases on network storage.
- All internal endpoints live under `/_/`; never reserve ordinary short names for operational routes.
- MCP write tools use the same validation as the HTTP API.

## Development

Use the `just` commands (see the `justfile`):

- `just install`: install the frontend dependencies.
- `just api`: run the Go API (serves the frontend embedded by the last `just build`).
- `just frontend`: run the Vite dev server with HMR (proxies the API).
- `just check`: gofmt, vet, race tests, oxfmt/oxlint, frontend tests and svelte-check.
- `just build`: build the single binary with the embedded frontend.
- `just run`: build, then run the binary.
- `just image-check`: build the container image and run `scripts/container-smoke.sh`.

## Frontend Conventions

`DESIGN.md` owns the theme, typography, layout, and interaction rules; read it before changing the frontend.

`frontend/src/design-system/` is synced from the shared design system and never edited by hand.
App-specific styles live in `frontend/src/app.css`.

main.ts imports the synced design-system CSS first, then palette.css if any, then app.css last, so app
rules win ties with the theme; prefer tokens and class selectors over !important; never edit
frontend/src/design-system/ (it is synced from the design-system repo with `pnpm sync`).
