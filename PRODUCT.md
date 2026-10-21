# Product

## Users

People or small teams looking for a simple golinks service to build memorable
shortcuts and a directory for finding, creating, and maintaining their links
from a desktop or phone.

Assumes deployment on a trusted network; there is no auth.

## Product Purpose

Turn a short name into the right destination immediately.

Keep link management, Go URL templates, visit counts, and backups in one self-hosted executable.

Success means opening a familiar shortcut quickly and changing its destination.

## Brand Personality

Direct, compact, and calm. Names and destinations carry the interface.

### Tone

Do not be overly explicit about explaining every part of the UI
with text. It should be intuitive. Only add explanatory messages when it's needed.

## Anti-references

Avoid promotional dashboards, invented sample links, account-management flows,
and layouts that hide the directory behind a separate editing screen.

## Design Principles

- Put short names and destinations first.
- Keep the directory visible while editing or previewing a template.
- Preserve encoded paths, query parameters, timestamps, and backup data.
- Distinguish loading, no links, no search results, and request failures.
- Preserve unsaved input and return focus after save, cancel, or import.

## Agent Access

Agents use the same product through typed MCP tools: search, inspect, and
preview links without counting visits; create and update with the same
validation as the HTTP API. Read-only mode omits the write tools.
