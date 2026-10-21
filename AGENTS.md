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

## Frontend Conventions

`DESIGN.md` owns the theme, typography, layout, and interaction rules; read it before changing the frontend.

## Deployment

Daniel (danielfrg/danrodriguez) deploys this service in his homelab; the deployment source is in `~/code/danielfrg/homelab`.
