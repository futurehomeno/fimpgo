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

#### Scenario: Pending attempt fails during the backoff wait

- **WHEN** the pending connect attempt finishes with an error while `Start` waits between attempts
- **THEN** `Start` starts a new connect attempt for that round instead of counting the finished one again
- **AND** the error of the failed pending attempt is logged

#### Scenario: Broker unreachable for every attempt

- **WHEN** all connect attempts time out
- **THEN** `Start` returns a timeout error
- **AND** the pending attempt is cancelled without waiting for it to wind down

### Requirement: Unparsable payloads are diagnosable

When an incoming payload cannot be decoded as a FIMP message, the transport SHALL log the topic, the payload length and the payload's top-level JSON keys, without logging any values. Keys SHALL be read from the decompressed payload, logged quoted, and bounded to 16 keys of at most 64 characters each.

#### Scenario: Payload without a value type

- **WHEN** a message arrives whose JSON has `serv`, `type` and `val` but no `val_t`
- **THEN** the error log line names the topic, the payload length and the keys `serv`, `type` and `val`

#### Scenario: Compressed payload without a value type

- **WHEN** the same message arrives gzip-compressed on a `j1c1` topic
- **THEN** the error log line names the keys `serv`, `type` and `val`

#### Scenario: Hostile keys

- **WHEN** a payload has a key containing a newline, a key longer than 64 characters, and more than 16 keys
- **THEN** the error log line is a single line, the long key is cut to 64 characters, and the key list ends with `...` after 16 keys
