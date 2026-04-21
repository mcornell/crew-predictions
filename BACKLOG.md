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

1. [x] **Mobile responsive layout** — shipped: cards stack correctly at 390px and 412px; header stays single-row; Predict button is full-width at 48px. Remaining polish:
   - [ ] Team name truncation on 360px (Galaxy S24) — "COLUMBUS CREW" clips to "COLUMBUS C..." at the narrowest CSS viewport; physical screen is 1080px but CSS pixels are 360px due to 3× device pixel ratio, so there's no way to use more pixels — need a layout solution (smaller inputs, abbreviated names, or two-line team display)

2. [x] **UX fixes — broken flows** — all shipped
   - [x] Predict while logged out silently fails — redirects to `/login` on 401
   - [x] Unknown routes render blank — 404 NotFoundView added with link home
   - [x] Score inputs visible to logged-out users — replaced with "Sign in to predict" button; inputs hidden until authenticated
   - [x] Profile display name not pre-populated — watches currentUser and seeds input on load
   - [ ] Google sign-in popup blocked on mobile — `signInWithPopup` is blocked by mobile browsers (tested: real device); switch to `signInWithRedirect` everywhere (desktop too — it does a full-page nav to Google and back, seamless and widely accepted). Drop the popup entirely for simplicity. Firebase's `getRedirectResult()` must be called on app load to complete the redirect flow.

3. [ ] **Match official Columbus Crew brand colors** — current `--gold: #ffc20e` is approximate; audit the official Columbus Crew website for exact hex values for gold, black, and any secondary palette colors (e.g. dark charcoal, accent red). Update CSS variables in `src/style.css`. Use the Playwright browser plugin to pull colors directly from the live site rather than guessing.

4. [ ] **UX gaps — missing content**
   - [ ] Leaderboard empty state is a dead end — "No predictions scored yet" with no explanation of what will appear or how scoring works; add context
   - [ ] No scoring rules explanation anywhere — new users have no idea what Aces Radio or Upper 90 Club scoring means; add a "How it works" section or `/rules` page

5. [ ] **Page `<title>` per route** — currently always "Crew Predictions" regardless of route; each view should set a meaningful `<title>` (e.g. "Leaderboard — Crew Predictions")

6. [x] **Auth UX polish** — all done
   - [x] Login/signup cross-links
   - [x] Error-state differentiation (sign-up only; login stays generic for security)
   - [x] Verify logout UI
   - [x] Password-reset flow
   - [x] Display name / profile page
   - [x] Email verification banner

7. [x] **Staging Cloud Run + artifact promotion**

8. [ ] **Custom domain migration** — Firebase Hosting custom domain + Cloud Run domain mapping. Update `authDomain` and OAuth redirect URIs.

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
