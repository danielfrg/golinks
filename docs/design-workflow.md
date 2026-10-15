# Golinks design workflow

Read `PRODUCT.md`, `DESIGN.md` and the target source before UI changes.
Use the optional local Impeccable skill for review; it adds no production
dependency. Preserve the compact Golinks workflow and Kit UI conventions.

## Shared controls and tokens

Import controls directly from `@kenn-io/kit-ui`. Add local components only for
app-specific behavior. Kit UI owns shared controls and light/dark tokens;
`frontend/src/app.css` maps semantic aliases and owns application layout.
Use Kit UI default control sizes, surfaces, typography and base theme; do not
restyle shared components. Leave text selection to the browser.
The dependency is pinned to an upstream commit, as in CI View. Theme state is
initialized in `frontend/src/main.ts` with storage key `golinks-theme`.

## Local review

Start `just api` and `just frontend` in separate terminals. Open
`http://127.0.0.1:5173/_/assets/`; verify that the page says Golinks.

```zsh
"$HOME/.codex/skills/impeccable/scripts/impeccable" context --target frontend/src/App.svelte
"$HOME/.codex/skills/impeccable/scripts/impeccable" live --target frontend/src/App.svelte
```

The shared Live configuration targets Vite's index.html. Document accepted
system changes in DESIGN.md. Stop Live Mode from the frontend root when finished:

```zsh
"$HOME/.codex/skills/impeccable/scripts/impeccable" live-server stop
```

Keep helper injections, sessions, screenshots, variants and caches out of builds
and Git. Never populate the application with invented links.

## Verification

Run the checks in AGENTS.md. Build with `just build`, then verify the actual
embedded server at desktop and 320/360px phone widths in light and dark themes.
Check forms, preview, search, sort, browser navigation, focus restoration,
read-only mode, import, long destinations, and loading/empty/error states.
Use intercepted browser responses for fixture data; do not seed a database.
