# Agent Instructions

Golinks is a self-hosted short-link service.
It builds one Go binary that serves redirects, a
Huma HTTP API, an MCP endpoint, and an embedded Svelte 5 frontend backed by SQLite.

Read `PRODUCT.md` for user intent and `DESIGN.md` before changing the frontend.

The `README.md` and `justfile` own setup and usage.
Do not make big changes to those as they are intended for humans but you can update them if they are user facing.

## Architecture

- Go for the API in /cmd and /internal
- TypeScript and Svelte 5 for the frontend in /frontend
- Should build into a single binary with an embedded frontend

## Behavior Invariants

- Keep SQLite rollback journaling; it supports databases on network storage.
- All internal endpoints live under `/_/`; never reserve ordinary short names for operational routes.
- MCP write tools use the same validation as the HTTP API

## Frontend Conventions

`DESIGN.md` owns the theme, typography, layout, and interaction rules; read it before changing the frontend

- Import shared controls directly from `@kenn-io/kit-ui`; local components own app-specific behavior, not pass-through wrappers
- Use Svelte 5 runes, typed props, and semantic CSS tokens
- Keep desktop and phone layouts, keyboard focus, and recoverable error states.

## Deployment

If you are working for Daniel/danielfrg/danrodriguez. He deploys this service in his homelab, you can find the source in ~/code/danielfrg/homelab to find more details.
