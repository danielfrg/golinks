# Golinks design workflow

UI review uses the optional local Impeccable workflow.
Impeccable is installed at `~/.codex/skills/impeccable` and adds no production
dependency. Read PRODUCT.md, DESIGN.md and the target source before UI changes.

Start `just api` and `just frontend` in separate terminals. Open
`http://127.0.0.1:5173/_/assets/`; verify that the page says Golinks.

```zsh
"$HOME/.codex/skills/impeccable/scripts/impeccable" context --target frontend/src/App.svelte
"$HOME/.codex/skills/impeccable/scripts/impeccable" live --target frontend/src/App.svelte
```

The frontend's shared Live configuration targets Vite's index.html. Keep
alternatives within the inherited neutral foundation, and document accepted
changes in DESIGN.md. Stop Live Mode from the frontend root when finished:

```zsh
"$HOME/.codex/skills/impeccable/scripts/impeccable" live-server stop
```

Confirm index.html has no helper injection before production builds. Session
state, captures, preview variants, receipts and caches are ignored. Preserve the
chosen design, remove unused variants and keep the machine-readable sidecar
aligned with the implemented tokens/components.

Run the checks in AGENTS.md. Review the actual browser at desktop and 320/360 px,
including long destinations, keyboard focus, failed forms and loading/empty states.
