## ADDED Requirements

### Requirement: Start on a started transport keeps the session

`Start` on a transport that is started and not stopped SHALL keep the running session and SHALL return once it is connected, or a timeout error after the timeout.

#### Scenario: Started twice

- **WHEN** `Start` is called again on a connected transport
- **THEN** it returns without error at once
- **AND** a later `Stop` returns
