# src/CLAUDE.md

Frontend guidance for the Vue 3 SPA in this directory.

---

## Vue Test Patterns

- `vitest.config.ts` scopes to `src/**/*.test.ts` — do not change this or Playwright BDD specs get picked up
- Always create a **fresh router per test** via a factory function — sharing a router instance causes watcher accumulation and phantom fetch calls
- Stub `fetch` in `beforeEach` with `vi.stubGlobal('fetch', vi.fn()...)` and restore with `vi.restoreAllMocks()` in `afterEach`
- Test files live in `__tests__/` subdirectories alongside the code they test

---

## Design Language

**Theme:** Industrial Black & Gold Brutalism — matchday program crossed with a construction-site bulletin board.

**Tokens live in `src/style.css`** — do not duplicate values here; the CSS file is the source of truth.

Key patterns:
- 3px gold left border on hovered/predicted match cards
- Score inputs: 52×52px, `DM Mono`, gold text, dark background, focus glows gold
- Locked state: blinking `▊` indicator in `--danger` red
- Noise texture overlay on `body::before`; gold stripe on `body::after` (top of viewport, fixed)

Typography: `Bebas Neue` (headings) · `DM Mono` (scores/metadata) · `Barlow` (body copy)
