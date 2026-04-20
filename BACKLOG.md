# Backlog

## Done

- [x] Go server with ESPN match fetching
- [x] Firestore prediction store
- [x] AcesRadio scoring engine (+15/+10/−15/0)
- [x] Upper 90 Club scoring engine (+1 result, +1 Columbus goals, stacking)
- [x] Leaderboard (both formats, JSON API + Vue view)
- [x] Firebase Auth — Email/Password + Google OAuth
- [x] Session cookies (`HttpOnly`, base64 JSON)
- [x] Vue 3 SPA: MatchesView, LoginView, LeaderboardView, AppHeader
- [x] BDD e2e suite — 11/11 Playwright scenarios green
- [x] Vite dev proxy for local development
- [x] Industrial Black & Gold Brutalism design applied
- [x] ESPN date parsing fix (`2026-04-12T23:00Z` no-seconds format)
- [x] Match listings — upcoming (next 7 days) + results with scores inline
- [x] Multi-competition support — MLS, US Open Cup, Leagues Cup, CONCACAF Champions
- [x] Server-side prediction locking — 403 after kickoff; ESPN fetcher injected into handler
- [x] FirestoreResultStore — match results persist across restarts
- [x] Seed endpoints (`/admin/seed-match`, `/admin/seed-prediction`) — deterministic e2e fixtures, no ESPN dependency in tests
- [x] Thread-safe in-memory stores — `sync.RWMutex` on prediction and result stores
- [x] E2e test isolation — `Before` hook resets all stores per scenario; serial workers prevent shared-state races
- [x] Remove templ — deleted `templates/` package and dead HTML-rendering `List` handlers after Vue SPA migration
- [x] Patch CVE-2026-34986 — upgraded `go-jose/v4` to 4.1.4 (transitive dep via Firebase → gRPC → SPIFFE)

---

## Before First Deploy

- [x] **FirestoreResultStore** — results persist across restarts
- [x] **Prediction locking** — enforce kickoff time deadline server-side; predictions must be submitted before kickoff
- [x] **GCP/Firebase env vars** — `FIREBASE_API_KEY`, `FIREBASE_AUTH_DOMAIN`, `FIREBASE_PROJECT_ID`, `GOOGLE_CLOUD_PROJECT` set on Cloud Run service; old OAuth vars removed
- [x] **Cloud Run deploy** — live at https://crew-predictions-937208344837.us-east5.run.app
- [x] **Firebase Hosting** — live at https://crew-predictions.web.app; rewrites `/api/**`, `/auth/**`, `/admin/**` to Cloud Run; SPA fallback for all other routes

---

## Post-Launch

- [ ] **GitHub Actions CI/CD** — push to `main` → test + deploy via Workload Identity Federation (keyless GCP auth, no stored service account keys).
- [ ] **Real-data scoring accuracy test** — e2e scenario using actual 2025 Columbus Crew match results to validate the scoring engine against real outcomes. Get match data from user before writing.
- [ ] **Firestore match cache + score polling** — cache ESPN schedule in Firestore. Weekly refresh fires Tuesday midnight ET (cron) regardless of whether there's a match that week. When `kickoff + 2h <= now` and match not yet `STATUS_FINAL`, poll ESPN every ~5 min and write to ResultStore when final, then stop. ESPN already returns `status.type.name`.

---

## Post-Launch (Test Infrastructure)

- [ ] **Per-worker server isolation** — if the e2e suite grows large enough that serial execution is too slow, give each Playwright worker its own server instance (separate ports) so parallel runs don't share in-memory state.

---

## Decisions Made / Won't Do

- **Match result entry UI** — admin page not needed; `POST /admin/results` API is sufficient for now.
- **Bluesky AT Protocol auth** — dropped for v1. Complex, adds no value for first release.
- **FirebaseUI** — dropped. Incompatible with `firebase@11` (requires `^9||^10`).
- **templ/HTMX/Alpine.js** — replaced by Vue 3 SPA; `templates/` package deleted.
- **GCP Cloud Build** — GitHub Actions preferred; simpler config, already where the code lives, Workload Identity Federation handles GCP auth without stored keys.
- **Frontend subdirectory** — kept as single project root (`package.json` at root, `src/` for Vue).
