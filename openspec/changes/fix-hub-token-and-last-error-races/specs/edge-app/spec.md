## ADDED Requirements

### Requirement: The hub token is safe for concurrent use

`FhOAuth2Client` SHALL read and write the hub token under one lock, so setting it, loading it from CloudBridge and exchanging tokens may run from any goroutine at the same time. A token exchange SHALL load the hub token only when it is still empty once the lock is held.

#### Scenario: Token set during an exchange

- **WHEN** `SetHubToken` runs while `ExchangeRefreshToken` sends its request
- **THEN** both complete without a data race

### Requirement: The last error is safe for concurrent use

`Lifecycle.LastError` and `SetLastError` SHALL be safe to call from any goroutine at the same time.

#### Scenario: Last error set while read

- **WHEN** one goroutine calls `SetLastError` while another calls `LastError`
- **THEN** both complete without a data race
