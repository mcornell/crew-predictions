# Backlog

## Up Next

### Security / Critical

1. [ ] **Guard admin endpoints in prod** — `/admin/results`, `/admin/poll-scores`, `/admin/refresh-matches`, and `/admin/backfill-users` are fully unauthenticated; anyone can POST fake results and corrupt the leaderboard. Add a shared-secret middleware: check an `X-Admin-Key` header against a Cloud Run env var secret. Only `/admin/reset` and `/admin/seed-*` are already TEST_MODE-gated. (`main.go:131–151`, `handlers/results.go`, `handlers/refresh_matches.go`, `handlers/poll_scores.go`, `handlers/backfill.go`)

2. [ ] **Add `Secure` flag to session cookie** — `writeSessionCookie` sets `HttpOnly` and `SameSite` but omits `Secure: true`; cookie transmits over plain HTTP. Set `Secure: os.Getenv("FIREBASE_AUTH_EMULATOR_HOST") == ""` so it's off in the emulator and on in prod. Same fix needed in the logout cookie. (`handlers/handle.go:54`, `handlers/auth.go:12`)

3. [ ] **Handle `Save()` error in predictions handler** — Firestore write failure is silently discarded; user gets a 302 success redirect while the prediction is lost. Return 500 if `h.store.Save()` errors. (`handlers/predictions.go:66`)

### Bugs

4. [ ] **Server-side goals range validation** — `strconv.Atoi` accepts negative numbers and values in the billions; HTML `min`/`max` is client-only. Reject `home_goals` or `away_goals` outside 0–99 with 400. (`handlers/predictions.go:34`, `handlers/results.go:21`)

5. [ ] **`localStorage` without try/catch** — `getItem`/`setItem`/`JSON.parse` on `localStorage` throw `SecurityError` in Safari Private Browsing and crash the entire `onMounted` handler. Wrap all `localStorage` access in try/catch and degrade gracefully. (`src/views/MatchesView.vue:187,213`)

6. [ ] **Check `/auth/session` response in sign-in flows** — `postSession(token)` never checks `res.ok`; if the server returns 401 the cookie is never set but the client still routes to `/matches`, leaving the user in a broken logged-in-looking-but-not state. (`src/views/LoginView.vue:38`, `src/views/SignupView.vue:35`, `src/App.vue:29`)

7. [ ] **Handle `r.ParseForm()` errors** — called without checking the return value in five handlers; a malformed or oversized body silently yields `""` for all fields, causing misleading downstream errors. Return 400 on parse failure. (`handlers/predictions.go:28`, `handlers/session.go:41`, `handlers/handle.go:25`, `handlers/results.go:19`, `handlers/seed.go:19`)

8. [ ] **Handle `json.NewEncoder(w).Encode()` errors** — encoding failures send truncated JSON with a 200 status because headers are already written; at minimum log the error. (`handlers/leaderboard.go`, `handlers/matches.go`, `handlers/profile.go`, `handlers/me.go`, `handlers/backfill.go`)

9. [ ] **Discard `json.Marshal` error in `serveFirebaseConfig`** — `cfg, _ := json.Marshal(...)` sends `null` to every browser on failure, breaking Firebase init on every page load. (`main.go:246`)

### Infrastructure / Medium

10. [ ] **Add `.dockerignore`** — without it `COPY . .` pulls `.env` files, `node_modules/`, `.git/`, and test fixtures into the Docker build context; local `.env` credentials can end up in image layers.

11. [ ] **Run container as non-root** — Dockerfile final stage has no `USER` directive; the Go server runs as root. Add `RUN useradd -r -u 1001 app && USER app`.

12. [ ] **Validate `match_id` for emptiness in results and seed handlers** — a result or seed saved with an empty `match_id` persists in the store and never matches any prediction lookup, silently poisoning the store. (`handlers/results.go:18`, `handlers/seed.go:19`)

13. [ ] **`inject()` without default** — `inject<Ref<...>>('currentUser')` returns `undefined` if called outside the provide tree; use `inject('currentUser', ref(null))` to make the default explicit. (`src/views/MatchesView.vue:115`, `src/views/ProfileView.vue:72`)

### Low

14. [ ] **Prod smoke suite** — unauthenticated-only scenarios (app loads, leaderboard/matches API responds, Vue hydrates); replaces current `curl` liveness check in `deploy-prod`.

15. [ ] **Add loading and error states to views** — Leaderboard, Profile, and Matches views render a blank page during fetch and on error; users see nothing with no feedback. Add a loading indicator and an error message on non-ok responses.

16. [ ] **Cache leaderboard scoring** — currently recalculated on every request; fine now but will need in-memory caching or pre-computation at scale.

17. [ ] **Remove stale `handle` from predictions** — `handle` on prediction documents is legacy; leaderboard now sources display names from `UserStore` by `userID`. Stop writing it on new predictions and remove the `p.Handle` fallback in the leaderboard once confirmed no UID-less predictions exist in prod.

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
