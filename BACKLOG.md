# Backlog

## Up Next

### Infrastructure

- [ ] **Decouple frontend from Docker image** — Go server currently embeds `dist/` and serves the SPA directly from Cloud Run as a fallback. Since Firebase Hosting is the real frontend entry point (and rewrites API paths to Cloud Run), the frontend doesn't need to be in the image. Refactor: remove `spaHandler`/`assetsHandler` from `main.go`, strip the Node/Vite build stage from the Dockerfile, build the Vue app as a separate CI step, and upload the artifact directly to GCS — no `docker cp` extraction needed. Smaller image, cleaner separation of concerns.

### Security

- [ ] **App Check** — register Crew Predictions web app with reCAPTCHA v3 attestation provider; enforce on Cloud Firestore and Authentication. Note: does not cover Go/Cloud Run endpoints (those are protected by session cookies). Wire `initializeAppCheck()` into `src/firebase.ts` before enforcing.

### Test Infrastructure

- [ ] **Per-worker server isolation** — current parallelism runs two Playwright projects (`auth` + `app`) against a shared server. If the app group grows too slow, give each worker its own Go server instance on a separate port so they don't share in-memory state.

### Low Priority

- [ ] **Prod smoke suite** — unauthenticated-only scenarios (app loads, leaderboard/matches API responds, Vue hydrates); replaces current `curl` liveness check in `deploy-prod`.
- [ ] **Real-data scoring accuracy test** — e2e scenario using actual 2025 Columbus Crew match results to validate the scoring engine against real outcomes. Get match data from user before writing.
- [ ] **Remove stale `handle` from predictions** — `handle` on prediction documents is legacy; leaderboard and profile now source display names from `UserStore` by `userID`. Stop writing it on new predictions and drop the field once confirmed no UID-less predictions exist in prod.

---

## Future / Exploratory

- [ ] **Migrate from Firestore to Firebase SQL Connect (PostgreSQL)** — relational model would be a better fit; SQL Connect wasn't available/mature when the project started. Consider when user base grows or query complexity increases.

---

## Decisions Made / Won't Do

- **Custom domain migration** — Firebase Hosting custom domain + Cloud Run domain mapping. Low priority — may never be needed.
- **Cloud Scheduler for match refresh** — `POST /admin/refresh-matches` is called manually after deploy. No cron job needed.
- **Match result entry UI** — admin page not needed; `POST /admin/results` API is sufficient for now.
- **Bluesky AT Protocol auth** — dropped for v1. Complex, adds no value for first release.
- **FirebaseUI** — dropped. Incompatible with `firebase@11` (requires `^9||^10`).
- **templ/HTMX/Alpine.js** — replaced by Vue 3 SPA; `templates/` package deleted.
- **GCP Cloud Build** — GitHub Actions preferred; simpler config, already where the code lives, Workload Identity Federation handles GCP auth without stored keys.
- **Frontend subdirectory** — kept as single project root (`package.json` at root, `src/` for Vue).

---

## Done

