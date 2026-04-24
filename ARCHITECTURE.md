# Architecture

## Stack

| Layer | Technology | Why |
|---|---|---|
| Language | Go | Fast cold starts, single binary, good GCP SDKs |
| Compute | GCP Cloud Run | Multi-route web app, free tier (2M req/mo) |
| Frontend | Vue 3 + TypeScript + Vite | Simpler than React, user preference over templ/HTMX |
| Database | Firestore | GCP always-free (50k reads/20k writes per day) |
| Auth | Firebase Auth — Email/Password + Google SSO | Custom form; FirebaseUI dropped (incompatible with firebase@11) |
| Static assets | Firebase Hosting | GCP-native CDN; rewrites API/auth paths to Cloud Run |
| Match data | ESPN unofficial API | Free, covers MLS/Columbus Crew |

**Firestore region:** `us-east5` (Columbus, Ohio — obviously)

---

## How the Pieces Fit

```
Browser (Vue SPA)
    │
    ├── /api/*  ──────────────────► Go server (Cloud Run :8080)
    │                                    │
    ├── /auth/* ──────────────────►      ├── Firebase Admin SDK (token verification)
    │                                    ├── Firestore (predictions, results)
    └── /assets/* ─────────────────      └── ESPN API (match data)
         (Vite build → dist/)
```

**Local dev:** Vite dev server (:5173) proxies `/api`, `/auth`, `/admin` to Go server (:8080). Firebase Auth + Firestore emulators run on :9099 and :8081.

