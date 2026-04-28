# Handoff Guide

This document is written for whoever takes over Crew Predictions. You do not need to know how to code. You need to be able to pay a small monthly bill, respond to an occasional alert, and know who to call when something is broken.

---

## What This Is

A predictions game for Columbus Crew fans. Before each match, users pick a scoreline. After the final whistle, the app scores everyone automatically and updates the leaderboard. Match data comes from ESPN for free — no manual entry needed.

The site runs at whatever domain is configured. As of writing, it is deployed at `crew-predictions.web.app` (Firebase Hosting) backed by `crew-predictions-937208344837.us-east5.run.app` (Google Cloud Run).

---

## What It Costs

**Roughly $1–3/month.** Broken down:

| Service | Cost | Notes |
|---|---|---|
| Google Cloud Run | ~$0–1/mo | Scales to zero between requests |
| Firestore | ~$0/mo | Within always-free tier |
| Firebase Hosting | Free | CDN for the frontend |
| Domain (if active) | ~$12/yr | Must renew annually or auto-renew |
| Billing killswitch | ~$0/mo | Cloud Function that shuts off billing if costs spike |

A budget alert is configured — if the bill ever exceeds $10/mo, Google will email the billing account owner. That should never happen under normal use.

---

## How to Check If It's Working

1. Open the site in a browser. If you see match cards, it's working.
2. If you see an error or blank page, try again in 5 minutes — Cloud Run sometimes has a cold start delay.
3. If it's still broken, check the GCP Console: [console.cloud.google.com](https://console.cloud.google.com) → Project: `crew-predictions` → Cloud Run → `crew-predictions` → Logs tab. Look for red error lines.

---

## Who Owns What

| Thing | Where It Lives | Current Owner |
|---|---|---|
| GCP project (billing, Cloud Run, Firestore) | Google Cloud | Mike Cornell (mcornell74@gmail.com) |
| GitHub repository (code) | github.com/mcornell/crew-predictions | Mike Cornell |
| Firebase project (auth, hosting) | Firebase Console | Mike Cornell |
| Domain (if active) | Registrar TBD | Mike Cornell |

All of these need to transfer to you. Steps below.

---

## How to Transfer Ownership

### Step 1 — Get access to GCP

Mike will add you as an Owner on the GCP project before transfer. You'll receive an email from Google Cloud — accept it and sign in with a Google account you control.

After you have Owner access:
- Go to [console.cloud.google.com](https://console.cloud.google.com)
- Select project `crew-predictions`
- Go to **Billing** → link your own credit card or billing account
- Remove Mike's billing account after yours is linked

### Step 2 — Transfer the GitHub repository

Mike will transfer the repository to your GitHub account or a GitHub organization you control:
- GitHub → Settings → Transfer repository
- You accept via the emailed link
- The code, history, and CI/CD pipeline all come with it

You do not need to do anything with the code after transfer. The pipeline runs automatically when code is pushed to `main`.

### Step 3 — Transfer the domain (if active)

Domain registrars allow transfers between accounts. This usually takes 5–7 days. Keep auto-renew enabled. Pre-pay multiple years if possible.

### Step 4 — Update GitHub Actions secrets

The CI/CD pipeline uses a GCP service account to deploy automatically. After the GCP transfer, you'll need to update one GitHub secret:

- Go to the repo → Settings → Secrets and variables → Actions
- The secret named `WORKLOAD_IDENTITY_PROVIDER` and `SERVICE_ACCOUNT` tie to the GCP project
- These may continue working after transfer if the project ID stays the same — test by pushing a small change to `develop` and watching the Actions tab

If CI fails with a permission error, open an issue on the repo and note the error message. Someone in the community with GCP knowledge can help recreate the service account.

---

## What Could Break and How to Fix It

### The site is down — "service unavailable"
Most likely: billing lapsed. Check [console.cloud.google.com/billing](https://console.cloud.google.com/billing). If the project is suspended, re-enable billing. The service restarts automatically within a few minutes.

### No new matches are appearing
The app fetches match data from ESPN's unofficial API automatically. If no Crew matches appear for an upcoming gameweek, try:
1. Wait 24 hours — the app refreshes match data daily at 4am ET
2. Trigger a manual refresh: send a POST request to `/admin/refresh-matches` with the admin key (see below)

If ESPN changes their API format (this has happened before), the match poller will silently stop finding new matches. This requires a code fix — post a GitHub issue describing the problem and the community may be able to help, or hire a developer for an hour.

### Google sign-in is broken
The most common cause: `FIREBASE_AUTH_DOMAIN` environment variable must match the serving domain. If you ever change the domain the site lives at, update this variable:

```
gcloud run services update crew-predictions \
  --region=us-east5 \
  --update-env-vars FIREBASE_AUTH_DOMAIN=<your-domain-here>
```

### The leaderboard stopped updating after a match
Results need to be entered manually after each match. There is no UI for this — it's done via the admin API:

```
POST /admin/results
X-Admin-Key: <admin key>
Body: match_id=<id>&home_score=<n>&away_score=<n>
```

The admin key is stored in GCP Secret Manager under `crew-predictions` → Secret Manager → `ADMIN_KEY`. You can view it in the GCP Console. After entering a result, scores recalculate automatically.

---

## The Admin Key

The admin key protects match result entry and other administrative operations. It lives in:

**GCP Console → Secret Manager → `crew-predictions` project → Secret: `ADMIN_KEY`**

Keep this private. Do not share it publicly. If it is ever exposed, rotate it by creating a new secret version and redeploying Cloud Run with the updated value.

---

## Keeping the Site Running Year to Year

The main ongoing tasks:

1. **Pay the bill** (~$1–3/mo, via GCP billing)
2. **Renew the domain** (annually, via registrar)
3. **Enter match results** after each Crew match (via admin API or if a UI gets built)
4. **At season end** — archive the leaderboard and reset for the new season (feature planned — see `BACKLOG.md`)

The app handles match fetching, scoring, and leaderboard updates automatically. The only human task during the season is entering final scores.

---

## Getting Help

- **Code questions:** Open an issue on the GitHub repository. The commit history and README explain how everything works.
- **GCP / Firebase questions:** [cloud.google.com/support](https://cloud.google.com/support) — free support is limited but the documentation is thorough.
- **Billing emergencies:** Google will not delete your data immediately if billing lapses — there is a grace period of roughly 30 days before anything is permanently lost.

---

## A Note From Mike

This was built for the Crew community I've known for thirty years. I hope it keeps running for a long time and brings some joy to matchdays. The codebase is clean and well-tested — don't let anyone tell you it's too complicated to maintain. It mostly runs itself.

Go Crew.
