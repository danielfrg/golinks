# Golinks

<!-- impeccable:product-schema 1 -->

## Platform

web

## Stack

Go, Huma, SQLite, Svelte 5, TypeScript and Tailwind 4. The built frontend is
embedded in one Go executable.

## Users

Daniel and others on the network who use memorable shortcuts from a
desktop or phone, and maintain the destinations through a compact browser UI.

## Product Purpose

Provide memorable link shortcuts with an independently built service.
Success means redirects and Go URL templates work, and links are easy to find,
create, edit.

## Operating Context

## Capabilities and Constraints

- Redirect immediately using HTTP 302; support appended paths, merged queries,
  aliases and Go destination templates.
- Resolve short names without case sensitivity or hyphens.
- Store links and click counts in SQLite. Export and import the app's own backups.
- Serve the UI and API from one executable.
- Show loading, empty, failed request and validation states; preserve unsaved
  input after a failed operation.

## Brand Personality

Focused, clean, friendly and minimal.

## Design Principles

- Put familiar short names and actual destinations first.
- Keep management compact and avoid overexplaining the UI.
- Show the result of a template through an explicit preview.
- Make every action usable with a keyboard and on a phone.
- Use actual data; never populate the service with invented links.
