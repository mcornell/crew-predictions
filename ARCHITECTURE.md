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

---

## Go Server

Entry point: `cmd/server/main.go`

| Package | Responsibility |
|---|---|
| `internal/handlers` | HTTP handlers — matches, predictions, leaderboard, auth, session |
| `internal/repository` | Data access — Firestore and in-memory stores |
| `internal/scoring` | Scoring engines — AcesRadio and Upper90Club |
| `internal/espn` | ESPN API client — fetches upcoming Crew matches |
| `internal/models` | Domain types |

---

## API Endpoints

| Method | Path | Auth | Description |
|---|---|---|---|
| `GET` | `/api/matches` | optional | Upcoming matches + caller's predictions |
| `POST` | `/api/predictions` | required | Submit a score prediction (form data) |
| `GET` | `/api/leaderboard` | none | Ranked scores for both formats |
| `GET` | `/api/me` | optional | Current session user `{handle}` or 401 |
| `POST` | `/auth/session` | — | Exchange Firebase ID token for session cookie (form data) |
| `GET` | `/auth/logout` | — | Clear session cookie, redirect to /matches |
| `GET` | `/auth/config.js` | — | Firebase client config as JS (`window.__firebaseConfig`) |
| `POST` | `/admin/refresh-matches` | — | Fetch matches from ESPN and populate the in-memory match cache |
| `DELETE` | `/admin/reset` | — | Reset in-memory stores (TEST_MODE=1 only) |
| `POST` | `/admin/results` | — | Record a final match result for scoring |
| `POST` | `/admin/seed-match` | — | Inject a fixture match (TEST_MODE=1 only) |
| `POST` | `/admin/seed-prediction` | — | Inject a fixture prediction (TEST_MODE=1 only) |

**Form data convention:** `POST /api/predictions` and `POST /auth/session` use `application/x-www-form-urlencoded` (Go's `r.ParseForm()`). Send via `URLSearchParams`, not JSON.

---

## Match Cache

The server holds an in-memory `MatchStore` (populated via `POST /admin/refresh-matches`). In production, call this endpoint after deploy and schedule it weekly via Cloud Scheduler. In `TEST_MODE=1`, the refresh fetcher reads from the seeded store rather than calling ESPN — so e2e tests inject fixtures via `POST /admin/seed-match` and trigger a refresh to populate the cache.

ESPN data is fetched via `internal/espn.FetchCrewMatches`, which hits four league endpoints (MLS, US Open Cup, Leagues Cup, CONCACAF Champions). The HTTP base URL is injectable for testing — `fetchCrewMatchesFrom(base)` is covered by `httptest.Server` + captured fixture JSON.

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
| `src/views/LeaderboardView.vue` | `/leaderboard` | Aces Radio + Upper 90 Club rankings |
| `src/views/ProfileView.vue` | `/profile` | Display name edit |
| `src/views/RulesView.vue` | `/rules` | Scoring format explainer |
| `src/views/NotFoundView.vue` | `*` | 404 catch-all |
| `src/components/AppHeader.vue` | (all) | Nav header; desktop nav + hamburger drawer at ≤480px |
| `src/App.vue` | — | Root: fetches `/api/me` on mount + route change; handles Google redirect result |
| `src/firebase.ts` | — | Firebase Auth SDK init + `signIn` / `signInWithGoogle` helpers |

**CSS:** `src/style.css` — Industrial Black & Gold Brutalism design tokens as CSS variables, imported in `src/main.ts`.

---

## Data Model (Firestore)

```
predictions/{predictionId}
  matchID:    string
  userID:     string   // "firebase:{uid}"
  handle:     string   // display name at time of prediction
  homeGoals:  int
  awayGoals:  int

results/{matchID}
  homeScore:  int
  awayScore:  int
```

---

## CI/CD Pipeline

All deploys flow through GitHub Actions (`.github/workflows/ci.yml`).

```
push to develop
    │
    ├── test job ──────── Go unit tests + Vue unit tests + e2e BDD suite (Firebase emulators)
    │
    └── deploy-staging ── Docker build → Artifact Registry
                          Cloud Run deploy (crew-predictions-staging, us-east5)
                          Firebase Hosting deploy → crew-predictions-staging.web.app
                          Smoke test suite (real staging URL, screenshots always on)
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
npm test               # e2e BDD suite (emulators must be running)
npm run test:smoke     # smoke suite against staging (STAGING_URL env var)
SMOKE_DEBUG=1 npm run test:smoke  # headed browser + video locally
```

**Note:** `GOOGLE_CLOUD_PROJECT` must NOT be set in the Playwright `webServer` env — it triggers Firestore, which conflicts with the in-memory test store.
