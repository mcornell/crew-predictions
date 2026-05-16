# Match Polling — External Trigger Architecture

Replaces reliance on the in-process goroutine-based match poller with externally-triggered Cloud Tasks chains, plus externally-triggered daily refreshes via Cloud Scheduler. The in-process poller remains as a fast-path optimization when the container happens to be warm.

Status: shipped via PR #52 (2026-05-16). Live in both staging and prod.

---

## Problem

Live match polling and the daily ESPN refresh are both in-process Go goroutines. They die when the Cloud Run container dies — which happens ~15 min after the last HTTP request, because the project runs with `--min-instances=0` and default CPU throttling.

For matches where someone is actively watching the site, the Vue SPA's 30s `/api/matches` polling keeps the container warm and the goroutine alive. For matches with no engagement (vacations, away games, midweek cup ties), the container dies mid-match, polling stops, and the live clock freezes on the last successful value until either a user hits the site or the next match window.

The 4 AM ET daily ESPN refresh has the same problem: nobody is browsing at 4 AM, so the container is dead and the refresh goroutine never runs.

Confirmed observation: NYRB at MSG on 2026-05-13. Container died twice during the match (29-min and 47-min gaps in polling logs); each gap began ~15 min after the last user-driven HTTP request. Same pattern is present in earlier matches but was masked when the user was actively watching (each page load restarted the container).

## Goal

External triggers that wake the container exactly when work is needed:

- Daily ESPN refresh at known times — 4 AM, 12 PM, 6 PM ET
- Match-window polling that runs **only while a match is in flight** and survives container deaths across the 2-min tick interval
- Safe to call when a match is already in progress (no duplicate chains, no `Reset()` races on the in-process poller)
- No 24/7 keepalive, no `--min-instances=1`, no `--no-cpu-throttling`. Stay within the free tier.

## Architecture

```
┌─────────────────────┐   POST /admin/refresh-matches
│  Cloud Scheduler    │ ─────────────────────────────────►  ┌──────────────┐
│  4am / 12pm / 6pm   │                                     │   Cloud Run  │
└─────────────────────┘                                     │   Go server  │
                                                            └──────┬───────┘
                                                                   │
                                          enqueue task at          │
                                          kickoff - 5min           │
                                                                   ▼
┌─────────────────────┐                                ┌──────────────────────┐
│  Cloud Tasks queue  │                                │   Firestore          │
│  `match-polling`    │ ◄──────── enqueue next ──────  │   matches/{id}       │
│                     │                                │     .lastPollAt      │
└──────────┬──────────┘                                │     .chainSeededFor  │
           │                                           │     .abandonedAt     │
           │ POST /admin/poll-scores?matchID=...       └──────────────────────┘
           │ (every 2 min until terminal)
           ▼
   ┌──────────────┐
   │  Cloud Run   │  ── poll ESPN ─►  update Firestore  ─► (if non-terminal)
   │  Go server   │                                          enqueue next task
   └──────────────┘
```

## Components

### Cloud Scheduler (3 jobs, exactly fills the free tier)

| Name | Cron (America/New_York) | Target | Purpose |
|---|---|---|---|
| `refresh-4am` | `0 4 * * *` | POST /admin/refresh-matches | Catch overnight ESPN updates; seed morning + afternoon matches |
| `refresh-noon` | `0 12 * * *` | POST /admin/refresh-matches | Catch midday updates; seed evening matches |
| `refresh-6pm` | `0 18 * * *` | POST /admin/refresh-matches | Catch evening updates; recover dead chains for matches still in progress |

Auth: existing `AdminAuth` middleware (`X-Admin-Key` header).

### Cloud Tasks (1 queue)

