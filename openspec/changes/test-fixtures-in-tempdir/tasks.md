## 1. Tests leave the tree clean

- [x] 1.1 Failing check: `go test ./edgeapp ./security` leaves `edgeapp/config_test.json` and `security/test-key-store.json` modified
- [x] 1.2 Save both test files under `t.TempDir()` and remove the tracked fixtures
- [x] 1.3 Rerun: tests pass and `git status --porcelain` is empty
