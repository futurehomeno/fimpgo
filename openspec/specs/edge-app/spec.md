# edge-app Specification

## Purpose
Building blocks for hub edge apps: lifecycle states, config files, the app manifest, partner OAuth
through the Futurehome proxy, network checks and log setup.

## Requirements
### Requirement: Lifecycle states are broadcast to subscribers

An app SHALL start in app state `STARTING`, config `NOT_CONFIGURED`, auth `NA` and connection `NA`.
Each change of the app, config, auth or connection state SHALL be sent to every subscriber without
blocking: a subscriber whose buffer is full SHALL miss the event. Waiting for a state SHALL return at
once when it already holds, and otherwise when an event reports it.

#### Scenario: Wait for running

- **WHEN** a component waits for app state `RUNNING` and the app later reaches it
- **THEN** the wait returns once the `RUNNING` event is published

#### Scenario: Slow subscriber

- **WHEN** a subscriber's buffer is full when the auth state changes
- **THEN** that subscriber misses the event and the change still takes effect

### Requirement: Configs fall back to the shipped defaults

The app config SHALL live at `<workdir>/data/config.json`; when it is missing it SHALL be copied from
`<workdir>/defaults/config.json`. Saving SHALL stamp `configured_at` and `configured_by` `auto`.
Restoring defaults SHALL replace the data file with the defaults file.

#### Scenario: First start

- **WHEN** an app starts with no `data/config.json`
- **THEN** the defaults are copied there and loaded

### Requirement: Partner tokens are exchanged through the Futurehome proxy

The OAuth client SHALL get the hub token from Cloud Bridge with `cmd.clbridge.get_auth_token`
(answered by `evt.clbridge.auth_token_report`), retrying 7 times with 30 s timeouts by default, and
SHALL send token exchanges to the partner proxy for the prod or beta environment with that token as
bearer. A request SHALL be retried 5 times, 60 s apart by default, on a transport error or a status of
400 or above, and SHALL fail with the last error or status when every attempt fails.

#### Scenario: Refresh token

- **WHEN** an app on beta exchanges a refresh token
- **THEN** it posts `refreshToken` and `partnerCode` to `https://partners-beta.futurehome.io/api/control/edge/proxy/refresh`

#### Scenario: Proxy down

- **WHEN** every attempt returns status 502
- **THEN** the exchange fails with `server status code=502`

### Requirement: Network readiness can be awaited

The network SHALL count as available when a non-loopback, non-bridge, non-`tap0` interface is up,
broadcast-capable and has an address; the internet SHALL count as available when the network is and
`google.com` resolves. Waiting SHALL poll every 5 s and return a timeout error after the given
timeout, or wait indefinitely when it is zero.

#### Scenario: Offline hub

- **WHEN** an app waits 30 s for the internet on a hub without a route
- **THEN** it gets a timeout error

### Requirement: Logs are configured from the app config

Log setup SHALL select a JSON, `budzik` or coloured text format, parse the level and fall back to
debug when it is invalid, and when a file is given rotate it at 5 MB keeping 2 backups.

#### Scenario: Invalid level

- **WHEN** the configured level is `verbose`
- **THEN** logging runs at debug level
