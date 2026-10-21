# Design System

Golinks is a simple self-hosted application.

The interface should use a small consistent vocabulary rather than
one-off styling for each new control.

## Component Rules

Use the shared desktop workspace theme in `frontend/src/theme.css`. We keep that file identical across the other self-hosted apps, so changes there should be generic and minimal.
Application layout and interaction rules belong in `frontend/src/app.css`.

- Use Svelte 5 runes, typed props, and callback props.
- Keep desktop and phone layouts, keyboard focus, and recoverable error states.
- Load Inter and IBM Plex Mono from Google Fonts with `display=swap`. Do not
embed font binaries in the repository or production assets.
- Inter is for interface text; IBM Plex Mono is for wordmarks, identifiers, and
technical data. Use neutral surfaces, subtle one-pixel borders, rounded controls,
and no decorative shadows.
- Focus is blue in light mode and amber in dark mode.
- Preserve `scrollbar-gutter: stable`.
- Import shared controls directly from `@kenn-io/kit-ui`.
- Use semantic theme tokens and keep the library's default control sizes.
- Keep local components for app-specific behavior, not pass-through wrappers.
- Give fields one visible focus treatment and keep pressed controls stationary.

## Appearance and Accessibility

Start in dark mode and respect explicit light or system preferences. The header
ThemeToggle saves the choice under the app's own storage key. System mode follows
the operating system. Keep focus visible, controls labeled, and reduced-motion
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
