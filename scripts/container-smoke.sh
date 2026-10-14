#!/usr/bin/env bash
set -euo pipefail
image_ref=${1:?usage: container-smoke.sh IMAGE}
smoke_dir=$(mktemp -d)
smoke_suffix=$(basename "$smoke_dir" | tr -cd '[:alnum:]')
containers=()
volumes=()
cleanup() {
  if ((${#containers[@]})); then docker rm -f "${containers[@]}" >/dev/null 2>&1 || true; fi
  if ((${#volumes[@]})); then docker volume rm "${volumes[@]}" >/dev/null 2>&1 || true; fi
  rm -rf "$smoke_dir"
}
trap cleanup EXIT
[[ $(docker image inspect "$image_ref" --format '{{.Config.User}}') == 65532:65532 ]]
new_volume() {
  smoke_volume="golinks-smoke-${smoke_suffix}-$1"
  volumes+=("$smoke_volume")
  docker volume create "$smoke_volume" >/dev/null
}
start_server() {
  local server_name="golinks-smoke-${smoke_suffix}-$1" data_volume=$2
  local listen_port
  listen_port=$(python3 -c 'import socket;s=socket.socket();s.bind(("127.0.0.1",0));print(s.getsockname()[1]);s.close()')
  smoke_base="http://127.0.0.1:$listen_port"
  containers+=("$server_name")
  docker run -d --name "$server_name" --read-only --tmpfs /tmp \
    -p "127.0.0.1:$listen_port:8080" -v "$data_volume:/data" \
    -e "GOLINKS_PUBLIC_URL=$smoke_base" "$image_ref" >/dev/null
  for _ in {1..60}; do
    if curl --fail --silent "$smoke_base/_/readyz" >/dev/null; then return; fi
    sleep .2
  done
  docker logs "$server_name"
  return 1
}
request() {
  curl --fail --silent -H 'Content-Type: application/json' -H 'X-Golink: 1' "$@"
}
new_volume source; source_volume=$smoke_volume
start_server source "$source_volume"; source_base=$smoke_base
curl --fail --silent "$source_base/_/healthz" >/dev/null
curl --fail --silent "$source_base/" > "$smoke_dir/index.html"
python3 - "$smoke_dir/index.html" > "$smoke_dir/asset" <<'PY'
import re,sys
text=open(sys.argv[1]).read()
assert '<title>Golinks</title>' in text
asset=re.search(r'src="(/_/assets/[^\"]+\.js)"',text)
assert asset, 'embedded frontend asset missing'
print(asset.group(1))
PY
curl --fail --silent "$source_base$(cat "$smoke_dir/asset")" >/dev/null
curl --fail --silent "$source_base/_/api/docs" >/dev/null
curl --fail --silent "$source_base/_/api/openapi.json" > "$smoke_dir/schema.json"
python3 - "$smoke_dir/schema.json" <<'PY'
import json,sys
schema=json.load(open(sys.argv[1]));assert '/_/api/links' in schema['paths'];assert '/_/api/import' in schema['paths']
PY
request -d '{"short":"wiki","long":"https://example.org/wiki/"}' "$source_base/_/api/links" >/dev/null
request -d '{"short":"search","long":"https://example.org/?q={{QueryEscape .Path}}"}' "$source_base/_/api/links" >/dev/null
request -d '{"short":"alias","long":"/wiki"}' "$source_base/_/api/links" >/dev/null
request -d '{"long":"https://example.org/?q={{QueryEscape .Path}}","path":"a b","query":""}' "$source_base/_/api/resolve" > "$smoke_dir/preview.json"
python3 - "$smoke_dir/preview.json" <<'PY'
import json,sys
assert json.load(open(sys.argv[1]))['url']=='https://example.org/?q=a+b'
PY
for uri in /wiki /wiki /search/hello /alias; do
  [[ $(curl --silent -o /dev/null -w '%{http_code}' "$source_base$uri") == 302 ]]
done
curl --fail --silent "$source_base/_/api/links" > "$smoke_dir/expected.json"
docker rm -f "golinks-smoke-${smoke_suffix}-source" >/dev/null
start_server source "$source_volume"; source_base=$smoke_base
curl --fail --silent "$source_base/_/api/links" > "$smoke_dir/restarted.json"
cmp "$smoke_dir/expected.json" "$smoke_dir/restarted.json"
curl --fail --silent "$source_base/_/export" > "$smoke_dir/golinks.jsonl"
curl --fail --silent "$source_base/_/export-stats" > "$smoke_dir/golinks-stats.csv"
python3 - "$smoke_dir" <<'PY'
import json,sys,pathlib
p=pathlib.Path(sys.argv[1]);links=(p/'golinks.jsonl').read_text();stats=(p/'golinks-stats.csv').read_text()
(p/'import.json').write_text(json.dumps({'links':links,'stats':stats}))
(p/'invalid.json').write_text(json.dumps({'links':json.dumps({'short':'fresh','long':'https://example.org/'})+'\n{broken'}))
PY
new_volume target; target_volume=$smoke_volume
start_server target "$target_volume"; target_base=$smoke_base
request --data-binary "@$smoke_dir/import.json" "$target_base/_/api/import" >/dev/null
request --data-binary "@$smoke_dir/import.json" "$target_base/_/api/import" > "$smoke_dir/repeated.json"
python3 - "$smoke_dir/repeated.json" <<'PY'
import json,sys
result=json.load(open(sys.argv[1]));assert result['imported']==0;assert result['skipped']==3
PY
[[ $(curl --silent -o /dev/null -w '%{http_code}' -H 'Content-Type: application/json' -H 'X-Golink: 1' --data-binary "@$smoke_dir/invalid.json" "$target_base/_/api/import") == 422 ]]
curl --fail --silent "$target_base/_/api/links" > "$smoke_dir/restored.json"
cmp "$smoke_dir/expected.json" "$smoke_dir/restored.json"
curl --fail --silent "$target_base/_/export-stats" > "$smoke_dir/restored.csv"
cmp "$smoke_dir/golinks-stats.csv" "$smoke_dir/restored.csv"
new_volume cli; cli_volume=$smoke_volume
chmod 755 "$smoke_dir"; chmod 644 "$smoke_dir/golinks.jsonl" "$smoke_dir/golinks-stats.csv"
docker run --rm --read-only --tmpfs /tmp -v "$cli_volume:/data" -v "$smoke_dir:/backup:ro" "$image_ref" import --links /backup/golinks.jsonl --stats /backup/golinks-stats.csv >/dev/null
start_server cli "$cli_volume"
curl --fail --silent "$smoke_base/_/api/links" > "$smoke_dir/cli.json"
cmp "$smoke_dir/expected.json" "$smoke_dir/cli.json"
if docker run --rm --read-only -v "$smoke_dir:/backup:ro" -e GOLINKS_DATABASE=/unwritable/golinks.db "$image_ref" import --links /backup/golinks.jsonl 2>/dev/null; then
  echo 'Unwritable database should fail' >&2; exit 1
fi
printf 'Container verified: embedded UI, API, redirects, nonroot storage, restart persistence and backup round trip.\n'
