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

#### Outer loop (Playwright/BDD scenario)

1. **Red** — Write one Gherkin scenario in a `.feature` file under `e2e/features/` describing the next observable user behavior. Run `npm test`. Confirm it fails for the expected reason (missing step definition = red, wrong behavior = red). Do not proceed until the failure matches intent.
2. **Inner loop** — Repeat until the Playwright test can pass:
   - **Red** — Write the smallest Go unit test for the next missing piece the scenario needs. One test at a time. Run it. Confirm it fails.
   - **Green** — Write the **minimum** production code to make that unit test pass. No speculative code. No implementing more than the test demands.
   - **Refactor** — Clean up only covered code. All unit tests must stay green.
3. **Green (scenario)** — Re-run the Playwright test. If still failing, identify the missing piece and return to the inner loop.
4. **Refactor (scenario)** — Refactor across modules if needed. All tests must stay green.
5. Repeat from step 1.

#### Test commands

```bash
go test ./...          # unit tests (Go testing package)
npm test               # e2e BDD outer loop (bddgen + playwright)
npx bddgen             # regenerate specs from .feature files only
templ generate         # must run before go build if .templ files changed
```

Feature files live in `e2e/features/`. Step definitions live in `e2e/steps/`.
Always run `bddgen` after editing a `.feature` file — it generates the `.features-gen/` specs that Playwright actually executes.

#### Discipline rules

- **Never skip red.** If you cannot articulate why a test fails, stop and re-read the requirement.
- **One test at a time.** Never write multiple tests before running them.
- **Minimum code.** Only write production code demanded by the current failing test. Stub everything else.
- **Ask before assuming.** If a design decision is unclear, ask the user before writing code.
- **Commit on every green step** — unit green, scenario green, AND after every refactor. Three distinct commits per cycle, not one.
- **Never skip the refactor commit.** Refactor → run full suite → commit. Do not proceed to the next red until this is done.
- **Run only the relevant test** after each green step; run the full suite before committing.

## Explaining Things to the User

The user is learning Go, HTMX, Alpine.js, and the GCP/Firebase ecosystem as we build. When introducing a new pattern or tool:
- Explain *why* it works that way, not just *what* to write
- Show the mental model before the code
- Call out Go idioms that differ from other languages
- Point to the relevant docs section

## Project Goal

Create a predictions ranking system for fans of Columbus Crew. Sarcastic tone, like #Crew96 fandom. Only Crew — everyone else can pound sand.

Scoring rules come in two flavors matching podcast formats: Aces Radio and Upper 90. Confirm full point tables before implementing `internal/scoring/`.

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
