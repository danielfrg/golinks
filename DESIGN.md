# Golinks design

## Foundation

Charcoal Index is the approved OpenCode-inspired direction drawn from the local
Homelab Launcher: an operational short-link directory with IBM Plex Mono,
charcoal surfaces, square edges and restrained amber feedback. A compact, full-width navigation header
precedes search, New link and Sort above one table. Keep the directory
visible while editing. No sidebar, summary tiles or decorative dashboard cards.

Import Button, TextInput, SearchInput, SelectDropdown, EmptyState, IconButton,
Table and TableHeaderCell directly from `@kenn-io/kit-ui`. Keep local
components for Golinks behavior rather than pass-through wrappers. Plain CSS in
`frontend/src/app.css` supplies the approved theme and directory/editor layout;
Tailwind is not used. Native textarea and file inputs use the same semantic tokens.

The fixed palette below is defined by `:root.dark`. The document always has the
`dark` class. The legacy token name `--accent-blue` carries amber.

| Kit UI token | Charcoal |
| --- | --- |
| `--bg-primary` | `#161616` |
| `--bg-surface` | `#242424` |
| `--bg-inset` | `#202020` |
| `--bg-surface-hover` | `#303030` |
| `--text-primary` | `#ffffff` |
| `--text-secondary` | `#b8b8b8` |
| `--text-muted` | `#a3a3a3` |
| `--border-default` | `#4a4a4a` |
| `--border-muted` | `#3a3a3a` |
| `--accent-blue` | `#e8a44b` |
| `--accent-red` | `#ff9b9b` |

Golinks aliases map `--background` to `--bg-primary`, `--foreground` to
`--text-primary`, `--card` to `--bg-surface`, `--muted` to `--bg-inset`,
`--muted-foreground` to `--text-secondary`, `--border` and `--input` to
`--border-default`, `--ring` to `--accent-blue`, and `--destructive` to
`--accent-red`. All three radius tokens are `0px`; `--radius` uses `--radius-lg`.
`--shadow-raised`, `--shadow-sunken` and `--shadow-panel` are `none`.

## Layout and type

The full-width header matches OpenCode Data's navigation bar: sticky at the top,
72px tall with 20px vertical padding, a charcoal background and no bottom
border. Its horizontal padding matches the centered 1080px content column,
including the 20px gutters at phone widths. The Golinks wordmark uses uppercase 16px medium type with a 20px line
height. The header gap is 32px, reducing to 16px at 760px. Search, New link and
Sort remain above the table rather than inside the navigation bar.

Below it, center a 1080px content column with shell padding
`32px max(24px, calc((100vw - 1080px)/2)) 120px`. At 760px the shell becomes
`24px 20px 110px`. The shell fills the viewport below the 72px header.

Self-host IBM Plex Mono regular (400) and medium (500) TTF files from
`frontend/src/lib/assets/fonts`, retaining the OFL notices. Both `--font-sans`
and `--font-mono` use IBM Plex Mono with `ui-monospace, SFMono-Regular, Menlo,
monospace` fallbacks. Body sizing uses `--font-size-root`; shared controls retain
Kit UI md sizing, with no app `--kit-control-height`. Buttons are naturally about
29px tall, TextInput is 28px and SelectDropdown is 26px on desktop. The shared
`--font-weight-medium` is 400. Main short paths use 14px and footer actions 11px.

The toolbar wraps search onto a full row at narrow widths (760px and 600px rules),
with New link and Sort beneath it. Most used is the initial sort; Name is
available. Sorting belongs in this toolbar; table headers are not interactive.
Table body cells have 14px padding and 13px text. Headers have `12px 14px`
padding, 11px text, .04em tracking and weight 400. Visits use tabular numerals.
The wrapper has a static 1px default border and scrolls horizontally on phones,
including 320px, without widening the page. Preserve default stripes and hover
backgrounds without hover outlines or border color changes. Read-only mode omits
the Actions column and editing/import.

The inline editor uses a square surface with a 1px default border and 24px padding
(20px at 760px; 16px at 400px). Fields stack at 760px; preview fields and actions
wrap at 400px. Long destinations and template results wrap within their containers.

## Interaction and states

The site uses only the charcoal dark design. There is no theme selector, theme
store initialization or saved appearance preference. System appearance and old
preferences do not change the palette. Native browser fields use the dark color
scheme, and text selection keeps the browser highlight.

Use a single visible focus border for text fields: TextInput's outer border turns
amber through `--ring`, while its outline, shadow and inner input outline remain
absent. Its inner input has no border; placeholders use `--text-secondary`.
Native textarea and file inputs also recolor their existing border on focus,
without another outline or shadow. Other controls retain visible Kit UI focus.
Buttons retain amber hover borders. `--press-transform: none` prevents buttons
and dropdown triggers moving downward when pressed. Do not add row hover border
changes or arrow animations. Respect reduced motion; preview disclosure rotation
is the existing 150ms transition.

Search/edit/help state lives in URL query parameters and survives refresh and
browser navigation. The inline editor shows immutable names when editing,
destination, Save/Cancel and optional path/query preview. Errors stay near the
form and preserve inputs. Delete confirms within the editor. Save restores focus
to Edit; closing a new editor restores focus to New link. An unknown shortcut
offers creation without automatically creating data.

Import accepts exported JSONL and optional click-history CSV. Preserve selected
files after errors, report imported/skipped counts, reload after success and
restore focus to Import links. Footer actions wrap on phones: Template help,
Export links, Import links and API docs (`/_/api/docs`). Distinguish loading,
no links, no results and failed requests; keep the previous list after refresh
errors and expose retry. Never add a Refresh action or invented links.
Impeccable remains local tooling; follow `docs/design-workflow.md`.
