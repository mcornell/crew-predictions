# Backlog

## Up Next

1. [ ] **Profile page needs context** — currently just a display name form floating in space; add current prediction count, scoring summary, or other stats to make it worth visiting.

2. [ ] **Custom domain migration** — Firebase Hosting custom domain + Cloud Run domain mapping. Update `authDomain` and OAuth redirect URIs.

3. [ ] **Leaderboard: show users with ≥1 prediction at 0 pts before results land** — seed leaderboard from unique userIDs in predictions rather than only from scored pairs; smoke test accounts won't appear since they never predict.

4. [ ] **Cache leaderboard scoring** — currently recalculated on every request; fine now but will need in-memory caching or pre-computation at scale.

5. [ ] **Remove stale `handle` from predictions** — `handle` on prediction documents is legacy; leaderboard now sources display names from `UserStore` by `userID`. Stop writing it on new predictions and remove the `p.Handle` fallback in the leaderboard once confirmed no UID-less predictions exist in prod.

---

## Data & Polling

- [ ] **Real-data scoring accuracy test** — e2e scenario using actual 2025 Columbus Crew match results to validate the scoring engine against real outcomes. Get match data from user before writing.
- [x] **Score polling** — `MatchPoller` schedules a per-match timer at kickoff; ticks ESPN every 2 minutes while active; writes ResultStore on terminal status then deactivates. Unknown/postponed matches run until 4am reset. `POST /admin/poll-scores` for manual triggers and e2e.
- [x] **Live match state** — `state` field (`"pre"/"in"/"post"`) parsed from ESPN `status.state`; propagated through model → API → Vue; pulsing LIVE badge on in-progress match cards.
- [x] **Daily match refresh at 4am ET** — replaces 24h-from-startup ticker; calls `poller.Reset` so match pollers are rescheduled from fresh data after each refresh. Manual `POST /admin/refresh-matches` has the same effect.

---

## Test Infrastructure

- [ ] **Playwright smoke suite for prod** — identify a small tagged subset of e2e scenarios that can run against the live prod URL after deploy (replaces the current `curl` liveness check in `deploy-prod`)
- [ ] **Per-worker server isolation** — current parallelism runs two Playwright projects (`auth` + `app`) against a shared server. If the app group grows too slow, give each worker its own Go server instance on a separate port so they don't share in-memory state.

---

## Decisions Made / Won't Do

- **Cloud Scheduler for match refresh** — `POST /admin/refresh-matches` is called manually after deploy. No cron job needed.


- **Match result entry UI** — admin page not needed; `POST /admin/results` API is sufficient for now.
- **Bluesky AT Protocol auth** — dropped for v1. Complex, adds no value for first release.
- **FirebaseUI** — dropped. Incompatible with `firebase@11` (requires `^9||^10`).
- **templ/HTMX/Alpine.js** — replaced by Vue 3 SPA; `templates/` package deleted.
- **GCP Cloud Build** — GitHub Actions preferred; simpler config, already where the code lives, Workload Identity Federation handles GCP auth without stored keys.
- **Frontend subdirectory** — kept as single project root (`package.json` at root, `src/` for Vue).

---

## Done

