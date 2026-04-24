# e2e/CLAUDE.md

BDD guidance for the Playwright e2e suite in this directory.

---

## BDD Dual-Loop TDD

Every feature increment starts from a failing Playwright scenario and is driven inward through unit-level red-green-refactor cycles.

### Outer loop

1. **Red** — Write one Gherkin scenario describing the next observable user behavior. Run `npm test`. Confirm it fails for the expected reason. Do not proceed until the failure matches intent.
2. **Inner loop** — Repeat until you believe the scenario can pass:
   - **Red** — Write the smallest failing unit test for the next missing piece. One test at a time. Run it. Confirm it fails.
   - **Green** — Write the minimum production code to pass that one test. Only what the test demands. Run it. **Commit.**
   - **Refactor** — Clean up covered code only. All tests stay green. **Commit.**
3. **Coverage gate** — Run `go test ./... -cover`. Any uncovered branch means production code was written without a test — go back to the inner loop.
4. **Green (scenario)** — Run `npm test`. Still failing? Identify the next missing piece and return to the inner loop. Passes? **Commit.**
5. **Refactor (scenario)** — Refactor across modules if needed. All tests stay green. **Commit.**
6. Repeat from step 1.

### Absolute rules

- **No production code without a red test.** Every `if`, every `return`, every error branch — driven by a failing test first. No exceptions for "wiring", "infrastructure", or "obvious" code.
- **This applies to Vue code too.** Before creating any `.vue` or `.ts` in `src/`, write a failing Vitest test first.
- **One branch, one test, one commit.** Never write a whole function and test it afterward.
- **Always run the full suite.** `go test ./...` and `npm run test:unit` at every red/green/refactor step.
- **Never skip red.** If you cannot articulate exactly why the test fails, stop.
- **Exception — external HTTP calls:** Note the gap as tech debt; do not silently accept low coverage.

---

## Conventions

**Feature files:** `e2e/features/` · **Step definitions:** `e2e/steps/`
**Smoke features:** `e2e/smoke/features/` · **Smoke steps:** `e2e/smoke/steps/`

Always run `npx bddgen` after editing a `.feature` file.

**@reset tag:** A `Before` hook calls `DELETE /admin/reset` before each `@reset` scenario. Any feature that mutates match/prediction/result stores must carry `@reset`. Auth-only features omit it and run in parallel via the `auth` Playwright project.

**Smoke suite:** `npm run test:smoke` runs in CI after `deploy-staging`. It is not a substitute for the local e2e suite — run `npm test` locally first.

---

## Environment

`npm test` builds Vue then starts its own Go server on :8082 (TEST_MODE). Emulators must already be running (`./dev.sh` handles that). `dev.sh` and `npm test` can run simultaneously — no port conflict.

Do **not** set `GOOGLE_CLOUD_PROJECT` in the playwright `webServer` env — it triggers Firestore and breaks test isolation.
