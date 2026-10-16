# Golinks

Self-hosted short links with Go URL templates, SQLite storage, and an embedded
Svelte UI. Anyone who can reach the service can manage links; use it on a
trusted network or behind an authenticated proxy.

![Golinks](docs/screenshots/homepage.png)

## Installation

### Docker

```zsh
docker run -d --name golinks \
  -p 127.0.0.1:8080:8080 \
  -v golinks-data:/data \
  -e GOLINKS_PUBLIC_URL=http://localhost:8080 \
  ghcr.io/danielfrg/golinks:latest
```

Open [http://localhost:8080](http://localhost:8080). The container runs as
UID/GID 65532 and stores SQLite data at `/data/golinks.db`. Bind-mounted data
directories must be writable by that user. Run one instance per database.

### Standalone

Build from source using the [build instructions](#build), then run:

```zsh
./golinks
```

Open [http://127.0.0.1:8080](http://127.0.0.1:8080).

## Usage

Names are case-insensitive and ignore hyphens: `meeting-notes` and `MeetingNotes`
resolve to the same link. Extra path segments are appended to plain destinations.
Incoming query parameters are merged with destination parameters.

Destination templates use `.Path`, `.Now` (UTC), and `PathEscape`, `QueryEscape`,
`TrimPrefix`, `TrimSuffix`, `ToLower`,
`ToUpper`, and `Match`. For example:

```text
https://www.google.com/{{if .Path}}search?q={{QueryEscape .Path}}{{end}}
```

Use Try a path in the editor to preview a destination. Add `+` to a short link
to open its details. API documentation is at `/_/api/docs`.

### Backups

Use **Export links** and **Import links** in the footer to transfer JSONL backups
between instances. Links, templates, timestamps and total visits are preserved.
An optional CSV from `/_/export-stats` also preserves click history. Existing names
are skipped; invalid backups make no partial changes. Browser imports accept up
to 8 MiB of files combined. Larger backups use the CLI:

```zsh
GOLINKS_DATABASE=./data/restored.db ./golinks import \
  --links /path/to/golinks.jsonl --stats /path/to/golinks-stats.csv
```


## Configuration

| Variable                | Default                                      |
| ----------------------- | -------------------------------------------- |
| `GOLINKS_ADDR`          | `127.0.0.1:8080`                             |
| `GOLINKS_DATABASE`      | `./data/golinks.db`                          |
| `GOLINKS_PUBLIC_URL`    | `http://127.0.0.1:8080`                      |
| `GOLINKS_READONLY`      | `false`                                      |
| `GOLINKS_API_TARGET`    | `http://127.0.0.1:8080` (frontend dev proxy) |

For another local port:

```zsh
GOLINKS_ADDR=127.0.0.1:8082 GOLINKS_PUBLIC_URL=http://127.0.0.1:8082 ./golinks
```

Readiness is `/_/readyz`; liveness is `/_/healthz`. The `/_/` namespace
keeps operational endpoints separate from short links. Readiness checks database
access. For a reverse proxy, set `GOLINKS_PUBLIC_URL` to its public origin.

## Development

Requires Go 1.25, Node 26, pnpm 10, and just.

```zsh
just install
just api
```

In another terminal:

```zsh
just frontend
```

Open [http://127.0.0.1:5173](http://127.0.0.1:5173). Vite proxies API
requests to the Go server. Run `just check` for Go race tests, vet, frontend
tests, and Svelte/TypeScript checks.

The repositories share this structure:

```text
cmd/<app>/              Configuration and server lifecycle
internal/               Go application packages
internal/httpapi/       HTTP API and health probes
internal/web/           Embedded frontend and static serving
frontend/src/theme.css  Shared palette, typography, and font definitions
frontend/src/app.css    Application layouts and interactions
scripts/                Container smoke test
docs/screenshots/       README screenshots
```

The shared theme loads Inter and IBM Plex Mono from Google Fonts, with dark, light,
and system preferences. Dark is the default; the header theme control saves
your choice. Keep `frontend/src/theme.css` identical across these projects.
`PRODUCT.md` describes users and purpose; `DESIGN.md` records interface rules.
Local plans, agent reports, and design-tool state stay out of Git.

## Build

```zsh
just build
```

The executable embeds the built frontend and needs no Node runtime. Build and
smoke-test a local container with:

```zsh
just image-check
```

Published images support Linux AMD64 and ARM64. `latest` tracks `main`;
`sha-<commit>` tags and image digests can pin a deployment.
