# Design System

Golinks is a simple self-hosted application.

The interface should use a small consistent vocabulary rather than
one-off styling for each new control.

## Component Rules

Theme and component styles come from the shared design system, synced into
`frontend/src/design-system/` (`theme/` and `custom/`). Never edit that folder by
hand: changes are made upstream and synced back in.

- Use Svelte 5 runes, typed props, and callback props.
- Keep desktop and phone layouts, keyboard focus, and recoverable error states.
- Import controls from `@ui-sdk/svelte/<component>` and shared custom components
  (Table, ThemeToggle) from `src/design-system/custom/`.
- Use semantic theme tokens (`--ink`, `--muted`, `--line`, `--space-*`) and keep the
  library's default control sizes.
- `src/app.css` holds only golinks layout and app-specific pieces; it adjusts a library
  control's layout (width, spacing) at most, never its look.
- Keep local components for app-specific behavior, not pass-through wrappers.
- Give fields one visible focus treatment and keep pressed controls stationary.
- Preserve `scrollbar-gutter: stable`.

## Appearance and Accessibility

Start in dark mode and respect explicit light or system preferences. The header
ThemeToggle cycles light, dark, and system and saves the choice under the shared
`color-scheme` key. Keep focus visible, controls labeled, and reduced-motion
preferences respected. Phone layouts must not widen the page.

## Application Patterns

Use an 88px header and begin the content at 128px on desktop, or 110px at 640px and below.
Align the title and content inside the same centered 1344px container.

Keep the compact wordmark header, search/New link/Sort toolbar, and one link
table. Short names use monospace; destinations use interface text. Preserve the
initial Most used sort and the optional Name sort. The table may scroll inside
its wrapper on small screens.

Keep the inline editor, template preview, import, and help near the directory.
Errors preserve input and selected files. Save or close restores focus to the
relevant action. Search, edit, and help state survive reload and browser history.
Read-only mode omits editing and import actions. Do not add a Refresh control.
