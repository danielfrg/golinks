# Golinks

A golinks service. Ready for self-hosting.

Go links are intuitive, easy-to-remember short links.
[Read more about them](https://golinks.github.io/golinks/).

![Golinks](docs/screenshots/homepage.png)

## Installation

### Docker

```zsh
docker run -it \
  -p 127.0.0.1:8080:8080 \
  -v golinks-data:/data \
  -e GOLINKS_PUBLIC_URL=http://localhost:8080 \
  ghcr.io/danielfrg/golinks:latest
```

- Open [http://localhost:8080](http://localhost:8080)
- The SQLite database is stored at `/data/golinks.db`
- The container runs as UID/GID 65532

Published images support Linux AMD64 and ARM64. `latest` tracks `main`.

### Standalone

Build a single executable with `just build` (see [Development](#development)),
then run:

```zsh
./golinks
```

Open [http://127.0.0.1:8080](http://127.0.0.1:8080).

## Usage

Create a link in the UI and `go/name` redirects to its destination.

Names are case-insensitive and ignore hyphens: `meeting-notes` and `MeetingNotes` are the
same link.

Destinations can also be Go URL templates. For example, `go/g/something`
searches Google:

```text
https://www.google.com/{{if .Path}}search?q={{QueryEscape .Path}}{{end}}
```

Use **Try a path** in the editor to preview a template. Add `+` to a short
link (`go/name+`) to open its details instead of redirecting. API documentation
is at `/_/api/docs`.

### Browser redirects

This works great with some of the [Redirect extensions](https://chromewebstore.google.com/detail/redirector/jegbdohdgebjljoljfeinojeobdabpjo?hl=en)

Config:

```
Redirect: http://go/*
To: https://go.example.net/$1
Example: http://go/page → https://go.example.net/page

## You might also need this in Chrome:
Redirect: ^https:\/\/www\.google\.com\/search\?.*q=go(?:%2F|\/)([^&]+).*$
To: https://go.example.net/$1
```

### Backups

Use **Export links** and **Import links** in the footer to move links between
instances as JSONL files. Or with the CLI:

```shell
export GOLINKS_DATABASE=./data/restored.db
golinks import --links /path/to/golinks.jsonl --stats /path/to/golinks-stats.csv
```

### MCP

Agents can manage links through the MCP endpoint at `/_/mcp` (Streamable HTTP).
Tools: `search_links`, `get_link`, `resolve_link`, `create_link`, and
`update_link`.

### Security

There is no authentication: anyone who can reach the service can manage links,
so run it on a trusted network or behind an authenticated proxy.

## Configuration

| Variable             | Default                                      |
| -------------------- | -------------------------------------------- |
| `GOLINKS_ADDR`       | `127.0.0.1:8080`                             |
| `GOLINKS_DATABASE`   | `./data/golinks.db`                          |
| `GOLINKS_PUBLIC_URL` | `http://127.0.0.1:8080`                      |
| `GOLINKS_READONLY`   | `false`                                      |
| `GOLINKS_API_TARGET` | `http://127.0.0.1:8080` (frontend dev proxy) |

For a reverse proxy, set `GOLINKS_PUBLIC_URL` to its public origin. Readiness
is `/_/readyz`; liveness is `/_/healthz`.

## Development

Requires Go 1.25, Node 26, pnpm 10, and just.

Start the API and Frontend in two terminals:

```shell
# Terminal 1:
just install
just api

# Terminal 2:
just frontend
```

Open [http://127.0.0.1:5173](http://127.0.0.1:5173). Vite proxies API requests
to the Go server. Run `just check` before sending changes. See `AGENTS.md` for
project structure and conventions.

## Build

```zsh
just build
```

```zsh
just image-check
```
