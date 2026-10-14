set dotenv-load := true

api:
    go run ./cmd/golinks

frontend:
    pnpm --dir frontend dev

check:
    go test ./...
    go vet ./...
    pnpm --dir frontend test
    pnpm --dir frontend check

build:
    pnpm --dir frontend install --frozen-lockfile
    pnpm --dir frontend test
    pnpm --dir frontend check
    pnpm --dir frontend build
    rm -rf internal/web/dist
    mkdir -p internal/web/dist
    touch internal/web/dist/.keep
    cp -R frontend/dist/. internal/web/dist/
    test -s internal/web/dist/index.html
    go test ./...
    go vet ./...
    CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o golinks ./cmd/golinks

image:
    docker build -t golinks:local .

image-check: image
    bash scripts/container-smoke.sh golinks:local
