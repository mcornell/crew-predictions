# Backlog

## Up Next

### Security

- [ ] **HMAC-sign session cookies** — the session cookie is currently base64-encoded JSON with no integrity check; any client-side modification is accepted as valid. Sign with HMAC-SHA256 using a secret loaded from Secret Manager (`SESSION_SECRET`). Verify signature in `UserFromSession`; reject unsigned/tampered cookies with 401.

- [ ] **Rate limit expensive endpoints** — `/api/leaderboard` and `/api/profile/:userID` hit Firestore on every request with no throttle. Add per-IP rate limiting (e.g. 60 req/min) using an in-memory token bucket. Cloud Run's single-instance concurrency makes in-process state viable; revisit if multi-instance needed. **Known risk:** staging smoke suite hits the leaderboard from a GitHub Actions runner IP — if it breaches 60 req/min the smoke tests will get 429s. If this happens, add `RATE_LIMIT_ENABLED=true` env var and only enable in prod.

### Low

- [ ] **Decouple frontend from Docker image** — Go server currently embeds `dist/` and serves the SPA directly from Cloud Run as a fallback. Since Firebase Hosting is the real frontend entry point (and rewrites API paths to Cloud Run), the frontend doesn't need to be in the image. Refactor: remove `spaHandler`/`assetsHandler` from `main.go`, strip the Node/Vite build stage from the Dockerfile, build the Vue app as a separate CI step, and upload the artifact directly to GCS — no `docker cp` extraction needed. Smaller image, cleaner separation of concerns.

14. [ ] **Prod smoke suite** — unauthenticated-only scenarios (app loads, leaderboard/matches API responds, Vue hydrates); replaces current `curl` liveness check in `deploy-prod`.

17. [ ] **Remove stale `handle` from predictions** — `handle` on prediction documents is legacy; leaderboard and profile now source display names from `UserStore` by `userID`. Stop writing it on new predictions and drop the field once confirmed no UID-less predictions exist in prod.

---

## Data & Polling

- [ ] **Real-data scoring accuracy test** — e2e scenario using actual 2025 Columbus Crew match results to validate the scoring engine against real outcomes. Get match data from user before writing.
- [x] **Score polling** — `MatchPoller` schedules a per-match timer at kickoff; ticks ESPN every 2 minutes while active; writes ResultStore on terminal status then deactivates. Unknown/postponed matches run until 4am reset. `POST /admin/poll-scores` for manual triggers and e2e.
- [x] **Live match state** — `state` field (`"pre"/"in"/"post"`) parsed from ESPN `status.state`; propagated through model → API → Vue; pulsing LIVE badge on in-progress match cards.
- [x] **Daily match refresh at 4am ET** — replaces 24h-from-startup ticker; calls `poller.Reset` so match pollers are rescheduled from fresh data after each refresh. Manual `POST /admin/refresh-matches` has the same effect.

---

## Test Infrastructure

- [ ] **Per-worker server isolation** — current parallelism runs two Playwright projects (`auth` + `app`) against a shared server. If the app group grows too slow, give each worker its own Go server instance on a separate port so they don't share in-memory state.

---

## Security (Low Priority)

- [ ] **App Check** — register Crew Predictions web app with reCAPTCHA v3 attestation provider; enforce on Cloud Firestore and Authentication. Note: does not cover Go/Cloud Run endpoints (those are protected by session cookies). Wire `initializeAppCheck()` into `src/firebase.ts` before enforcing.

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

- [x] **Artifact Registry cleanup policy** — `infra/artifact-policy.json` keeps `prod`-tagged and `latest`-tagged images; deletes everything else after 4 hours. `prod` tag applied in CI after staging smoke test passes. Separate `gcf-artifact-policy.json` for Cloud Functions repo: keep 3 most recent, delete after 4 hours.

- [x] **Firestore scoring fields read back correctly** — `toUser()` now deserialises all scoring fields (`acesRadioPoints`, `upper90Points`, `grouchyPoints`, `predictionCount`) from Firestore documents. Bug: fields were written correctly but never read back, so staging leaderboard showed zero scores for all users. Rule added to `internal/CLAUDE.md`: every new `FirestoreUserStore` method requires an integration test that round-trips through the emulator.

