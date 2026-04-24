# CLAUDE.md

Guidance for Claude Code when working in this repository.
See also: [ARCHITECTURE.md](ARCHITECTURE.md) for system design, [BACKLOG.md](BACKLOG.md) for what's next.

---

## Absolute Rules

**When the user establishes a new long-term rule, immediately audit the entire CLAUDE.md and all memory files for anything that contradicts or undermines it — explicit conflicts, implied conflicts, and anything that could be read as permission to violate it. Fix everything found before moving on.**

**When the user sends a message mid-task — feedback, a question, a correction — stop all tool use immediately and respond in text only. Do not resume running tools until the user explicitly signals to continue (e.g. "continue", "go ahead", "yes"). The BDD loop steps, TDD cycle, and any other workflow descriptions do not override this: they describe what to do, not permission to ignore the user.**

**Never `git push` without the user explicitly saying to push in that message.** Commit, then stop. Say what was committed. Wait. This rule has been violated repeatedly and trust has been lost because of it. No exceptions — not for "small" changes, not for docs, not for anything.

The failure mode is always the same: after committing, `git push` feels like the natural way to finish the task, so it runs automatically. That feeling is wrong. After `git commit`, the next tool call must not be `git push` unless the user said "push" in the message that triggered this work. Check: did the user say "push" in this turn? If not, stop.

**Never write production code without a failing test first.** Full BDD loop detail in [e2e/CLAUDE.md](e2e/CLAUDE.md).

---

## Project Goal

Predictions ranking system for Columbus Crew fans. Sarcastic #Crew96 tone. Only Crew — everyone else can pound sand.

**Scoring formats** documented in [README.md](README.md), implemented in `internal/scoring/`.

**Copy tone:** "Pick your scores. Be wrong in public. It's tradition."

---

## Test Commands

```bash
go test ./...          # Go unit tests — always run the full suite
npm run typecheck      # TypeScript type check (vue-tsc --noEmit)
npm run test:unit      # Vitest unit tests for Vue components
npm test               # e2e BDD outer loop — runs against local emulator
npm run test:smoke     # post-deploy smoke suite against staging
```

**Before pushing:** `npm run typecheck`, `npm run test:unit`, and `npm test` must all be green locally first.

---

## Subdirectory Guidance

- [`e2e/CLAUDE.md`](e2e/CLAUDE.md) — BDD dual-loop TDD detail, feature file conventions, @reset isolation, environment notes
- [`src/CLAUDE.md`](src/CLAUDE.md) — Vue test patterns, design language
- [`internal/CLAUDE.md`](internal/CLAUDE.md) — Go backend patterns, Firestore integration test requirements