- [x] **Live match experience** — Now Playing section shows in-progress (`state=in`) and delayed matches above Upcoming; pulsing gold LIVE badge with match clock (`48'`, `HT`), blinking red DELAYED badge; no prediction inputs; STATUS_DELAYED rejects predictions with 403. Now Playing card is a clickable link to match detail. Match detail page shows LIVE indicator bar with clock, projected scores computed on-the-fly from live ESPN data (`isProjected: true` flag, projected label above table), smart polling every ~30s (self-rescheduling setTimeout, active window = live/delayed or within 30min of kickoff up to 2h after kickoff). ESPN live scores parsed correctly — scores arrive as plain strings (`"2"`) during play, not objects; `scoreField.UnmarshalJSON` handles both forms. `displayClock` field propagated through model → Firestore → API → Vue.
- [x] **Scoring engines** — AcesRadio (+15 exact, +10 correct result, −15 flipped scoreline, 0 otherwise); Upper 90 Club (+1 correct result, +1 correct Crew goals, +1 correct opponent goals, max +3); Grouchy™ (+1 for correct Columbus margin bucket: Win 2+, Win 1, Draw, Lose 1, Lose 2+). Rules page matches the actual engines. `internal/scoring` package; 100% coverage.
- [x] **Leaderboard** — unified sortable grid table (RANK · PREDICTOR · ACES RADIO · UPPER 90 CLUB · GROUCHY™); click headers to sort; dynamic tied ranks; shows users with ≥1 prediction at 0 pts before results land; profile link disabled for legacy handle-only users (`hasProfile` field). Mobile: stacked cards, sort buttons, active format score shown in gold. Precomputed via `Recalculate()` — O(U) reads instead of O(P×R) per request.
- [x] **Match detail page** — `/matches/:matchId`; unified sortable predictions table (RANK · PREDICTOR+PICK · ACES RADIO · UPPER 90 CLUB · GROUCHY™); prediction shown below handle; result cards link here, upcoming cards do not. `GET /api/matches/:matchId` returns per-format scores, `scoringFormats` array, live state, `isProjected`.
- [x] **Match data & polling infrastructure** — in-memory `MatchStore` backed by `WriteThroughMatchStore` (Firestore durable writes + fast memory reads); survives restarts. ESPN client fetches four league endpoints (MLS, US Open Cup, Leagues Cup, CONCACAF Champions); `fetchCrewMatchesFrom(base)` injectable for tests. `MatchPoller` schedules per-match kickoff timers, ticks ESPN every 2 minutes while active, writes `ResultStore` on terminal status. Daily refresh at 4am ET resets pollers and runs `Backfill()` to catch results that finalized during downtime. `POST /admin/refresh-matches` and `POST /admin/poll-scores` for manual triggers and e2e.
- [x] **Score recalculation** — `internal/recalculator.Recalculate()` recomputes all three format totals and prediction count per user from scratch; upserts to `UserStore`. Triggered after every match final (via `MatchPoller.SetOnResultSaved`) and on startup after `Backfill()`. Leaderboard and profile read precomputed values.
- [x] **CI/CD pipeline** — GitHub Actions; Go unit tests + Vue unit tests + TypeScript typecheck + e2e BDD suite on every develop push; deploy-staging (Docker → Artifact Registry → Cloud Run staging → Firebase Hosting staging → smoke tests); deploy-prod on main merge (promotes staging artifact — no rebuild; Cloud Run prod → Firebase Hosting prod → liveness check → automatic rollback on failure). GCS frontend artifacts (`sha-{SHA}.zip` / `latest.zip` / `prod.zip` / `prod-previous.zip`). Concurrency cancel on develop. Workload Identity Federation for GCP auth (no stored keys). Firebase config env vars (`FIREBASE_PROJECT_ID`, `FIREBASE_API_KEY`, `FIREBASE_AUTH_DOMAIN`, `GOOGLE_CLOUD_PROJECT`) set correctly for both staging and prod via `--set-env-vars`.
- [x] **Staging environment** — separate GCP project (`crew-predictions-staging`); Cloud Run + Firebase Hosting (`crew-predictions-staging.web.app`); staging smoke suite hits real staging URL with permanent test accounts (no account creation); `authDomain` correctly set to `.firebaseapp.com`.
- [x] **Artifact Registry cleanup** — `infra/artifact-policy.json` keeps `prod`-tagged and `latest`-tagged images indefinitely; deletes everything else after 4 hours. `prod` tag applied only after staging smoke passes. Separate policy for Cloud Functions repo: keep 3 most recent, delete after 4 hours.
- [x] **Firebase Analytics** — `initAnalytics()` guarded by `measurementId`, `appId`, and `projectId` (all three required — Firebase Installations needs `projectId` or it throws); called inside try/catch in `bootstrap()` so a crash never blocks `mount('#app')`.
- [x] **Auth** — email/password sign-in + sign-up; Google SSO via `signInWithRedirect` (redirect not popup, for mobile); `getRedirectResult()` in App.vue on mount; password reset flow. Session cookies: HttpOnly, HMAC-SHA256 signed (`SESSION_SECRET` from Secret Manager), `Secure` flag on in prod. `postSession()` checks `res.ok` before navigating.
- [x] **User & profile** — `users/{userID}` Firestore collection as source of truth for display names; `UserStore` interface split: `Upsert` (profile fields) and `UpdateScores` (scoring fields) are separate so auth handlers can't accidentally zero scores. `GET /api/me` lazily upserts on app open. `POST /auth/handle` updates name + location. `/profile/:userID` page shows 4-stat grid (Predictions · Aces Radio · Upper 90 · Grouchy™ each with rank); edit form on own profile only. `POST /admin/backfill-users` seeds `users` from existing predictions.
- [x] **TwoOneBot** — `internal/bot` package; `bot:twooonebot` user predicts Columbus 2-1 (home) / 1-2 (away) on every upcoming match at each refresh and daily tick. Appears on leaderboard and match detail like any other user.
- [x] **Guest predictions** — logged-out users enter picks stored in localStorage; `flushGuestPredictions()` auto-submits on login/signup so predictions aren't lost. Safari Private Browsing handled gracefully.
- [x] **UI/UX** — Industrial Black & Gold Brutalism design (noise texture, gold stripe, DM Mono scores); official Crew gold `#fedd00`; Barlow Condensed 800 display font; mobile responsive layout (cards stack, full-width Predict button); hamburger nav drawer at ≤600px; match ordering stable (sort by kickoff ascending, results descending); loading and error states on all views; autocomplete attributes on auth forms; team name stacks vertically at ≤360px.
- [x] **Security hardening** — admin endpoints guarded by `AdminAuth` middleware (X-Admin-Key, `crypto/subtle` compare); server refuses to start if `ADMIN_KEY` unset in prod. Rate limiting on leaderboard and profile endpoints. Container runs as non-root (uid 1001). `.dockerignore` excludes credentials and test artifacts. Goals range validated (0–99). `match_id` validated for emptiness. `ParseForm` errors handled. `Save()` error returned to caller. Patch CVE-2026-34986 (`go-jose/v4` → 4.1.4).
- [x] **Billing killswitch** — Cloud Function (`infra/billing-killswitch/`, nodejs24, gen2, `us-east5`) disables billing on both GCP projects when a budget alert fires (`costAmount > budgetAmount` on `billing-alerts` Pub/Sub topic).
- [x] **BDD e2e suite** — Playwright; two parallel projects (`auth` + `app`); `@reset` tag gates the reset hook to app features; covers auth, predictions, leaderboard, match detail, mobile layout, guest predictions, live match state.
