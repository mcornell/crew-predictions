# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

# Development Approach

## BDD Dual-Loop TDD

Every feature increment starts from a failing **Playwright** (browser) scenario and is driven inward through unit-level red-green-refactor cycles.

### Outer loop (Playwright scenario)

1. **Red** — Write one Playwright test describing the next observable user behavior. Run it. Confirm it fails for the expected reason. Do not proceed until the failure matches intent.
2. **Inner loop** — Repeat until the Playwright test can pass:
   - **Red** — Write the smallest Vitest unit test for the next missing piece the scenario needs. One test at a time. Run it. Confirm it fails.
   - **Green** — Write the **minimum** production code to make that unit test pass. No speculative code. No implementing more than the test demands.
   - **Refactor** — Clean up only covered code. All unit tests must stay green.
3. **Green (scenario)** — Re-run the Playwright test. If still failing, identify the missing piece and return to the inner loop.
4. **Refactor (scenario)** — Refactor across modules if needed. All tests must stay green.
5. Repeat from step 1.

### Discipline rules

- **Never skip red.** If you cannot articulate why a test fails, stop and re-read the requirement.
- **One test at a time.** Never write multiple tests before running them.
- **Minimum code.** Only write production code demanded by the current failing test. Stub everything else.
- **Ask before assuming.** If a design decision is unclear, ask the user before writing code.
- **Commit on every green step** (unit or scenario).
- **Run only the relevant test** after each green step; run the full suite before committing.

## Project Goal

Create a predictions ranking system for fans of Columbus Crew. If this is actually popular, we might expand to other teams...but i mean, right now, only Crew, everyone else can pound sand.  This should be fun, and have a sarcastic tone, like #Crew96 fandom.

We have access to Cloudflare Pages and the mcornell.dev domain to host this. If it's super popular, we'll get another domain.  Because I need to stay on free tiers we can also examine GCP / AWS / Azure's free offerings.  I'm not locked down to using typescript. I'm interested in exploring other languages for fun. At one point I considered micronaut, and recently I thought about using Go or Rust for the heck of it. So let's evaluate what the right choice is and document it here.  We will delete all of the stale code and start fresh.

## Sources

Please always provide sources when responding. That way I can look at the source of truth if I want to learn more details.

