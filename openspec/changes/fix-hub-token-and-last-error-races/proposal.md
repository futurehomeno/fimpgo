## Why

Review of the 1.18.0 release (PR #69) found two remaining data races in `edgeapp`:

- `FhOAuth2Client` serializes `LoadHubTokenFromCB`, but `SetHubToken` writes and `postMsg` reads `hubToken` without that lock, so a token exchange running beside `SetHubToken` or a token load races on the field. Two exchanges starting without a token also fetch it from CloudBridge twice.
- `Lifecycle.LastError` and `SetLastError` access `lastError` without `busMux`, against the requirement that lifecycle state reads and changes are safe from any goroutine.

`TestLifecycle_WaitForStateChangedWhileStarting` holds `busMux` before the wait subscribes, so it only exercises the already-at-target path; `TestLifecycle_WaitForStateRechecksDroppedEvent` covers the intended case.

## What Changes

- `SetHubToken`, `LoadHubTokenFromCB` and `postMsg` read and write `hubToken` under one mutex; `postMsg` loads the token only if it is still empty once it holds the lock.
- `LastError` and `SetLastError` take `busMux`.
- Remove `TestLifecycle_WaitForStateChangedWhileStarting`.

## Impact

- `edgeapp/auth.go`: `SetHubToken`, `LoadHubTokenFromCB`, `postMsg`.
- `edgeapp/lifecycle.go`: `LastError`, `SetLastError`.
