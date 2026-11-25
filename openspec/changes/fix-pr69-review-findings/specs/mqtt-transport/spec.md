## ADDED Requirements

### Requirement: Start makes at least one attempt

`MqttTransport.Start` SHALL treat a start retry count below 1 as 1, so it always waits up to its timeout for a connect attempt.

#### Scenario: Retry count 0

- **WHEN** the start retry count is 0
- **AND** the broker is reachable
- **THEN** `Start` connects and returns no error
