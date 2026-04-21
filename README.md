# Crew Predictions

Public humiliation as a service for #Crew96 fandom.

Pick your scores for upcoming Columbus Crew matches, get ranked against other fans using podcast scoring formats, and be wrong in public. It's tradition.

Predictions lock at published kickoff time. No excuses.

---

## Scoring Formats

**Aces Radio**
| Outcome | Points |
|---|---|
| Exact score | +15 |
| Correct result (wrong score) | +10 |
| Flipped same scoreline (predict Crew 3–2, actual opponent 3–2 Crew) | −15 |
| Anything else | 0 |

**Upper 90 Club** — two independent points per match:
| Condition | Points |
|---|---|
| Correct match result (win/draw/loss) | +1 |
| Correct Columbus Crew goal count | +1 |

---

## Running Locally

**Prerequisites:** Go, Node.js, Firebase CLI

```bash
# Install dependencies
npm install

# Terminal 1 — emulators + Go server
./dev.sh

# Terminal 2 — Vue dev server (hot reload)
npm run dev
```

Open `http://localhost:5173`.

**Run tests:**
```bash
go test ./...          # Go unit tests
npm run test:unit      # Vue component tests (Vitest)
npm test               # Full e2e suite (Playwright BDD, runs against local emulator)
npm run test:smoke     # Staging smoke suite (runs against crew-predictions-staging.web.app)
```

**Staging smoke suite** requires `STAGING_API_KEY` and `SMOKE_TEST_PASSWORD` env vars. In CI these are GitHub secrets. Locally:
```bash
STAGING_API_KEY=... SMOKE_TEST_PASSWORD=... npm run test:smoke
```

**Google sign-in — manual verification required on staging:**
The smoke suite verifies that clicking "Sign in with Google" initiates the redirect (page navigates to Google's OAuth flow). It cannot complete the full Google sign-in automatically because Google's OAuth consent screen blocks automated browsers. After any deploy that touches the auth flow, manually verify on `https://crew-predictions-staging.web.app`:
1. Click "Sign in with Google" on desktop Chrome
2. Complete the OAuth flow
3. Confirm you land on `/matches` and your email appears in the header
4. Repeat on mobile (iOS Safari and Android Chrome)

---

## Status

See [BACKLOG.md](BACKLOG.md) for what's done and what's next.

See [ARCHITECTURE.md](ARCHITECTURE.md) for how the pieces fit together.

---

## What This Is Not

- Official. The podcasts don't know this exists.
- Finished. It's a side project built for fun and to learn Go + Vue + GCP.

