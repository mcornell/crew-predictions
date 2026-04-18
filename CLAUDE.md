# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Stack

| Layer | Technology |
|---|---|
| Language | Go |
| Compute | GCP Cloud Run (serverless, no container management) |
| Frontend | templ + HTMX + Alpine.js |
| Database | Firestore (GCP always-free) |
| Auth (MVP) | Firebase Auth — Google OAuth |
| Auth (follow-on) | Bluesky AT Protocol OAuth |
| Static assets | Firebase Hosting |
| Match data | ESPN unofficial API → Firestore cache (daily refresh) |

## Frontend Pattern

**HTMX** handles server round-trips: submitting a prediction, loading a partial leaderboard update. The server returns an HTML fragment and HTMX swaps it into the page — no JSON, no client-side fetch.

**Alpine.js** handles purely client-side state that needs no server involvement: sorting a leaderboard table, toggling a UI element. No round-trip needed.

**When to use which:** If the action needs data from the server → HTMX. If it's purely presentational and the data is already on the page → Alpine.js.

**templ** is a type-safe Go templating language that compiles to Go code. Run `templ generate` before `go build`. templ files live in `templates/` and have a `.templ` extension.

## Development Approach

### BDD Dual-Loop TDD

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

- **No production code without a red test.** This is absolute. A failing test must exist and have been run before any production file is created or edited. This means every `if`, every `return`, every error branch — each one must be driven by its own failing test first.
- **One branch, one test, one commit.** Every error path, happy path, and edge case is its own red-green-refactor cycle. Never write a whole function and test it afterward. If you do this correctly, coverage is never something to check at the end — it's guaranteed by construction.
- **Always run `go test ./...` (full suite).** Never run a subset of packages. Run the full suite at every red, green, and refactor step.
- **Never skip red.** If you cannot articulate exactly why the test fails, stop.
- **Exception — external HTTP calls:** Handlers that make live HTTP calls can't be fully unit-tested without injecting the HTTP client. Note the gap explicitly as tech debt; do not silently accept low coverage.

#### Test commands

```bash
go test ./...          # unit tests (Go testing package)
npm test               # e2e BDD outer loop (bddgen + playwright)
npx bddgen             # regenerate specs from .feature files only
templ generate         # must run before go build if .templ files changed
```

Feature files live in `e2e/features/`. Step definitions live in `e2e/steps/`.
Always run `bddgen` after editing a `.feature` file — it generates the `.features-gen/` specs that Playwright actually executes.

## Explaining Things to the User

The user is learning Go, HTMX, Alpine.js, and the GCP/Firebase ecosystem as we build. When introducing a new pattern or tool:
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
| Correct winner (wrong score) | +10 |
| Predicted wrong team wins by the exact same scoreline (e.g. predict Crew 3–2 Portland, actual Portland 3–2 Crew) | −15 |
| Anything else | 0 |

**Upper 90 Club** (confirmed):
| Outcome | Points |
|---|---|
| Exact score | +2 |
| Correct winner (wrong score) | +1 |
| Anything else | 0 |

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
go run ./cmd/server                                          # local dev
firebase emulators:start --only firestore                    # local Firestore
gcloud run deploy crew-predictions --source . --region us-east5
firebase deploy --only hosting
```