- Queue: `match-polling`, region `us-east4` (Cloud Tasks doesn't ship in `us-east5` where Cloud Run lives — `us-east4` is the closest east-coast region for the queue; cross-region HTTP target is fine), both projects
- Free tier: 1M ops/month (expected usage: ~1k ops/month)
- Built-in retries on delivery failure (exponential backoff, default up to ~1h) — handles container cold-start transparently
- Task target: `POST /admin/poll-scores?matchID={id}` with admin-key header

### Firestore — new fields on `matches/{matchID}`

```
matches/{matchID}
  ...existing fields...
  lastPollAt:        timestamp   // updated by every /admin/poll-scores call
  chainSeededFor:    timestamp   // kickoff time the current chain was seeded for; dedup key
  abandonedAt:       timestamp   // set by 4h safety bailout for diagnostic visibility
```

### Go code changes

| File | Change |
|---|---|
| `internal/tasks/client.go` (new) | Cloud Tasks SDK wrapper. `EnqueuePoll(matchID, runAt) error`. Task name encodes `matchID + runAt.Unix()` for deterministic uniqueness within the 1h dedup window. Interface-based so tests can swap a fake. |
| `internal/handlers/refresh_matches.go` | After existing fetch+cache logic, iterate matches and apply the state-based enqueue rules (table below). |
| `internal/handlers/poll_scores.go` | Take `matchID` query param (or default to all-active for backwards compat). Write `lastPollAt = now` on the match. After polling, if non-terminal status, call `EnqueuePoll(matchID, now+2min)`. If terminal or `(now - kickoff) > 4h`, do nothing (chain ends; safety stop also writes `abandonedAt`). |
| `internal/repository/matches.go` | Add `LastPollAt`, `ChainSeededFor`, `AbandonedAt` fields to Match struct + Firestore mapping. |
| `internal/poll/match_poller.go` | `Reset()` becomes a *soft* reset — leaves entries in the `active` map alone if their match is state=in. Only clears+reschedules pre matches. |
| `cmd/server/config.go` | Add `CloudTasksQueue` config field (env var). |

### New IAM

| SA | Grant | Why |
|---|---|---|
| Cloud Run runtime SA (default Compute SA today) | `roles/cloudtasks.enqueuer` on the `match-polling` queue | Server can enqueue follow-up tasks |
| Cloud Scheduler SA (auto-created) | none required | Auth is via admin-key header on the request |

---

## Refresh logic by match state

`/admin/refresh-matches` iterates ESPN matches and applies:

| State | Condition | Action |
|---|---|---|
| **pre** | kickoff in next ~8h AND `chainSeededFor != kickoff` | Enqueue task at `kickoff - 5min`, set `chainSeededFor = kickoff` |
| **pre** | kickoff already seeded (`chainSeededFor == kickoff`) | Nothing — task is already in the queue |
| **pre** | kickoff > 8h out | Nothing — next refresh will pick it up |
| **in** | `lastPollAt > now - 5min` | Nothing — chain is alive, don't disturb it |
| **in** | `lastPollAt` stale or unset | Enqueue immediate revival task |
| **post** | (terminal) | Nothing — chain has ended naturally |

## Chain termination

`/admin/poll-scores` polls ESPN for one matchID, writes `lastPollAt = now`, then:

| ESPN status | Action |
|---|---|
| `STATUS_FULL_TIME`, `STATUS_FINAL_AET`, `STATUS_FINAL_PEN` | Save result, run Recalculate, **end chain** (no enqueue) |
| `STATUS_POSTPONED`, `STATUS_CANCELED`, `STATUS_ABANDONED` | **End chain.** Next refresh picks up any rescheduled kickoff and seeds a fresh chain. |
| `STATUS_IN_PROGRESS`, `STATUS_HALFTIME`, `STATUS_DELAYED`, `STATUS_END_PERIOD`, ... | **Enqueue next task** at `now + 2min` |
| any non-terminal AND `now - kickoff > 4h` | **Safety bailout** — write `abandonedAt = now`, **end chain.** Longest legitimate soccer match (90 + ET 30 + PKs + buffer) is well under 4h. Past that, assume something's wrong with ESPN data; next refresh will re-evaluate. |

---

## Edge cases addressed

| Scenario | Behavior |
|---|---|
| **Sunday 4:30 PM kickoff, 6 PM refresh during match** | Refresh sees state=in, `lastPollAt` is fresh (2 min old), takes no action. Chain ticks undisturbed. |
| **Wednesday 7:30 PM kickoff, vacation (no traffic)** | 6 PM refresh seeds task at 7:25 PM. Task fires → wakes container → polls. Every 2 min thereafter, next task wakes container as needed. Match ends → chain ends. |
| **Container dies mid-chain between ticks** | Next scheduled task fires, hits Cloud Run, wakes new container, polls, enqueues next task. Cloud Tasks retries handle any wake delay. |
| **Match delayed 30 min for weather** | `STATUS_DELAYED` is non-terminal → chain keeps ticking through delay. Resumes when match resumes. |
| **Match postponed entirely** | `STATUS_POSTPONED` ends chain. Next 4am/12pm/6pm refresh sees the new kickoff date on ESPN and seeds a fresh chain for that day. |
| **Chain task delivery fails (Cloud Run 503)** | Cloud Tasks retries with exponential backoff (~30s → ~1m → ~5m → up to ~1h). Container eventually wakes. |
| **Chain enqueue silently dies** (e.g., poll succeeds but Cloud Tasks call crashes before next enqueue) | `lastPollAt` goes stale. Next 4am/12pm/6pm refresh sees state=in + stale `lastPollAt` → enqueues revival task. Worst case: up to 6 hours of dead chain. |
| **Same match seeded twice by overlapping refreshes** | `chainSeededFor` field is the dedup key. Second refresh sees field already matches current kickoff, skips. |
| **Same task enqueued twice** | Cloud Tasks deduplication via task name (`matchID + unix-ts`) within the 1h dedup window. |
| **Match goes to extra time + PKs** | All non-terminal until `STATUS_FINAL_PEN`. Chain ticks through. 4h safety stop protects against runaway. |

---

## Cost

- Cloud Tasks: ~1k ops/month vs 1M free tier
- Cloud Scheduler: 3 jobs vs 3 free tier (exact fit)
- Cloud Run: ~30 wakeups per match × 10 matches × ~200ms each ≈ 60 vCPU-sec/month vs 180k free tier
- **Net incremental cost: $0**

---

## Rollout sequence

1. Code changes locally + staging — add Cloud Tasks client, modify refresh + poll handlers, add Firestore fields. Unit + integration tests.
2. Enable APIs on both projects — Cloud Tasks, Cloud Scheduler. Create `match-polling` queue in both.
3. IAM grants — `roles/cloudtasks.enqueuer` to the Cloud Run runtime SA on both projects.
4. Deploy to staging — Cloud Run picks up new code. Manually hit `/admin/refresh-matches` against staging. Watch staging logs to confirm task enqueue → fire → poll → next enqueue.
5. Add Cloud Scheduler jobs in staging — 4am/12pm/6pm ET. Verify they fire and refresh successfully.
6. Merge to main, deploy to prod.
7. Add Cloud Scheduler jobs in prod — same three crons.
8. Watch one full match cycle to confirm end-to-end behavior in prod.
9. Update ARCHITECTURE.md, infra/wif-setup.md, BACKLOG.md, internal/CLAUDE.md.

---

## Design decisions (locked in)

| Question | Decision |
|---|---|
| Should the in-process `MatchPoller` goroutine stay at all? | **Yes, as a fast-path.** Costs nothing when the container is already warm, gives sub-2-min latency for users actively watching. The Cloud Tasks chain is the durability layer; the goroutine is the optimization layer. |
| Should we delete `match-polling` tasks when a chain ends? | **No.** They've already executed or been pruned by then. Cloud Tasks handles cleanup. |
| Should the 4h safety stop write a Firestore flag? | **Yes** — `abandonedAt`. Cheap and helpful for debugging. |
| Should refresh skip the in-process `poller.Reset()` entirely? | **No** — keep it for cold-start matches that haven't been chain-seeded yet. Just make it soft (skip in-progress matches). |

---

## Effort estimate

~2-3 hours of code + ~30 min infra setup. Most of the complexity is the first-time Cloud Tasks integration.
