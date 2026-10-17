## MODIFIED Requirements

### Requirement: Partner tokens are exchanged through the Futurehome proxy

The OAuth client SHALL get the hub token from Cloud Bridge with `cmd.clbridge.get_auth_token`
(answered by `evt.clbridge.auth_token_report`), retrying 7 times with 30 s timeouts by default, and
SHALL send token exchanges to the partner proxy for the prod or beta environment with that token as
bearer. A request SHALL be retried 5 times, 60 s apart by default, on a transport error or a status of
400 or above, and SHALL fail when every attempt fails.

#### Scenario: Refresh token

- **WHEN** an app on beta exchanges a refresh token
- **THEN** it posts `refreshToken` and `partnerCode` to `https://partners-beta.futurehome.io/api/control/edge/proxy/refresh`

#### Scenario: Authorization code

- **WHEN** an app on beta exchanges an authorization code
- **THEN** it posts `code` and `partnerCode` to `https://partners-beta.futurehome.io/api/control/edge/proxy/auth-code`

#### Scenario: Proxy down

- **WHEN** every attempt returns status 502
- **THEN** the exchange fails after 5 attempts
