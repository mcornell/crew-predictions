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

This rule has been violated repeatedly. The specific failure pattern: new UI behavior (e.g. "Now Playing" section) gets implemented by changing frontend filtering logic without first writing a failing e2e scenario. The result is the feature ships with zero e2e coverage AND the filtering change breaks 18 existing tests. Both outcomes happened in the same session. The user was angry.

A second recurring failure pattern: infrastructure and "helper" code — repository implementations (`MemoryConfigStore`, `MemoryUserStore.Reset`), seed handlers (`SeedUserHandler`), Vue component additions (AppHeader flyout/drawer) — gets written without unit tests first because it feels like plumbing rather than logic. It is not exempt. Every function, method, branch, and error path requires a failing test before the line of production code that satisfies it. If `go test ./... -cover` shows any function below 80% coverage that isn't a Firestore adapter (covered by integration tests) or external SDK wrapper, that is a violation.

**The outer BDD loop is not optional.** Before touching any production file for a new feature:
1. Write the Gherkin scenario. Immediately create stub step definitions so the suite can run (missing steps block all other tests). Run `npm test`. Confirm it fails for the right reason — the new scenario fails, all others pass.
2. Only then open any `.vue`, `.go`, or `.ts` production file.

If you find yourself writing production code and cannot point to the specific failing e2e scenario that demands it, stop. You are violating this rule.

**After any feature or fix, run `go test ./... -cover` and `npm run test:unit -- --coverage`. Check every function for coverage.** Uncovered functions and branches that aren't Firestore adapters or external SDK wrappers must be covered before declaring the work done. Do not wait for the user to ask.

**Run all three test suites at every commit, not just at the end.** `go test ./...`, `npm run typecheck && npm run test:unit`, and `npm test` must all be green before every commit during a feature build. "I only changed Go code" is not an excuse to skip `npm test` — a feature file with missing step definitions blocks the entire e2e suite and hides regressions in unrelated tests.

**Before declaring any work done — feature OR refactor — run `go test ./...`, `npm run typecheck && npm run test:unit`, and `npm test`. All three must be green.** This applies to refactors and removals too, not just new features — observable behavior can regress without a new scenario failing first.

---

## Project Goal

Predictions ranking system for Columbus Crew fans. Sarcastic #Crew96 tone. Only Crew — everyone else can pound sand.

**Scoring formats** documented in [README.md](README.md), implemented in `internal/scoring/`.

**Copy tone:** "Pick your scores. Be wrong in public. It's tradition."

---

## Test Commands

```bash
go test ./...                                           # Go unit tests — always run the full suite
FIRESTORE_EMULATOR_HOST=localhost:8081 go test -tags integration ./internal/repository/...  # Firestore integration tests (requires emulator)
npm run typecheck      # TypeScript type check (vue-tsc --noEmit)
npm run test:unit      # Vitest unit tests for Vue components
npm test               # e2e BDD outer loop — runs against local emulator
npm run test:smoke     # post-deploy smoke suite against staging
```

**Before pushing:** `go test ./...`, `npm run typecheck`, `npm run test:unit`, and `npm test` must all be green locally first. When emulators are running, also run the integration tests — they use `//go:build integration` so `go test ./...` silently skips them.

---

## Subdirectory Guidance

- [`e2e/CLAUDE.md`](e2e/CLAUDE.md) — BDD dual-loop TDD detail, feature file conventions, @reset isolation, environment notes
- [`src/CLAUDE.md`](src/CLAUDE.md) — Vue test patterns, design language
- [`internal/CLAUDE.md`](internal/CLAUDE.md) — Go backend patterns, Firestore integration test requirements
