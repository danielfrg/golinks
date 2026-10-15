import test from "node:test";
import assert from "node:assert/strict";
import { createServer as createHTTPServer } from "node:http";
import { once } from "node:events";
import type { AddressInfo } from "node:net";
import { createServer as createViteServer } from "vite";

test("dev proxy preserves browser Host and Origin, including rejected foreign origins", async () => {
 const upstream = createHTTPServer((request, response) => {
  if (request.method === "POST" && request.headers.origin !== `http://${request.headers.host}`) {
   response.writeHead(403).end("Origin does not match the request host");
   return;
  }
  response.setHeader("Content-Type", "application/json");
  response.end(JSON.stringify({host: request.headers.host, origin: request.headers.origin, path: request.url}));
 });
 upstream.listen(0, "127.0.0.1");
 await once(upstream, "listening");
 const previousTarget = process.env.GOLINKS_API_TARGET;
 process.env.GOLINKS_API_TARGET = `http://127.0.0.1:${(upstream.address() as AddressInfo).port}`;
 let vite: Awaited<ReturnType<typeof createViteServer>> | undefined;
 try {
  const { default: config } = await import("../../vite.config.ts");
  vite = await createViteServer({
   ...config,
   configFile: false,
   logLevel: "silent",
   server: {...config.server, port: 0, strictPort: false, open: false},
  });
  await vite.listen();
  const host = `127.0.0.1:${(vite.httpServer!.address() as AddressInfo).port}`;
  const origin = `http://${host}`;
  const response = await fetch(`${origin}/_/api/resolve`, {
   method: "POST",
   headers: {Origin: origin, "Content-Type": "application/json", "X-Golink": "1"},
   body: JSON.stringify({long: "https://example.com/", path: "", query: ""}),
  });
  const body = await response.text();
  assert.equal(response.status, 200, body);
  assert.deepEqual(JSON.parse(body), {host, origin, path: "/_/api/resolve"});
  // A separate read verifies the export proxy also honors the configured target
  // and preserves the browser-facing request host.
  const exported = await fetch(`${origin}/_/export`, {headers: {Origin: origin}});
  assert.equal(exported.status, 200);
  assert.deepEqual(await exported.json(), {host, origin, path: "/_/export"});
  const foreign = await fetch(`${origin}/_/api/resolve`, {
   method: "POST",
   headers: {Origin: "https://other.test", "Content-Type": "application/json", "X-Golink": "1"},
   body: "{}",
  });
  assert.equal(foreign.status, 403);
 } finally {
  if (previousTarget === undefined) delete process.env.GOLINKS_API_TARGET;
  else process.env.GOLINKS_API_TARGET = previousTarget;
  await vite?.close();
  upstream.closeAllConnections();
  await new Promise<void>((resolve, reject) => upstream.close(error => error ? reject(error) : resolve()));
 }
});
