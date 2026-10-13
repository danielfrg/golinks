# Golinks

A Go link service.

## Local development

Requires Go 1.25, Node 26, pnpm 10 and just:

```zsh
pnpm --dir frontend install --frozen-lockfile
```

Start the API in one terminal:

```zsh
just api
```

Start the UI in another:

```zsh
just frontend
```

Open [http://127.0.0.1:5173/](http://127.0.0.1:5173/).

## Build

Build one executable containing the frontend:

```zsh
just build

./golinks
```

Open [http://127.0.0.1:8080](http://127.0.0.1:8080)

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

## Links and templates

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

## Deployment

Homelab infrastructure lives in `~/code/danielfrg/homelab`. This milestone runs
locally; container builds, release automation and homelab cutover come later.
