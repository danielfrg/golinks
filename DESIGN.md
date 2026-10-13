---
name: Golinks
description: Compact homelab short-link manager
colors:
  background: "oklch(0.985 0 0)"
  foreground: "oklch(0.145 0 0)"
  card: "oklch(1 0 0)"
  primary: "oklch(0.205 0 0)"
  primary-foreground: "oklch(0.985 0 0)"
  muted: "oklch(0.97 0 0)"
  muted-foreground: "oklch(0.45 0 0)"
  destructive: "oklch(0.577 0.245 27.325)"
  border: "oklch(0.922 0 0)"
  ring: "oklch(0.556 0 0)"
typography:
  headline:
    fontFamily: "ui-sans-serif, system-ui, sans-serif"
    fontSize: "30px"
    fontWeight: 600
    lineHeight: "1.2"
    letterSpacing: "-0.025em"
  body:
    fontFamily: "ui-sans-serif, system-ui, sans-serif"
    fontSize: "14px"
    fontWeight: 400
    lineHeight: "1.5"
rounded:
  md: "8px"
  lg: "10px"
spacing:
  sm: "12px"
  md: "20px"
  lg: "24px"
components:
  button-primary:
    backgroundColor: "{colors.primary}"
    textColor: "{colors.primary-foreground}"
    rounded: "{rounded.md}"
    padding: "0 16px"
    height: "44px"
  button-outline:
    backgroundColor: "transparent"
    textColor: "{colors.foreground}"
    rounded: "{rounded.md}"
    height: "44px"
  input:
    backgroundColor: "{colors.card}"
    textColor: "{colors.foreground}"
    rounded: "{rounded.md}"
    height: "44px"
---

# Design System: Golinks

## Overview

**Creative North Star: "Compact short-link manager"**

Read PRODUCT.md before changing the UI. Help people find and maintain a memorable
link quickly. Preserve the neutral light foundation shared by Budget and Homelab
Dashboard; this is a compact operational tool.

## Colors

Use the inherited neutral light surfaces and dark text. Destructive red is
reserved for errors and deletion. Semantic values are defined above and in CSS.

## Typography

Svelte 5, Tailwind 4, system sans-serif and local shadcn-svelte Button/Input
components. Semantic colors and radius live in frontend/src/app.css. Use these
tokens consistently instead of inventing component palettes. Monospace is used
only for short-link paths, destination templates and code examples.

## Layout

The page has a Golinks heading, search and New link action. A single bordered list
shows short name, destination, visit count and Edit. Rows become stacked on phones.
The alphabetical/most-used selector aligns right above the directory; omit a
directory subtitle. Keep the read-only indicator when applicable. Export links,
Import links, Template help and API docs are the footer actions. API docs links to `/_/api/docs`.
Do not add a Refresh action, summary
tiles or example data.

## Elevation & Depth

Flat surfaces separated by quiet 1 px borders. Do not add decorative shadows.

## Shapes

List/editor corners use 10 px; inputs and buttons use 8 px.

## Components

The inline editor shows immutable short names when editing, a destination field,
Save/Cancel controls and an optional path/query preview. Errors stay near the
relevant form. Delete asks for confirmation within that editor. A successful save
returns to the list and restores focus. Keep the list visible during editing.

Import opens an inline form for the exported JSONL and an optional click-history
CSV. Keep selected files after errors, report imported/skipped counts, reload the
list after success and restore focus to Import links. Hide Import in read-only
mode. The footer wraps on narrow screens.

## Do's and Don'ts

- Use labeled native controls, visible focus, 44 px primary touch controls,
  restrained hover feedback and reduced-motion support.
- Wrap long URLs/templates without page overflow at 320 px.
- Distinguish loading, no links, no results and request errors.
- Preserve form values after a failure and preserve the list after refresh errors.
- Query parameters keep search/edit/help state navigable and shareable.
- An unknown shortcut offers creation without automatically creating data.

Impeccable is local development tooling, with no production dependency. Follow
docs/design-workflow.md. Keep live sessions, captures and helper injections out
of builds. The sidecar describes this app's own foundation, not Budget data.
