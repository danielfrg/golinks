# Design System

## Foundation

Use the shared desktop workspace theme in `frontend/src/theme.css`. Keep that
file identical across Golinks, CI View, and Homelab Launcher. Application layout
and interaction rules belong in `frontend/src/app.css`.

Inter is for interface text; IBM Plex Mono is for wordmarks, identifiers, and
technical data. Use neutral surfaces, subtle one-pixel borders, rounded controls,
and no decorative shadows. Focus is blue in light mode and amber in dark mode.
Preserve `scrollbar-gutter: stable`.

## Component Rules

- Import shared controls directly from `@kenn-io/kit-ui`.
- Use semantic theme tokens and keep the library's default control sizes.
- Keep local components for app-specific behavior, not pass-through wrappers.
- Use Svelte 5 runes, typed props, and callback props.
- Give fields one visible focus treatment and keep pressed controls stationary.

## Appearance and Accessibility

Start in dark mode and respect explicit light or system preferences. The header
ThemeToggle saves the choice under the app's own storage key. System mode follows
the operating system. Keep focus visible, controls labeled, and reduced-motion
preferences respected. Phone layouts must not widen the page.

## Application Patterns

Use an 88px header and begin the content at 128px on desktop, or 110px at
640px and below, matching CI View and Homelab Launcher. Align the title and
content inside the same centered 1344px container.

Keep the compact wordmark header, search/New link/Sort toolbar, and one link
table. Short names use monospace; destinations use interface text. Preserve the
initial Most used sort and the optional Name sort. The table may scroll inside
its wrapper on small screens.

Keep the inline editor, template preview, import, and help near the directory.
Errors preserve input and selected files. Save or close restores focus to the
relevant action. Search, edit, and help state survive reload and browser history.
Read-only mode omits editing and import actions. Do not add a Refresh control.

Use `PRODUCT.md` for user intent and `README.md` for setup and commands.
