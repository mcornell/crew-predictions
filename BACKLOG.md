# Backlog

## Up Next

1. [ ] **Profile page needs context** — currently just a display name form floating in space; add current prediction count, scoring summary, or other stats to make it worth visiting.

2. [ ] **Custom domain migration** — Firebase Hosting custom domain + Cloud Run domain mapping. Update `authDomain` and OAuth redirect URIs.

3. [ ] **Allow unlocking a pick up to kickoff** — after a user locks in a prediction, let them edit and re-submit it any time before kickoff. Server already enforces the 403 after kickoff; this is purely a UI unlock flow.
   - **Countdown to lock** — show a live countdown on each match card so users know how long they have. The tricky part: browser clocks drift and can be manipulated, so the countdown is cosmetic only. Cheat-proofing options: (a) fetch server time once on page load and compute an offset (`serverTime - clientTime`), then apply that offset to all countdowns — cheap and good enough for casual play; (b) re-validate the kickoff time on submit via the server's 403 (already in place) so even a manipulated countdown can't actually lock in a late pick.

---

## Data & Polling

- [ ] **Real-data scoring accuracy test** — e2e scenario using actual 2025 Columbus Crew match results to validate the scoring engine against real outcomes. Get match data from user before writing.
- [ ] **Firestore match cache + score polling** — cache ESPN schedule in Firestore. Weekly refresh fires Tuesday midnight ET (cron) regardless of whether there's a match that week. When `kickoff + 2h <= now` and match not yet `STATUS_FINAL`, poll ESPN every ~5 min and write to ResultStore when final, then stop. ESPN already returns `status.type.name`.

---

## Test Infrastructure

- [ ] **Playwright smoke suite for prod** — identify a small tagged subset of e2e scenarios that can run against the live prod URL after deploy (replaces the current `curl` liveness check in `deploy-prod`)
- [ ] **Per-worker server isolation** — if the e2e suite grows large enough that serial execution is too slow, give each Playwright worker its own server instance (separate ports) so parallel runs don't share in-memory state.

---

## Decisions Made / Won't Do

- **Match result entry UI** — admin page not needed; `POST /admin/results` API is sufficient for now.
- **Bluesky AT Protocol auth** — dropped for v1. Complex, adds no value for first release.
- **FirebaseUI** — dropped. Incompatible with `firebase@11` (requires `^9||^10`).
- **templ/HTMX/Alpine.js** — replaced by Vue 3 SPA; `templates/` package deleted.
- **GCP Cloud Build** — GitHub Actions preferred; simpler config, already where the code lives, Workload Identity Federation handles GCP auth without stored keys.
- **Frontend subdirectory** — kept as single project root (`package.json` at root, `src/` for Vue).

---

## Done

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
