# internal/CLAUDE.md

Go backend guidance for packages in this directory.

---

## Firestore Integration Tests

Every new method on `FirestoreUserStore` (or any Firestore-backed store) requires an integration test that round-trips through the emulator **before** writing the implementation. The e2e suite runs against `MemoryUserStore` — it will never catch a Firestore read/write bug.

- Integration tests live in `*_test.go` files with `//go:build integration`
- Run with: `FIRESTORE_EMULATOR_HOST=localhost:8081 go test -tags integration ./internal/repository/...`
- A test must exercise the full write → read cycle (e.g. `UpdateScores` then `GetByID`/`GetAll`) to catch missing fields in `toUser()`

This was learned the hard way: `UpdateScores` wrote scoring fields to Firestore correctly, but `toUser()` never read them back. Every unit and e2e test passed because they used `MemoryUserStore`. The bug only surfaced on the live staging environment.
