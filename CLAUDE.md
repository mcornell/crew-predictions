# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Approach

### BDD Dual-Loop TDD

**This is the highest-priority rule in this file. Every feature starts here.**

Every feature increment starts from a failing **Playwright** (browser) scenario and is driven inward through unit-level red-green-refactor cycles.

#### Outer loop

1. **Red** — Write one Gherkin scenario describing the next observable user behavior. Run `npm test`. Confirm it fails for the expected reason. Do not proceed until the failure matches intent.
2. **Inner loop** — Repeat until you believe the scenario can pass:
   - **Red** — Write the smallest failing Go unit test for the next missing piece. One test at a time. Run it. Confirm it fails. No production code exists yet for this behavior.
   - **Green** — Write the minimum production code to pass that one test. Only what the test demands — nothing more. Run the test. **Commit.**
   - **Refactor** — Clean up covered code only. All tests stay green. **Commit.**
3. **Coverage gate** — Run `go test ./... -cover`. Any uncovered branch means production code was written without a test — go back to the inner loop. Do not proceed until coverage is clean.
4. **Green (scenario)** — Run the Playwright test. Still failing? Identify the next missing piece and return to the inner loop. Passes? **Commit and push.**
5. **Refactor (scenario)** — Refactor across modules if needed. All tests stay green. **Commit and push.**
6. Repeat from step 1.

#### Absolute rules

- **No production code without a red test.** This is absolute. A failing test must exist and have been run before any production file is created or edited. This means every `if`, every `return`, every error branch — each one must be driven by its own failing test first. There are no exceptions for "wiring", "infrastructure", or "obvious" code.
- **This applies to Vue code too.** Before creating any `.vue` or `.ts` file in `src/`, write a failing Vitest test. Run `npm run test:unit`. See it fail. Only then write the component. This rule was violated during the Vue migration and required significant rework. Do not repeat that mistake.
- **One branch, one test, one commit.** Every error path, happy path, and edge case is its own red-green-refactor cycle. Never write a whole function and test it afterward.
- **Always run the full test suite.** `go test ./...` and `npm run test:unit` at every red, green, and refactor step. Never run a subset.
- **Never skip red.** If you cannot articulate exactly why the test fails, stop. Do not proceed.
- **Exception — external HTTP calls:** Handlers that make live HTTP calls can't be fully unit-tested without injecting the HTTP client. Note the gap explicitly as tech debt; do not silently accept low coverage.

#### Test commands

```bash
go test ./...          # Go unit tests — always run the full suite
npm run test:unit      # Vitest unit tests for Vue components
npm test               # e2e BDD outer loop (bddgen + playwright)
npx bddgen             # regenerate specs from .feature files only
```

Feature files live in `e2e/features/`. Step definitions live in `e2e/steps/`.
Always run `bddgen` after editing a `.feature` file — it generates the `.features-gen/` specs that Playwright actually executes.

#### Vue inner loop (same discipline, different tool)

The inner loop applies to Vue code exactly as it does to Go:
- Write a **failing Vitest test** (`npm run test:unit`) before creating or editing any `.vue` or `.ts` file in `src/`
- Confirm the test fails for the right reason before writing the component
- Write minimum code to pass; run `npm run test:unit` again
- Then check the e2e outer loop

**Vue test patterns:**
- `vitest.config.ts` scopes to `src/**/*.test.ts` — do not change this or Playwright BDD specs will be picked up
- Always create a **fresh router per test** via a factory function — sharing a router instance across tests causes watcher accumulation and phantom fetch calls
- Mock `fetch` with `vi.stubGlobal('fetch', vi.fn()...)` and restore with `vi.restoreAllMocks()` in `beforeEach`
- Test files live alongside components in `__tests__/` subdirectories

### Environment & Tooling

**Local dev (two terminals):**
- Terminal 1: `./dev.sh` — starts Firestore + Auth emulators (ports 8081/9099) then Go server on :8080
- Terminal 2: `npm run dev` — Vite dev server on :5173 with hot reload; proxies `/api`, `/auth`, `/admin` to :8080
- Open `http://localhost:5173` in the browser

**E2e tests:** `npm test` — builds Vue (`vite build`) then starts Go server via `playwright.config.ts` webServer. Emulators must already be running (`firebase emulators:start --only firestore,auth`).

- Playwright `globalSetup` runs **after** `webServer` is ready — server endpoints can be called from it.
- Kill a stale server before debugging: `kill $(lsof -ti :8080) 2>/dev/null`

### Firebase Admin SDK + Emulator

