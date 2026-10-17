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

### Requirement: Stop returns once disconnected

`Stop` SHALL return only after the MQTT client has left the connected state, so a `Start` right after it connects again.

#### Scenario: Start right after Stop

- **WHEN** `Start` is called right after `Stop` returned, with a start retry count of 1
- **THEN** it returns without error