- [x] **Handle management + UserStore** — `users/{userID}` Firestore collection as source of truth for display names; leaderboard groups by `userID` and joins `UserStore` for current handle; `POST /auth/handle` upserts on profile save; `GET /api/me` lazily upserts returning users on app open (catches users who were logged in before the feature shipped); `POST /admin/backfill-users` seeds `users` from existing predictions.
- [x] **Match persistence across restarts** — `FirestoreMatchStore` persists full season match data to `matches/{matchID}`; `WriteThroughMatchStore` wraps memory (fast reads) + Firestore (durable writes); on startup, loads stored matches from Firestore into memory before ESPN fires; past-kickoff matches immediately scheduled for catch-up polling.
- [x] **Unlock picks + countdown** — Unlock button clears a saved prediction and pre-populates inputs with the previous pick; client-side only (server 403 after kickoff is the real gate). Live "locks in Xd Yh / Xh Ym / Xm" countdown on each match card using browser clock (cosmetic). Upcoming window extended to 8 days.
- [x] **Fix Upper 90 Club scoring rules** — +1 correct outcome, +1 correct Crew goals, +1 correct opponent goals (max 3 pts). Real-data tests updated.
- [x] **Match cache + ESPN fetch** — in-memory `MatchStore` populated via `POST /admin/refresh-matches`; ESPN fetcher injected (TEST_MODE reads seeded store); `fetchCrewMatchesFrom(base)` tested via `httptest.Server` + captured fixture JSON; 97% ESPN package coverage
- [x] **E2e parallelisation** — two Playwright projects (`auth` + `app`) run in parallel; `@reset` tag gates the `Before` reset hook to app features only
- [x] **Coverage drive** — handlers 95%, espn 97%, scoring 100%; extracted `toPredictions`/`toResult`/`isNotFound`/`tokenToFirebaseToken` to make SDK mapping logic unit-testable
- [x] **Guest predictions** — logged-out users can enter picks stored in localStorage; never hits server or leaderboard; nudge to create account after predicting
- [x] **Vue Router warnings in unit tests** — shared `makeRouter()` utility with catch-all route silences all "No match found" warnings across test files
- [x] **Team name truncation on 360px (Galaxy S24)** — stacks team names vertically at ≤600px; e2e covered
- [x] **Typography overhaul** — replaced Bebas Neue with Barlow Condensed 800 (closer to official MLS/Crew aesthetic, works in mixed case); bumped font sizes across the board; fixed button vertical centering (flexbox, replacing asymmetric padding hack for Bebas Neue baseline quirk)
- [x] **Official Crew brand colors** — `--gold: #fedd00` (sourced from columbuscrew.com computed styles); nav link and muted text contrast improved (`#888`)
- [x] **Autocomplete attributes** — `autocomplete="email"` on all email inputs, `current-password` on login, `new-password` on signup; fixes password manager autofill and removes browser console warning
- [x] **Google sign-in redirect** — switched from `signInWithPopup` to `signInWithRedirect` everywhere; `getRedirectResult()` called in App.vue onMounted with try/catch so fetchUser always runs even if redirect result fails
- [x] **Mobile responsive layout** — cards stack at 390px and 412px; header stays single-row; Predict button full-width at 48px
- [x] **UX fixes — broken flows** — logged-out predict redirects to `/login`; 404 NotFoundView; score inputs gated behind auth; profile display name pre-populated; Google sign-in popup fixed
- [x] **Auth UX polish** — login/signup cross-links; error-state differentiation; logout UI; password-reset flow; display name / profile page
- [x] **Sign-up flow** — `/signup` view, `createUserWithEmailAndPassword`, reuses `/auth/session` cookie flow
- [x] **Google SSO** — Google provider + OAuth web client secret configured in Firebase Console
- [x] **Staging Cloud Run + artifact promotion** — CI deploys develop → staging (crew-predictions-staging GCP project); main promoted from staging artifact
- [x] **Staging Firebase Hosting** — frontend deployed to `crew-predictions-staging.web.app` so app origin matches authDomain; fixes Google sign-in redirect on staging
- [x] **Mobile hamburger nav** — full-width drawer with fast slide animation; closes on link tap or outside tap; hamburger shown at ≤480px, desktop nav hidden; user handle visible in header on mobile; e2e covered
- [x] **Results reverse chronological order** — completed matches sorted most recent first on matches page
- [x] **Remove email verification banner** — banner removed; users go straight to predictions after sign-up
- [x] **Staging smoke suite** — post-deploy Playwright suite hits real staging URL; covers email sign-in/sign-up on desktop + mobile viewports; Google redirect initiation; two permanent test accounts with setup/teardown
- [x] **GitHub Actions CI/CD** — Go + Vitest + Playwright on push; Workload Identity Federation for GCP auth
- [x] **BDD e2e suite** — Playwright scenarios covering auth, predictions, leaderboard, mobile layout
- [x] **Server-side prediction locking** — 403 after kickoff; ESPN fetcher injected into handler
- [x] **Multi-competition support** — MLS, US Open Cup, Leagues Cup, CONCACAF Champions
- [x] **Match listings** — upcoming (next 7 days) + results with scores inline
- [x] **Leaderboard** — Aces Radio and Upper 90 Club formats; JSON API + Vue view
- [x] **Scoring engines** — AcesRadio (+15/+10/−15/0) and Upper 90 Club (+1 result, +1 Columbus goals, stacking)
- [x] **FirestoreResultStore** — match results persist across restarts; thread-safe stores with `sync.RWMutex`
- [x] **Firebase Auth** — email/password; session cookies (`HttpOnly`); Firebase Admin SDK with emulator support
- [x] **Vue 3 SPA** — MatchesView, LoginView, SignupView, LeaderboardView, ProfileView, AppHeader; Vite dev proxy
- [x] **Go server** — ESPN match fetching; Firestore prediction store; seed endpoints for e2e fixtures
- [x] **Industrial Black & Gold Brutalism design** — noise texture, gold stripe, Bebas Neue (now replaced), DM Mono scores, match card hover states
- [x] **Patch CVE-2026-34986** — upgraded `go-jose/v4` to 4.1.4
