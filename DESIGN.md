# Golinks design

## Foundation

Operate mode: a compact short-link manager with system typography, neutral
surfaces and restrained boundaries. The header holds the title and theme toggle.
Search, New link and Sort share a toolbar directly above the Kit UI table. Keep the list visible while editing. No sidebar,
summary tiles or decorative dashboard cards.

Kit UI (`@kenn-io/kit-ui`) supplies plain CSS tokens and reusable controls.
The dependency is pinned to the same upstream GitHub commit as CI View in
`frontend/package.json`. `frontend/src/app.css` maps its semantic palette to
Golinks aliases and owns the directory/editor layout. Tailwind is not used.
Light/dark mode follows the system initially; the header ThemeToggle persists
an explicit choice under `golinks-theme`.

Import Kit UI controls directly: Button, TextInput, SearchInput, SelectDropdown,
EmptyState, IconButton, Table, TableHeaderCell and ThemeToggle. Keep local components for Golinks behavior, rather
than pass-through aliases. Native textarea and file inputs use semantic fields
and the same visible focus. Destructive color is reserved for errors/deletion.
Use `--background`, `--foreground`, `--card`, `--muted`, `--muted-foreground`,
`--border`, `--input`, `--ring` and `--destructive` aliases mapped to Kit UI tokens.

## Layout and type

Content is capped at 1040px with 28px desktop and 16–20px phone gutters.
Kit UI Table rows align short name, destination, visits and Edit. TableHeaderCell
supplies column headers, including numeric alignment for Visits. Keep the default
striped rows, hover styling, cell spacing and typography. On narrow screens the
Kit UI wrapper scrolls horizontally without widening the page.
The directory toolbar keeps search and New link on the left and Sort on the right;
below 600px search takes a full row with New link and Sort beneath it. Forms stack fields.
At 400px preview fields and actions wrap. Long URLs/templates wrap without
page overflow at 320px.

Use the system sans-serif stack, a 1.5rem semibold page heading and compact
supporting labels. Use monospace only for paths, templates and code examples.
Visit numerals are tabular. Quiet 1px borders separate flat surfaces without
shadows. Editor corners use Kit UI radius tokens. Shared controls use
Kit UI default sizes, typography, focus styles and surfaces without app overrides.
The sort selector defaults to Most used and aligns right above
the directory; sorting stays in this toolbar, so table headers are not interactive.
Read-only mode keeps its indicator and omits the Actions column and editing/import.

## Interaction and states

Labeled Kit UI search and keyboard-accessible sort controls share semantic
fields and visible focus. Search/edit/help state lives in URL query parameters
and survives refresh and browser navigation. Theme selection uses the shared
Kit UI store and default base theme; native browser fields follow its light/dark
color scheme. Text selection uses the browser default highlight. Do not override
Kit UI component sizes, placeholders, surfaces or palette in application CSS.

The inline editor shows immutable names when editing, destination, Save/Cancel,
and optional path/query preview. Errors stay near the form and preserve inputs.
Delete asks for confirmation within the editor. Save returns to the list and
restores focus to Edit; closing a new editor restores focus to New link.
An unknown shortcut offers creation without automatically creating data.

Import opens an inline form for exported JSONL and optional click-history CSV.
Keep selected files after errors, report imported/skipped counts, reload after
success and restore focus to Import links. Footer actions wrap on phones:
Template help, Export links, Import links and API docs (`/_/api/docs`).

Distinguish loading, no links, no results and failed requests. Preserve the
previous list after refresh errors and expose retry. Never add a Refresh action
or invented links. Respect reduced motion. Impeccable is local tooling with no
production dependency; follow `docs/design-workflow.md`.
