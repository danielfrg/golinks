# Golinks

Self-hosted short links with Go URL templates and SQLite storage.

## Register

product

## Users and purpose

A self-hoster and others on their trusted network open memorable shortcuts and
maintain their destinations from a desktop or phone. Success means reaching the
right destination quickly and finding, creating or editing links in a compact UI.
Anyone who can reach the service can manage links; there are no accounts or
identity settings. The Go executable embeds the Svelte frontend.

## Principles

- Put familiar short names and actual destinations first; never seed invented links.
- Redirect immediately with HTTP 302. Preserve appended paths, merged queries,
  aliases and Go destination templates.
- Resolve names without case sensitivity or hyphens; prevent collisions.
- Store links and successful GET click counts in SQLite.
- Keep management compact, with one list and an inline editor that leaves it visible.
- Show the result of a template through an explicit preview without counting a visit.
- Preserve unsaved input and selected files after failed operations.
- Keep loading, empty lists, no search results and failed requests distinct.
- Search, edit and help state survive refresh and browser navigation.
- Make actions usable with a keyboard and on a phone in the fixed charcoal design.
- Export and import the app's own backups atomically, preserving timestamps,
  total visits and optional click history.
- Serve the UI and API from one standalone executable and publish container images
  with persistent SQLite storage.
