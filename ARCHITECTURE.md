# Architecture

## Stack

| Layer | Technology | Why |
|---|---|---|
| Language | Go | Fast cold starts, single binary, good GCP SDKs |
| Compute | GCP Cloud Run | Multi-route web app, free tier (2M req/mo), source deploy |
| Frontend | Vue 3 + TypeScript + Vite | Simpler than React, user preference over templ/HTMX |
| Database | Firestore | GCP always-free (50k reads/20k writes per day) |
| Auth | Firebase Auth — Email/Password | Custom form; FirebaseUI dropped (incompatible with firebase@11) |
| Static assets | Firebase Hosting | GCP-native, pairs with Cloud Run |
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
    └── /assets/* ─────────────────      └── ESPN API (match data, cached)
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
| `templates/` | templ components (leaderboard HTML — legacy, being phased out) |

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
| `DELETE` | `/admin/reset` | — | Reset in-memory stores (TEST_MODE=1 only) |
| `POST` | `/admin/results` | — | Record a final match result for scoring |

**Form data convention:** `POST /api/predictions` and `POST /auth/session` use `application/x-www-form-urlencoded` (Go's `r.ParseForm()`). Send via `URLSearchParams`, not JSON.

---

## Auth Flow

1. User submits email + password on `/login`
2. Vue calls `signInWithEmailAndPassword` (Firebase Auth SDK → emulator in dev, real in prod)
3. Gets ID token → POSTs to `POST /auth/session` as form data
4. Go server verifies token via Firebase Admin SDK, sets `HttpOnly` session cookie
5. Session cookie = base64-encoded JSON `{ userID, handle, provider }`
6. Subsequent requests: Go reads cookie via `UserFromSession(r)`

---

## Vue SPA

Entry: `src/main.ts` → loads `/auth/config.js` → mounts Vue app

| File | Route | Purpose |
|---|---|---|
| `src/views/MatchesView.vue` | `/` `/matches` | Upcoming matches + prediction inputs |
| `src/views/LoginView.vue` | `/login` | Email/password sign-in form |
| `src/views/LeaderboardView.vue` | `/leaderboard` | Aces Radio + Upper 90 Club rankings |
| `src/components/AppHeader.vue` | (all) | Nav header with auth state |
| `src/App.vue` | — | Root: fetches `/api/me` on mount + route change |
| `src/firebase.ts` | — | Firebase Auth SDK init + `signIn` helper |

**CSS:** `src/style.css` — design tokens as CSS variables, imported in `src/main.ts`. Not yet written.

---

## Data Model (Firestore)

```
predictions/{predictionId}
  matchID:    string
  userID:     string   // "firebase:{uid}"
  handle:     string   // display name at time of prediction
  homeGoals:  int
  awayGoals:  int
```

Results are currently stored in `MemoryResultStore` (not persisted). See BACKLOG.md.
