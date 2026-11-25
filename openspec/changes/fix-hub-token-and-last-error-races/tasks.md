## 1. Hub token

- [x] 1.1 Failing test: `SetHubToken` beside `ExchangeRefreshToken` is reported by the race detector
- [x] 1.2 Guard every `hubToken` access with the hub token mutex

## 2. Last error

- [x] 2.1 Failing test: `SetLastError` beside `LastError` is reported by the race detector
- [x] 2.2 Guard `lastError` with `busMux`

## 3. Tests

- [x] 3.1 Remove `TestLifecycle_WaitForStateChangedWhileStarting`