**E2e test server** runs on :8082 (`PORT=8082` in Playwright's `webServer` env) so `npm test` can run alongside `dev.sh` without port conflicts.

---

## Go Server

Entry point: `cmd/server/main.go`

| Package | Responsibility |
|---|---|
| `internal/handlers` | HTTP handlers — matches, predictions, leaderboard, profile, auth, session, handle update, match detail |
| `internal/repository` | Data access — Firestore and in-memory stores; `WriteThroughMatchStore` |
| `internal/scoring` | Scoring engines — AcesRadio, Upper90Club, and Grouchy |
| `internal/espn` | ESPN API client — fetches upcoming Crew matches |
| `internal/poll` | Score polling — `MatchPoller` schedules per-match kickoff timers; `PollOnce` for manual/test triggers |
| `internal/bot` | TwoOneBot — predicts Columbus 2-1 (home) / 1-2 (away) on every upcoming match at each refresh and daily tick |
| `internal/models` | Domain types |

---

## API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/matches` | optional | Upcoming matches + caller's predictions |
| `GET` | `/api/matches/:matchId` | none | Match detail: match info + all predictions with per-format scores + `scoringFormats` array |
| `POST` | `/api/predictions` | required | Submit a score prediction (form data) |
| `GET` | `/api/leaderboard` | none | `{entries: [{userID, handle, acesRadioPoints, upper90ClubPoints, grouchyPoints, hasProfile}]}` sorted by Aces Radio desc; all users with ≥1 prediction appear (0 pts until results land); `hasProfile: bool` false for legacy handle-only users |
| `GET` | `/api/me` | optional | Current session user `{userID, handle}` or 401; lazily upserts user to `UserStore` |
| `GET` | `/api/profile/:userID` | required | Public profile: handle, location, predictionCount, Aces Radio + Upper 90 Club + Grouchy™ standing |
| `POST` | `/auth/handle` | required | Update display name + location; upserts to `UserStore`, rewrites session cookie |
| `POST` | `/auth/session` | — | Exchange Firebase ID token for session cookie (form data) |
| `GET` | `/auth/logout` | — | Clear session cookie, redirect to /matches |
| `GET` | `/auth/config.js` | — | Firebase client config as JS (`window.__firebaseConfig`) |
| `POST` | `/admin/refresh-matches` | — | Fetch matches from ESPN, populate match cache, reschedule pollers |
| `POST` | `/admin/poll-scores` | — | Trigger a score poll immediately (fetch ESPN, update store, write terminal results) |
| `DELETE` | `/admin/reset` | — | Reset in-memory stores (TEST_MODE=1 only) |
| `POST` | `/admin/results` | — | Record a final match result for scoring |
| `POST` | `/admin/backfill-users` | — | One-time: populate `users` collection from existing predictions |
| `POST` | `/admin/seed-match` | — | Inject a fixture match (TEST_MODE=1 only) |
| `POST` | `/admin/seed-prediction` | — | Inject a fixture prediction (TEST_MODE=1 only) |

**Form data convention:** `POST /api/predictions` and `POST /auth/session` use `application/x-www-form-urlencoded` (Go's `r.ParseForm()`). Send via `URLSearchParams`, not JSON.

---

## Match Cache

The server holds a `MatchStore` backed by `WriteThroughMatchStore` — an in-memory primary (fast reads) wrapped around a `FirestoreMatchStore` secondary (durable writes). On startup, stored matches are loaded from Firestore into memory so match data survives restarts without waiting for the ESPN fetch. In `TEST_MODE=1`, a bare `MemoryMatchStore` is used and the seed handler writes directly to it.

ESPN data is fetched via `internal/espn.FetchCrewMatches`, which hits four league endpoints (MLS, US Open Cup, Leagues Cup, CONCACAF Champions). The HTTP base URL is injectable for testing — `fetchCrewMatchesFrom(base)` is covered by `httptest.Server` + captured fixture JSON.

**Daily refresh:** `startDailyRefresh` fires at 4am ET on startup and every subsequent 24h. It fetches ESPN, updates `MatchStore` (writing through to Firestore), and calls `poller.Reset(matches)` to reschedule all match pollers from fresh data.

**Score polling:** `internal/poll.MatchPoller` schedules a `time.AfterFunc` at each match's kickoff time. When the timer fires, the match enters the active set and `Tick()` polls ESPN every 2 minutes. On a terminal status (`STATUS_FULL_TIME` / `STATUS_FINAL_AET` / `STATUS_FINAL_PEN`), the result is written to `ResultStore` and the match is deactivated. Matches with unknown/postponed status stay active until the next 4am reset clears them. Matches loaded from Firestore with a past kickoff are scheduled at zero delay (immediate catch-up polling).

---

## Auth Flow

### Email/Password
1. User submits email + password on `/login` (or creates account on `/signup`)
2. Vue calls `signInWithEmailAndPassword` / `createUserWithEmailAndPassword`
3. Gets ID token → POSTs to `POST /auth/session` as form data
4. Go server verifies token via Firebase Admin SDK, sets `HttpOnly` session cookie
5. Session cookie = base64-encoded JSON `{ userID, handle, provider }`
6. Subsequent requests: Go reads cookie via `UserFromSession(r)`

### Google SSO
1. User clicks "Sign in with Google" on `/login` or `/signup`
2. Vue calls `signInWithRedirect` (redirect — not popup — for mobile compatibility)
3. After redirect back, `App.vue` calls `getRedirectResult()` on mount
4. Same session cookie flow as email/password from step 3

---

## Vue SPA

Entry: `src/main.ts` → loads `/auth/config.js` → mounts Vue app

| File | Route | Purpose |
|---|---|---|
| `src/views/MatchesView.vue` | `/` `/matches` | Upcoming matches + prediction inputs; completed matches reversed chronological |
| `src/views/LoginView.vue` | `/login` | Email/password sign-in + Google SSO |
| `src/views/SignupView.vue` | `/signup` | New account creation (email/password + Google SSO) |
| `src/views/ResetView.vue` | `/reset` | Password reset request |
| `src/views/LeaderboardView.vue` | `/leaderboard` | Unified sortable table (RANK · PREDICTOR · ACES RADIO · UPPER 90 CLUB · GROUCHY™); click column headers to sort; handles link to `/profile/:userID`; mobile: stacked cards, sort buttons above, active format score shown only |
| `src/views/ProfileView.vue` | `/profile/:userID` | Public profile; 4-stat grid (Predictions · Aces Radio · Upper 90 Club · Grouchy™, each with rank); edit form shown only on own profile |
| `src/views/MatchDetailView.vue` | `/matches/:matchId` | Unified sortable table (RANK · PREDICTOR+PICK · ACES RADIO · UPPER 90 CLUB · GROUCHY™); click column headers to sort; prediction shown below handle; result cards link here, upcoming cards do not; mobile: stacked cards with active format score |
| `src/views/RulesView.vue` | `/rules` | Scoring format explainer |
| `src/views/NotFoundView.vue` | `*` | 404 catch-all |
| `src/components/AppHeader.vue` | (all) | Nav header; desktop nav + hamburger drawer at ≤600px |
| `src/App.vue` | — | Root: fetches `/api/me` on mount + route change; handles Google redirect result |
| `src/firebase.ts` | — | Firebase Auth SDK init + `signIn` / `signInWithGoogle` helpers |
| `src/guestPredictions.ts` | — | Shared guest prediction helpers: `readGuestPredictions`, `writeGuestPredictions`, `flushGuestPredictions` (auto-submits localStorage predictions on login; called from LoginView, SignupView, App.vue Google redirect) |

**CSS:** `src/style.css` — Industrial Black & Gold Brutalism design tokens as CSS variables, imported in `src/main.ts`.

---

## Data Model (Firestore)

```
predictions/{predictionId}
  MatchID:    string   // PascalCase — matches Go struct field names (no firestore: tags)
  UserID:     string   // "firebase:{uid}" or "bot:twooonebot"
  Handle:     string   // display name at time of prediction
  HomeGoals:  int
  AwayGoals:  int

results/{matchID}
  homeScore:  int
  awayScore:  int

users/{userID}
  handle:     string   // current display name (source of truth)
  provider:   string   // "google.com", "password", etc.
  location:   string   // optional, user-supplied (e.g. "Columbus, OH")

matches/{matchID}
  homeTeam:   string
  awayTeam:   string
  kickoff:    timestamp
  status:     string
  homeScore:  string
  awayScore:  string
  state:      string   // "pre" / "in" / "post"
```

---

## CI/CD Pipeline

All deploys flow through GitHub Actions (`.github/workflows/ci.yml`).

```
push to develop
    │
    ├── test job ──────── Go unit tests (coverage report + pass/skip/fail summary → step summary)
    │                    Vue unit tests (v8 coverage table → step summary)
    │                    TypeScript type check
    │                    e2e BDD suite — Playwright, Firebase emulators (JUnit → step summary via dorny/test-reporter@v3)
    │
    └── deploy-staging ── Docker build → Artifact Registry
                          Cloud Run deploy (crew-predictions-staging, us-east5)
                          Firebase Hosting deploy → crew-predictions-staging.web.app
                          Smoke test suite (real staging URL, permanent accounts only — no account creation)
                          Frontend artifact uploaded (retained 90 days)

push to main (merge from develop)
    │
    └── deploy-prod ─────  Promote Docker image from staging artifact (no rebuild)
                           Cloud Run deploy (crew-predictions, us-east5)
                           Firebase Hosting deploy → crew-predictions.web.app
                           curl liveness check
```

**Artifact promotion:** prod deploys reuse the Docker image built for staging — no rebuild on merge. The frontend dist is downloaded from the staging workflow artifact and deployed directly.

**GCP auth:** Workload Identity Federation (no stored service account keys).

---

## Environments

| Environment | Frontend URL | Cloud Run | GCP Project |
|---|---|---|---|
| Prod | https://crew-predictions.web.app | `crew-predictions` service, `us-east5` | `crew-predictions` |
| Staging | https://crew-predictions-staging.web.app | `crew-predictions-staging` service, `us-east5` | `crew-predictions-staging` |
| Local | http://localhost:5173 (Vite proxy) | Go server :8080 | — (emulators) |

**Why staging needs its own GCP project:** Firebase Hosting rewrites to Cloud Run require both to be in the same GCP project. Staging Cloud Run lives in `crew-predictions-staging` so the staging Hosting config can rewrite to it without touching prod infrastructure.

---

## Environment Variables (Cloud Run)

Set via `gcloud run services update <service> --region us-east5 --update-env-vars KEY=VALUE`.

| Variable | Purpose |
|---|---|
| `GOOGLE_CLOUD_PROJECT` | Activates Firestore (predictions + results) |
| `FIREBASE_PROJECT_ID` | Firebase Admin SDK init |
| `FIREBASE_API_KEY` | Served to browser via `/auth/config.js` |
| `FIREBASE_AUTH_DOMAIN` | Served to browser via `/auth/config.js` |

---

## Local Dev Commands

```bash
./dev.sh               # start Firebase emulators (:8081/:9099) + Go server (:8080)
npm run dev            # Vite dev server (:5173) with hot reload + API proxy
go test ./...          # Go unit tests
npm run test:unit      # Vitest unit tests
npm test               # e2e BDD suite — starts own server on :8082; dev.sh can stay running
npm run test:smoke     # smoke suite against staging (STAGING_URL env var)
SMOKE_DEBUG=1 npm run test:smoke  # headed browser + video locally
```

**Note:** `GOOGLE_CLOUD_PROJECT` must NOT be set in the Playwright `webServer` env — it triggers Firestore, which conflicts with the in-memory test store.
