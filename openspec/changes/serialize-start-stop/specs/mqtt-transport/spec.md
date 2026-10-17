## ADDED Requirements

### Requirement: Start and Stop are serialized

Concurrent `Start` and `Stop` calls on one transport SHALL take effect one after another, with at most one receive loop running. A `Stop` SHALL make a `Start` that is still connecting return at once.

#### Scenario: Concurrent Start on a fresh transport

- **WHEN** several goroutines call `Start` at once on a transport that was never started
- **THEN** every call returns without error
- **AND** a later `Stop` returns

#### Scenario: Start while Stop waits for the reader

- **WHEN** `Start` is called while `Stop` waits for a reader that is still handling a message
- **THEN** `Stop` returns once the handler finishes
- **AND** a later `Stop` returns

#### Scenario: Stop while Start waits for the old reader

- **WHEN** `Stop` is called while a `Start` waits for the reader of a stopped session and the broker is down
- **THEN** that `Start` returns a connection-lost error and `Stop` returns once the handler finishes

#### Scenario: Concurrent Start while reconnecting

- **WHEN** several goroutines call `Start` on a started transport that lost its connection
- **THEN** they wait for the connection in parallel, each for at most its own timeout

## MODIFIED Requirements

### Requirement: Stop ends delivery

Stop SHALL disconnect immediately and return only after the MQTT client has left the connected state and the
receive loop has exited, so no handler runs after Stop returns. A stopped transport SHALL be startable again.

#### Scenario: Restart

- **WHEN** a transport is stopped and started again
- **THEN** it connects and delivers messages on its remembered subscriptions

#### Scenario: Start right after Stop

- **WHEN** `Start` is called right after `Stop` returned, with a start retry count of 1
- **THEN** it returns without error
