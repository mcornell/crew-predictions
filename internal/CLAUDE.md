# internal/CLAUDE.md

Go backend guidance for packages in this directory.

---

## Firestore Integration Tests

Every new method on `FirestoreUserStore` (or any Firestore-backed store) requires an integration test that round-trips through the emulator **before** writing the implementation. The e2e suite runs against `MemoryUserStore` — it will never catch a Firestore read/write bug.

- Integration tests live in `*_test.go` files with `//go:build integration`
- Run with: `FIRESTORE_EMULATOR_HOST=localhost:8081 go test -tags integration ./internal/repository/...`
- A test must exercise the full write → read cycle (e.g. `UpdateScores` then `GetByID`/`GetAll`) to catch missing fields in `toUser()`

This was learned the hard way: `UpdateScores` wrote scoring fields to Firestore correctly, but `toUser()` never read them back. Every unit and e2e test passed because they used `MemoryUserStore`. The bug only surfaced on the live staging environment.

---

## Firestore Security Rules

This app has **no client-side Firestore access**. All reads/writes go through the Go service via the Admin SDK, which bypasses security rules. The frontend never imports the Firestore web SDK — its only contract with the backend is the `/api/*` and `/auth/*` HTTP endpoints. The project ID is publicly visible (served at `/auth/config.js` for Firebase Auth), so any rule that permits client access is an external surface anyone with the project ID can hit.

**Posture: deny-all.** The deployed rules must be:

```
rules_version = '2';
service cloud.firestore {
  match /databases/{database}/documents {
    match /{document=**} {
      allow read, write: if false;
    }
  }
}
```

Audited 2026-05-07 via `/firestore-security-rules-auditor` — score 5/5. The only finding was operational: rules are not yet checked into the repo, so they can drift between prod and staging without a PR trail. Tracked in BACKLOG.

**When you change Firestore rules:**
1. Edit `firestore.rules` in the repo (once it lands — see BACKLOG).
2. Run `/firestore-security-rules-auditor` against the diff. Address findings before deploy.
3. Verify deployed rules with `firebase deploy --only firestore:rules` and read back via Firebase MCP `firebase_get_security_rules` to confirm the deploy took.
4. Prod and staging rules must be identical unless there is a documented reason in the PR description.

**If a frontend Firestore SDK call is ever introduced** (e.g., real-time listeners), the deny-all rule will silently block it. Either re-architect the rules to support that specific access pattern (with a fresh audit), or refactor back to an HTTP endpoint that the Go server proxies. Do not relax deny-all without an audit.
