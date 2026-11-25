## Why

The edgeapp config test and the security key-store test save to tracked files (`edgeapp/config_test.json`, `security/test-key-store.json`), stamping a new timestamp and a freshly generated key on every run. Each `go test ./...` leaves both files modified, and the churn ends up committed (PR 70; PR 72 reverted it, the next test run redid it).

## What Changes

- Both tests save to and load from a file in the test's temporary directory.
- The two tracked fixture files and the now-moot `.gitignore` entry are removed.

## Impact

- `edgeapp/configs_test.go`, `security/key-store_test.go`; tests only, no library code.
