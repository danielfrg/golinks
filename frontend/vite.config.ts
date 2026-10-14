import { fileURLToPath, URL } from "node:url";
import { svelte } from "@sveltejs/vite-plugin-svelte";
import tailwindcss from "@tailwindcss/vite";
import { defineConfig } from "vite";

export default defineConfig({
 plugins:[tailwindcss(),svelte()],
 base:"/_/assets/",
 resolve:{alias:{$lib:fileURLToPath(new URL("./src/lib",import.meta.url))}},
 server:{host:"127.0.0.1",port:5173,strictPort:true,proxy:{
  "/_/api":process.env.GOLINKS_API_TARGET ?? "http://127.0.0.1:8080",
  "/_/export":process.env.GOLINKS_API_TARGET ?? "http://127.0.0.1:8080",
 }},
});