- [x] **UserStore interface split** — `Upsert` (profile fields: handle, provider, location) and `UpdateScores` (scoring fields: points + count) are separate methods. Auth handlers calling `Upsert` can never accidentally zero out computed scores. `MemoryUserStore.Upsert` preserves existing scoring fields; `Recalculate()` calls `UpdateScores` exclusively for score writes.

- [x] **Billing killswitch** — Cloud Function (`infra/billing-killswitch/`, nodejs24, gen2) subscribes to `billing-alerts` Pub/Sub topic; disables billing on both GCP projects when `costAmount > budgetAmount`. Deployed to `us-east5`. Manual step: create a $10 budget in GCP Console and wire to `billing-alerts` topic.

- [x] **Precomputed user scores (leaderboard O(U) reads)** — `internal/recalculator.Recalculate()` computes AcesRadio, Upper90Club, Grouchy points and PredictionCount for every user from scratch and upserts to `UserStore`; triggered after every match final (`MatchPoller.SetOnResultSaved`) and on daily refresh startup (`startDailyRefresh` backfill). Leaderboard and profile handlers now read precomputed values from user docs — O(U) reads instead of O(P×R) per request.

- [x] **Frontend UX fixes (soft launch batch)** — Grouchy™ added to profile page (4-stat grid: Predictions · Aces Radio · Upper 90 Club · Grouchy™, each with rank); guest localStorage predictions auto-submitted on login via `flushGuestPredictions()` in shared `src/guestPredictions.ts` utility — users no longer need to manually resubmit after signing up; auth link styling added (`.auth-alt` was falling back to browser defaults); mobile breakpoints unified at 600px (hamburger was 480px, table cards 600px — now both 600px).
- [x] **Grouchy™ scoring format** — third scoring format based on Columbus margin-of-victory bucket (Win 2+, Win 1, Draw, Lose 1, Lose 2+); +1 for matching bucket, 0 otherwise. `internal/scoring.Grouchy()`; leaderboard and match detail APIs include `grouchyPoints`; all three tables show GROUCHY™ column with mobile sort button; Rules page updated; prediction displayed below handle in match detail (removed separate PICK column). 64 e2e scenarios green.
- [x] **Leaderboard + match detail table redesign** — unified sortable grid table (RANK · PREDICTOR · ACES RADIO · UPPER 90 CLUB · GROUCHY™) replaces separate sections; click column headers to sort; dynamic tied ranks; mobile stacked cards with sort buttons above and only the active format's score shown in gold. `GET /api/leaderboard` response unified to `{entries: [{acesRadioPoints, upper90ClubPoints, grouchyPoints, ...}]}`.
- [x] **TwoOneBot** — `internal/bot` package; `bot:twooonebot` user that predicts Columbus 2-1 (home) or 1-2 (away) on every upcoming match at every refresh and daily tick. Shown on leaderboard and match detail pages as any other user.
- [x] **Match detail page** — `/matches/:matchId` route; per-match predictions leaderboard with sort buttons for Aces Radio and Upper 90 Club; result cards link to detail; upcoming cards do not. `GET /api/matches/:matchId` handler.
- [x] **`localStorage` without try/catch** — `readGuestPredictions`/`writeGuestPredictions` helpers wrap all localStorage access; Safari Private Browsing degrades gracefully.
- [x] **Check `/auth/session` response in sign-in flows** — `postSession(token)` returns `Promise<boolean>`; all callers (`LoginView`, `SignupView`, `App.vue` Google redirect) check `res.ok` before navigating.
- [x] **Handle `json.NewEncoder(w).Encode()` errors** — all six JSON response handlers (`leaderboard`, `matches`, `profile`, `me`, `backfill`, `match_detail`) log encode errors.
- [x] **Add `.dockerignore`** — excludes `.env*`, `node_modules/`, `.git/`, test artifacts, `.claude`; prevents credential leakage in image layers.
- [x] **Run container as non-root** — final Dockerfile stage creates `app` user (uid 1001) and runs the server as that user.
- [x] **Validate `match_id` for emptiness in results and seed handlers** — `handlers/results.go` and `handlers/seed_match.go` reject empty `match_id`/`id` with 400.
- [x] **`inject()` without default** — `inject('currentUser', ref(null))` explicit defaults in `MatchesView`, `ProfileView`.
- [x] **Add loading and error states to views** — `LeaderboardView`, `MatchesView`, `ProfileView`, `MatchDetailView` all show a loading indicator during fetch and an error message on non-ok responses.
- [x] **Guard admin endpoints in prod** — `AdminAuth` middleware (X-Admin-Key, `crypto/subtle` compare) wraps `/admin/results`, `/admin/poll-scores`, `/admin/refresh-matches`, `/admin/backfill-users`. Server refuses to start if `ADMIN_KEY` unset in prod. **Deploy note: set `ADMIN_KEY` in Cloud Run (staging + prod) before merging.**
- [x] **Add `Secure` flag to session cookie** — `Secure: os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") == ""` in `writeSessionCookie` and `Logout`; off in local dev, on in prod.
- [x] **Handle `Save()` error in predictions handler** — returns 500 and logs on Firestore write failure.
- [x] **Server-side goals range validation** — predictions and results handlers reject `home_goals`/`away_goals` outside 0–99 with 400.
- [x] **Handle `r.ParseForm()` errors** — all 5 handlers already return 400 on parse failure; regression tests added.
- [x] **Discard `json.Marshal` error in `serveFirebaseConfig`** — marshaling `map[string]string` cannot fail; existing test covers the response.
- [x] **STATUS_DELAYED support** — blinking red DELAYED badge; match moves to new "Now Playing" section above Upcoming; server rejects predictions with 403; no Predict/Unlock buttons. Confirmed live vs LA Galaxy 2026-04-22.
- [x] **Match ordering stability** — `sort.Slice` by `Kickoff` ascending in `APIList`; `completedMatches` Vue computed sorts descending explicitly; no longer relies on ESPN return order.
- [x] **Client-side kickoff lock** — reactive `nowMs` ref updated each countdown tick; `isLocked()` gates Predict/Unlock at kickoff time without reload; covers STATUS_DELAYED and state='in' too.
- [x] **Now Playing section** — in-progress (state='in') and delayed matches shown above Upcoming in dedicated section; LIVE badge pulses gold, DELAYED badge blinks red; no prediction inputs on these cards.
- [x] **Leaderboard: show users with ≥1 prediction at 0 pts before results land** — seeded from all predictions on every request; smoke/admin accounts never appear since they don't predict.
- [x] **Leaderboard: disable profile link for legacy handle-only users** — `hasProfile` field added to leaderboard API response; false when user has no `UserStore` entry (predates UserStore). Frontend renders plain `<span>` instead of `<RouterLink>` to avoid a 404 profile page.
- [x] **Profile page** — `/profile/:userID` shows handle, location, prediction count, and leaderboard standing (points + rank for both formats); edit form on own profile only; leaderboard handles link to profiles; location field added to `POST /auth/handle`. Full Industrial Black & Gold Brutalism styling: stats grid with gold top borders, DM Mono values, Barlow Condensed handle.
- [x] **Staging smoke cleanup** — switched to permanent accounts only; no more account creation in smoke tests; `users` collection no longer accumulates stale entries per CI run.
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
- [x] **CI/CD pipeline redesign** — GCS frontend artifacts (`sha-{SHA}.zip` / `latest.zip` / `prod.zip` / `prod-previous.zip`); automatic prod rollback on failure; manual rollback workflow dispatch; concurrency cancel on develop; `prod`/`prod-previous` Docker + GCS tags; fallback to `latest` when exact SHA not found. Tested: fallback path exercised on first prod deploy (md-only commit had no artifacts); manual rollback tested — Cloud Run rollback succeeded, Firebase guard for missing `prod-previous.zip` added after first-deploy edge case discovered
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
