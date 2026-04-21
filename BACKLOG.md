# Backlog

## Done

- [x] Go server with ESPN match fetching
- [x] Firestore prediction store
- [x] AcesRadio scoring engine (+15/+10/−15/0)
- [x] Upper 90 Club scoring engine (+1 result, +1 Columbus goals, stacking)
- [x] Leaderboard (both formats, JSON API + Vue view)
- [x] Firebase Auth — Email/Password sign-in (Google SSO still pending)
- [x] Session cookies (`HttpOnly`, base64 JSON)
- [x] Vue 3 SPA: MatchesView, LoginView, LeaderboardView, AppHeader
- [x] BDD e2e suite — 24/24 Playwright scenarios green
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
- [x] GitHub Actions CI/CD — push runs Go + Vitest + Playwright; main deploys to Cloud Run + Firebase Hosting via Workload Identity Federation
- [x] Sign-up flow — `/signup` view + `firebase.signUp` (createUserWithEmailAndPassword), reuses `/auth/session` cookie flow
- [x] Google SSO — `signInWithPopup(GoogleAuthProvider)` button on both `/login` and `/signup`; Google provider + OAuth web client secret were enabled in the Firebase Console in an earlier session

---

## Before First Deploy

- [x] **FirestoreResultStore** — results persist across restarts
- [x] **Prediction locking** — enforce kickoff time deadline server-side; predictions must be submitted before kickoff
- [x] **GCP/Firebase env vars** — `FIREBASE_API_KEY`, `FIREBASE_AUTH_DOMAIN`, `FIREBASE_PROJECT_ID`, `GOOGLE_CLOUD_PROJECT` set on Cloud Run service; old OAuth vars removed
- [x] **Cloud Run deploy** — Go server deployed to Cloud Run
- [x] **Firebase Hosting** — SPA shell + static assets; rewrites `/api/**`, `/auth/**`, `/admin/**` to Cloud Run; SPA fallback for all other routes

---

## Next Up (in order)

1. [ ] **Mobile responsive layout** — match cards have no mobile breakpoint and collapse completely on all phone sizes (tested: iPhone 15 390px, Samsung Galaxy S24 360px, Pixel 8 412px); header also wraps at narrow widths. Full mobile-first pass needed:
   - Match card layout: stack home/away vertically on mobile, keep score prominent
   - Header: collapse nav or reduce to icon-only on narrow widths
   - Score inputs: increase tap target size and font size on mobile
   - Typography: team names and scores too small to read comfortably on 360–412px widths

2. [ ] **UX fixes — broken flows**
   - [ ] Predict while logged out silently fails — clicking Predict fires a 401 but shows no feedback; redirect to `/login` or show an inline "sign in to predict" prompt
   - [ ] Unknown routes render blank — `/notapage` shows just the header on a black void; add a 404/NotFound view with a link home
   - [ ] Score inputs visible to logged-out users — users can type scores before being asked to sign in; either hide inputs or replace Predict button with "Sign in to predict" until authenticated
   - [ ] Profile display name not pre-populated — the input on `/profile` is empty even when the user already has a handle set; load current handle on mount

3. [ ] **UX gaps — missing content**
   - [ ] Leaderboard empty state is a dead end — "No predictions scored yet" with no explanation of what will appear or how scoring works; add context
   - [ ] No scoring rules explanation anywhere — new users have no idea what Aces Radio or Upper 90 Club scoring means; add a "How it works" section or `/rules` page

4. [ ] **Page `<title>` per route** — currently always "Crew Predictions" regardless of route; each view should set a meaningful `<title>` (e.g. "Leaderboard — Crew Predictions")

5. [ ] **Auth UX polish** — remaining sub-items:
   - [x] Login/signup cross-links
   - [x] Error-state differentiation (sign-up only; login stays generic for security)
   - [x] Verify logout UI — fixed broken `/logout` href (Go route is `/auth/logout`), new BDD scenario asserts clicking Sign out logs the user out
   - [x] Password-reset flow — `/reset` view, `sendPasswordResetEmail`, "Forgot password?" link on login; emulator requires existing user so e2e seeds one first
   - [x] Display name / profile page — `/profile` route + `ProfileView`; session `handle` now prefers Firebase `name` claim (displayName) over email; user handle in header is a link to `/profile`; `waitForCurrentUser()` helper waits for Firebase auth state restoration before calling `updateProfile`
   - [x] Email verification — banner shown to unverified users; emailVerified surfaced through FirebaseToken → session cookie → /api/me → App.vue
2. [x] **Staging Cloud Run + artifact promotion** — develop builds Docker image tagged by SHA, pushes to Artifact Registry, deploys to `crew-predictions-staging`, smoke tests staging, uploads `dist/` as artifact. main promotes same image to prod (no rebuild), downloads matching `dist/` artifact, deploys Firebase Hosting, smoke tests prod.
3. [ ] **Custom domain migration** — Firebase Hosting custom domain + Cloud Run domain mapping. Update `authDomain` and OAuth redirect URIs.

---

## Polish (prioritize later)

- [ ] **Password field `autocomplete` attribute** — login and signup password inputs missing `autocomplete="current-password"` / `"new-password"`; breaks password manager autofill and triggers browser console warnings
- [ ] **Profile page needs context** — currently just a display name form floating in space; add current prediction count, scoring summary, or other stats to make it worth visiting

---

## Data & Polling

- [ ] **Real-data scoring accuracy test** — e2e scenario using actual 2025 Columbus Crew match results to validate the scoring engine against real outcomes. Get match data from user before writing.
- [ ] **Firestore match cache + score polling** — cache ESPN schedule in Firestore. Weekly refresh fires Tuesday midnight ET (cron) regardless of whether there's a match that week. When `kickoff + 2h <= now` and match not yet `STATUS_FINAL`, poll ESPN every ~5 min and write to ResultStore when final, then stop. ESPN already returns `status.type.name`.

---

## Test Infrastructure

- [ ] **Playwright smoke suite for prod** — identify a small tagged subset of e2e scenarios that can run against the live prod URL after deploy (replaces the current `curl` liveness check in `deploy-prod`)
- [ ] **Fix Vue Router warnings in unit tests** — test `makeRouter()` factories only register a few routes; components with `<router-link>` to unregistered paths (e.g. `/login`, `/signup`, `/reset`) trigger "No match found" warnings. Add stub routes to silence them.
- [ ] **Per-worker server isolation** — if the e2e suite grows large enough that serial execution is too slow, give each Playwright worker its own server instance (separate ports) so parallel runs don't share in-memory state.

---

## Decisions Made / Won't Do

- **Match result entry UI** — admin page not needed; `POST /admin/results` API is sufficient for now.
- **Bluesky AT Protocol auth** — dropped for v1. Complex, adds no value for first release.
- **FirebaseUI** — dropped. Incompatible with `firebase@11` (requires `^9||^10`).
- **templ/HTMX/Alpine.js** — replaced by Vue 3 SPA; `templates/` package deleted.
- **GCP Cloud Build** — GitHub Actions preferred; simpler config, already where the code lives, Workload Identity Federation handles GCP auth without stored keys.
- **Frontend subdirectory** — kept as single project root (`package.json` at root, `src/` for Vue).