When the server runs without ADC credentials (e.g. in the Playwright `webServer` subprocess), `firebase.NewApp(ctx, nil)` fails silently and falls back to `NoopTokenVerifier`. Fix is already in `cmd/server/main.go`: pass `option.WithoutAuthentication()` when `FIREBASE_AUTH_EMULATOR_HOST` is set.

- `FIREBASE_PROJECT_ID` — used for Firebase Admin SDK init; does **not** trigger Firestore
- `GOOGLE_CLOUD_PROJECT` — triggers Firestore; do **not** set this in playwright `webServer` env

### E2E Test Isolation

In-memory store state persists across test runs when the server is reused (`reuseExistingServer: true`). `e2e/global-setup.ts` resets both by calling `DELETE /admin/reset` (only registered when `TEST_MODE=1`) and clearing the Firebase Auth emulator accounts.

### Known Tech Debt

- **ESPN date parsing** — ESPN returns times without seconds (e.g. `2026-04-12T23:00Z`), which fails Go's `time.RFC3339` (requires seconds). All kickoff times currently parse to Go zero value (`0001-01-01`). Fix in `internal/espn/client.go`: try `"2006-01-02T15:04Z07:00"` as fallback.
- **firebaseui incompatible** — `firebaseui@6` requires `firebase@^9||^10`; we use `firebase@11`. Do not attempt to install it. Auth uses a custom email/password form + Firebase Auth SDK (`signInWithEmailAndPassword`) directly.
- **Frontend unstyled** — Vue SPA is functional but has no CSS. Design language (Industrial Black & Gold Brutalism) defined in the Design Language section below has not been applied yet.

---

# Stack

| Layer | Technology |
|---|---|
| Language | Go |
| Compute | GCP Cloud Run (serverless, no container management) |
| Frontend | Vue 3 + TypeScript + Vite |
| Database | Firestore (GCP always-free) |
| Auth | Firebase Auth — Email/Password (custom form; no FirebaseUI) |
| Static assets | Firebase Hosting |
| Match data | ESPN unofficial API → Firestore cache (daily refresh) |

## Frontend Pattern

**Vue 3 SPA** (`src/`) served from Go's `dist/` static file handler. All data flows through JSON APIs (`/api/*`). Go server handles auth session cookies.

**API conventions:**
- `GET /api/matches` → `{ matches: [...], predictions: {...} }`
- `POST /api/predictions` → form data (`URLSearchParams`), not JSON — Go uses `r.ParseForm()`
- `POST /auth/session` → form data (`URLSearchParams`) with `idToken` field
- `GET /api/me` → `{ handle }` or 401
- `GET /api/leaderboard` → `{ acesRadio: [...], upper90Club: [...] }`

**CSS** lives in `src/style.css` (imported in `src/main.ts`). Design tokens are CSS variables defined on `:root`.

## Explaining Things to the User

The user is learning Go, Vue 3, TypeScript, and the GCP/Firebase ecosystem as we build. When introducing a new pattern or tool:
- Explain *why* it works that way, not just *what* to write
- Show the mental model before the code
- Call out Go idioms that differ from other languages
- Point to the relevant docs section

## Project Goal

Create a predictions ranking system for fans of Columbus Crew. Sarcastic tone, like #Crew96 fandom. Only Crew — everyone else can pound sand.

Scoring rules come in two flavors matching podcast formats:

**Aces Radio** (confirmed):
| Outcome | Points |
|---|---|
| Exact score | +15 |
| Correct result (wrong score) | +10 |
| Scores exactly mirrored — predicted wrong team to win by the same scoreline (e.g. predict Crew 3–2 Portland, actual Portland 3–2 Crew) | −15 |
| Anything else | 0 |

Note: checks apply in order — exact score is evaluated before the mirror check.

**Upper 90 Club** (confirmed) — two independent points that stack:
| Condition | Points |
|---|---|
| Correct match result (win/draw/loss) | +1 |
| Correct Columbus Crew goal count | +1 |

Max 2 points per match. You do **not** need an exact score for +2 — e.g. predict 1–0 Crew win with Crew scoring 0, actual 2–0 Crew win: +1 (correct result) + 0 (wrong Crew goals) = 1 pt. See `internal/scoring/upper90club.go` and its tests for the full implementation.

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
- CSS lives in `static/style.css`; served via Go's `http.FileServer`

**Tone in copy:** Sarcastic #Crew96 fandom. "Pick your scores. Be wrong in public. It's tradition."

## Sources

Always provide sources when responding.

## Deployment

```bash
# Local dev
firebase emulators:start --only firestore,auth   # terminal 1
./dev.sh                                          # terminal 1 (or separately)
npm run dev                                       # terminal 2, open :5173

# Production build
npm run build                                     # outputs to dist/
gcloud run deploy crew-predictions --source . --region us-east5
firebase deploy --only hosting
```
