# CLAUDE.md

Guidance for Claude Code when working in this repository.
See also: [ARCHITECTURE.md](ARCHITECTURE.md) for system design, [BACKLOG.md](BACKLOG.md) for what's next.

---

## Development Approach

### BDD Dual-Loop TDD

**This is the highest-priority rule in this file. Every feature starts here.**

Every feature increment starts from a failing **Playwright** (browser) scenario and is driven inward through unit-level red-green-refactor cycles.

#### Outer loop

1. **Red** — Write one Gherkin scenario describing the next observable user behavior. Run `npm test`. Confirm it fails for the expected reason. Do not proceed until the failure matches intent.
2. **Inner loop** — Repeat until you believe the scenario can pass:
   - **Red** — Write the smallest failing unit test for the next missing piece. One test at a time. Run it. Confirm it fails. No production code exists yet for this behavior.
   - **Green** — Write the minimum production code to pass that one test. Only what the test demands — nothing more. Run the test. **Commit.**
   - **Refactor** — Clean up covered code only. All tests stay green. **Commit.**
3. **Coverage gate** — Run `go test ./... -cover`. Any uncovered branch means production code was written without a test — go back to the inner loop.
4. **Green (scenario)** — Run the Playwright test. Still failing? Identify the next missing piece and return to the inner loop. Passes? **Commit and push.**
5. **Refactor (scenario)** — Refactor across modules if needed. All tests stay green. **Commit and push.**
6. Repeat from step 1.

#### Absolute rules

- **No production code without a red test.** A failing test must exist and have been run before any production file is created or edited. Every `if`, every `return`, every error branch — each one driven by a failing test first. No exceptions for "wiring", "infrastructure", or "obvious" code.
- **This applies to Vue code too.** Before creating any `.vue` or `.ts` file in `src/`, write a failing Vitest test. Run `npm run test:unit`. See it fail. Only then write the component. This rule was violated during the Vue migration — do not repeat it.
- **One branch, one test, one commit.** Never write a whole function and test it afterward.
- **Always run the full suite.** `go test ./...` and `npm run test:unit` at every red/green/refactor step.
- **Never skip red.** If you cannot articulate exactly why the test fails, stop.
- **Exception — external HTTP calls:** Handlers that make live HTTP calls can't be fully unit-tested without injecting the HTTP client. Note the gap as tech debt; do not silently accept low coverage.

#### Test commands

```bash
go test ./...          # Go unit tests — always run the full suite
npm run test:unit      # Vitest unit tests for Vue components
npm test               # e2e BDD outer loop (bddgen + playwright)
npx bddgen             # regenerate specs from .feature files only
```

Feature files live in `e2e/features/`. Step definitions live in `e2e/steps/`.
Always run `bddgen` after editing a `.feature` file.

#### Vue test patterns

- `vitest.config.ts` scopes to `src/**/*.test.ts` — do not change this or Playwright BDD specs get picked up
- Always create a **fresh router per test** via a factory function — sharing a router instance causes watcher accumulation and phantom fetch calls
- Mock `fetch` with `vi.stubGlobal('fetch', vi.fn()...)` and restore with `vi.restoreAllMocks()` in `beforeEach`
- Test files live in `__tests__/` subdirectories alongside the code they test

---

## Environment & Tooling

**Local dev (two terminals):**
- Terminal 1: `./dev.sh` — emulators (ports 8081/9099) + Go server (:8080)
- Terminal 2: `npm run dev` — Vite dev server (:5173) with hot reload; proxies `/api`, `/auth`, `/admin` to :8080
- Open `http://localhost:5173`

**E2e tests:** `npm test` builds Vue then starts Go server. Emulators must already be running.
- Kill a stale server: `kill $(lsof -ti :8080) 2>/dev/null`

**Firebase Admin SDK:** Pass `option.WithoutAuthentication()` when `FIREBASE_AUTH_EMULATOR_HOST` is set — otherwise `firebase.NewApp` fails silently in subprocess (no ADC credentials) and falls back to `NoopTokenVerifier`.
- `FIREBASE_PROJECT_ID` — Firebase Admin SDK init only, does **not** trigger Firestore
- `GOOGLE_CLOUD_PROJECT` — triggers Firestore; do **not** set in playwright `webServer` env

**E2e test isolation:** `e2e/global-setup.ts` calls `DELETE /admin/reset` (TEST_MODE=1) and clears Firebase Auth emulator accounts between runs.

---

## Project Goal

Predictions ranking system for Columbus Crew fans. Sarcastic #Crew96 tone. Only Crew — everyone else can pound sand.

**Scoring formats** are documented in [README.md](README.md) and implemented in `internal/scoring/`.

**Copy tone:** "Pick your scores. Be wrong in public. It's tradition."

---

## Design Language

**Theme:** Industrial Black & Gold Brutalism — matchday program crossed with a construction-site bulletin board.

| Token | Value | Use |
|---|---|---|
| `--black` | `#0c0c0c` | Page background |
| `--dark` | `#141414` | Header, input backgrounds |
| `--card` | `#1a1a1a` | Match card backgrounds |
| `--border` | `#2a2a2a` | Card/input borders |
| `--gold` | `#ffc20e` | Primary accent — Crew name, CTAs, active states |
| `--muted` | `#555` | Secondary text, metadata |
| `--danger` | `#e03c3c` | Locked/error states |

**Typography:**
- `Bebas Neue` — headings, team names, buttons (condensed, stadium-board energy)
- `DM Mono` — scores, timestamps, metadata (scoreboard digits)
- `Barlow` — body copy

**Key patterns:**
- 3px gold left border on hovered/predicted match cards
- Score inputs: 52×52px, `DM Mono`, gold text, dark background, focus glows gold
- Locked state: blinking `▊` indicator in `--danger` red
- Noise texture overlay via inline SVG `feTurbulence` filter on `body::before`
- Gold stripe on `body::after` (top of viewport, `position: fixed`)
- CSS lives in `src/style.css`, imported in `src/main.ts`

---

## Explaining Things

The user is learning Go, Vue 3, TypeScript, and the GCP/Firebase ecosystem through this project. When introducing a new pattern or tool:
- Explain *why* it works that way, not just *what* to write
- Show the mental model before the code
- Call out Go idioms that differ from other languages
