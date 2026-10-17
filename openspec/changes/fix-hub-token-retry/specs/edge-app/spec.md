## ADDED Requirements

### Requirement: Token exchange retries use the current hub token

Each attempt of a token exchange SHALL send the hub token current when the attempt starts.

#### Scenario: Hub token replaced between attempts

- **WHEN** the first attempt of `ExchangeRefreshToken` is rejected and `SetHubToken` sets a new token before the next attempt
- **THEN** the next attempt sends the new token
