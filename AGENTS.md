# Golinks

Read `PRODUCT.md` and `DESIGN.md` before changing the frontend.
Use Impeccable for UI review and follow `docs/design-workflow.md`.
Do not edit `README.md` unless the user explicitly requests it.

## Architecture

- `cmd/golinks` owns environment configuration, CLI import and server lifecycle.
- `internal/links` defines link contracts, validation and Go URL templates.
- `internal/store` owns SQLite persistence, click history and atomic backups.
- `internal/httpapi` serves redirects and Huma JSON endpoints under `/_/api`.
- `internal/web` embeds the frontend in the standalone Go executable.
- `frontend` uses plain Vite, Svelte 5, TypeScript and Kit UI with pnpm.

## Conventions

Import shared controls directly from `@kenn-io/kit-ui`. Keep local components
for Golinks behavior rather than pass-through wrappers. Preserve semantic tokens,
Svelte runes, typed props, callback props and snippets. Tailwind and SvelteKit are
not used. Keep the compact list, inline editor and phone support.

The app has no accounts or identity configuration. Never populate the application
with invented links. Preserve inputs and selected files after errors, focus
restoration, and navigable search/edit/help URL state. Use only the fixed charcoal
dark design, without a theme selector or saved appearance preferences. Keep pressed
controls stationary, input focus to a single border, and table hover borders unchanged.

Redirect paths bypass ServeMux cleaning. Preserve encoded destinations, repeated
slashes, query merging and Go template helper behavior. Templates expose `.Path`
and `.Now`. Count successful GET redirects, not HEAD requests or previews.
Short-name lookup ignores case and hyphens; collisions cannot overwrite data.
The normalized name `_` is reserved for internal routes; `api` is available.

SQLite uses one connection, rollback journaling and FULL synchronization; do not
enable WAL on network storage. JSONL/CSV backups import atomically and idempotently
through the browser, `/_/api/import` and CLI. Preserve timestamps and clicks.
Browser files are limited to 8 MiB combined; API import bodies to 16 MiB.

Local development uses `just api` and `just frontend`; `GOLINKS_API_TARGET`
overrides the frontend proxy target. Configuration uses `GOLINKS_*` environment
variables; the default database is `data/golinks.db`. `just build` produces one Go
executable containing the frontend. Keep databases, local config, live tooling,
snapshots, dependency caches, specs, plans and agent reports out of Git.

## Deployment

GitHub Actions validates pull requests and publishes public multi-platform images
to `ghcr.io/danielfrg/golinks`. Deployment configuration belongs in the consuming
infrastructure repository. All internal routes use `/_/`, including assets,
exports and health probes. Containers run as UID/GID 65532 with writable storage.

## Check

```zsh
go test -race ./... -count=1
go vet ./...
pnpm --dir frontend test
pnpm --dir frontend check
just build
```

Server lifecycle tests bind an ephemeral localhost port and require socket access.
For UI work, verify the actual embedded server in a browser at desktop and phone
widths, including forms, preview, search, navigation and error states.
