# Backlog

## Done

- [x] Go server with ESPN match fetching
- [x] Firestore prediction store
- [x] AcesRadio scoring engine (+15/+10/−15/0)
- [x] Upper 90 Club scoring engine (+1 result, +1 Columbus goals, stacking)
- [x] Leaderboard (both formats, JSON API + Vue view)
- [x] Firebase Auth — Email/Password (custom form, no FirebaseUI)
- [x] Session cookies (`HttpOnly`, base64 JSON)
- [x] Vue 3 SPA: MatchesView, LoginView, LeaderboardView, AppHeader
- [x] BDD e2e suite — 10/10 Playwright scenarios green
- [x] Vite dev proxy for local development
- [x] Industrial Black & Gold Brutalism design applied
- [x] ESPN date parsing fix (`2026-04-12T23:00Z` no-seconds format)
- [x] Match listings — upcoming (next 7 days) + results with scores inline
- [x] Multi-competition support — MLS, US Open Cup, Leagues Cup, CONCACAF Champions

---

## Before First Deploy

- [ ] **FirestoreResultStore** — results currently in-memory only; don't persist across restarts. Implement alongside MemoryResultStore.
- [ ] **GCP/Firebase env vars** — add `FIREBASE_API_KEY`, `FIREBASE_AUTH_DOMAIN`, `FIREBASE_PROJECT_ID` to `.env` and Cloud Run service env
- [ ] **Cloud Run deploy** — first manual `gcloud run deploy crew-predictions --source . --region us-east5`
- [ ] **Firebase Hosting** — `firebase deploy --only hosting` for static assets

---

## Post-Launch

- [ ] **Real-data scoring accuracy test** — e2e scenario using actual 2025 Columbus Crew match results to validate the scoring engine against real outcomes. Get match data from user before writing.
- [ ] **GitHub Actions CI/CD** — push to `main` → test + deploy. Deferred until first successful manual deploy establishes the baseline.
- [ ] **Google OAuth** — add as second sign-in option alongside email/password
- [ ] **Match result entry UI** — admin page for podcast hosts to enter final scores (currently API-only via `POST /admin/results`)
- [ ] **Prediction locking** — enforce kickoff time deadline server-side (currently no lock)
- [ ] **Firestore match cache** — cache ESPN results in Firestore daily instead of fetching live on every request

---

## Decisions Made / Won't Do

- **Bluesky AT Protocol auth** — dropped for v1. Complex, adds no value for first release.
- **FirebaseUI** — dropped. Incompatible with `firebase@11` (requires `^9||^10`).
- **templ/HTMX/Alpine.js** — replaced by Vue 3 SPA. The `templates/` package still exists for the leaderboard HTML handler but is effectively legacy.
- **Frontend subdirectory** — kept as single project root (`package.json` at root, `src/` for Vue).
