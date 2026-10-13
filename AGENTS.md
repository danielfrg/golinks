## General

Golinks is a standalone Go/Svelte redirect service.

- Do not edit README.md unless the user explicitly requests it
- Read PRODUCT.md and DESIGN.md before changing the frontend; preserve the
  semantic tokens and local component conventions.
- Use Impeccable for UI review and follow docs/design-workflow.md.

Local development uses `just api` and `just frontend`. `just build` produces
one Go executable containing the frontend. Configuration uses `GOLINKS_*`
environment variables; the default local database is `data/golinks.db`.

## Capabilities and Constraints

- Go/net/http and Huma serve redirects, the UI and the JSON API under `/_/api`.
  The normalized short name `_` is reserved for internal routes; `api` is available.
- SQLite uses one connection, rollback journaling and FULL synchronization;
  do not enable WAL on network storage.
- Redirect paths bypass ServeMux cleaning. Preserve encoded destinations,
  repeated slashes, query merging and Go template helper behavior.
- Short-name lookup ignores case and hyphens; collisions cannot overwrite data.
- The app has no accounts or identity configuration. Templates expose `.Path`
  and `.Now`.
- Count successful GET redirects, not HEAD requests or previews.
- The app's JSONL/CSV backups import atomically and idempotently through the
  browser, `/_/api/import` and CLI. Preserve timestamps and clicks. Browser files
  are limited to 8 MiB combined; API import bodies to 16 MiB.
- Svelte 5, TypeScript, Tailwind 4 and local shadcn-svelte components render a
  compact list and inline editor. Preserve inputs after errors and support phones.
- Never populate the application with invented links
- Keep live tooling, databases, snapshots and dependency caches out of Git.

## Deployment


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
