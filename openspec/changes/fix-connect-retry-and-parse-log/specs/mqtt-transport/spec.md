## ADDED Requirements

### Requirement: Start waits on a pending connect attempt

While a connect attempt is still in progress, the transport SHALL keep waiting on it rather than starting another one. It SHALL start a new attempt only after the previous attempt has finished with an error. When every attempt has been used up without a connection, `Start` SHALL return a timeout error and abandon the pending attempt.

#### Scenario: Broker unreachable for longer than the first attempt

- **WHEN** the first connect attempt times out while the client is still connecting
- **THEN** later attempts wait on that same connection attempt
- **AND** no `status can only transition to connecting from disconnected` warning is logged

#### Scenario: Pending attempt connects during the backoff wait

- **WHEN** the pending connect attempt succeeds while `Start` waits between attempts
- **THEN** `Start` returns without waiting for the rest of the backoff

#### Scenario: Broker unreachable for every attempt

- **WHEN** all connect attempts time out
- **THEN** `Start` returns a timeout error
- **AND** the client is no longer connecting, so a later `Start` begins a fresh attempt

### Requirement: Unparsable payloads are diagnosable

When an incoming payload cannot be decoded as a FIMP message, the transport SHALL log the topic, the payload length and the payload's top-level JSON keys, without logging any values.

#### Scenario: Payload without a value type

- **WHEN** a message arrives whose JSON has `serv`, `type` and `val` but no `val_t`
- **THEN** the error log line names the topic, the payload length and the keys `serv`, `type` and `val`
