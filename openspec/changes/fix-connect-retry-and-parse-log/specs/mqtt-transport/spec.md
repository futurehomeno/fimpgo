## ADDED Requirements

### Requirement: Start waits on a pending connect attempt

While a connect attempt is still in progress, the transport SHALL keep waiting on it rather than starting another one; with paho's connect retry on, that single attempt does the retrying. It SHALL start a new attempt, after a growing backoff, only once the previous attempt has finished with an error. `Start` SHALL give up connecting after an overall budget of the start retry count times the timeout plus the backoffs between those attempts; once connected, it SHALL wait up to one more timeout for the subscriptions to be restored. When `Start` fails, it SHALL stop the transport: the pending attempt is cancelled and the transport reports disconnected even if that attempt completes afterwards.

#### Scenario: Broker unreachable for longer than one timeout

- **WHEN** the broker cannot be reached and the client has connect retry on
- **THEN** `Start` waits on a single connection attempt
- **AND** no `status can only transition to connecting from disconnected` warning is logged

#### Scenario: Pending attempt connects late

- **WHEN** the pending connect attempt succeeds after the timeout but within the budget
- **THEN** `Start` returns at once without error

#### Scenario: Attempt fails without connect retry

- **WHEN** a connect attempt of a client without connect retry finishes with an error
- **THEN** `Start` starts a new connect attempt after the backoff
- **AND** the error of the failed attempt is logged

#### Scenario: Broker unreachable for the whole budget

- **WHEN** the budget runs out without a connection
- **THEN** `Start` returns a timeout error
- **AND** the pending attempt is cancelled without waiting for it to wind down

#### Scenario: Budget runs out after a failed attempt

- **WHEN** the budget runs out during the backoff after an attempt failed, e.g. with `not Authorized`
- **THEN** `Start` returns that attempt's error rather than a timeout error

#### Scenario: Subscriptions not restored in time

- **WHEN** the connection comes up within the budget but the broker does not confirm the subscriptions within one more timeout
- **THEN** `Start` returns a timeout error and the transport reports disconnected

#### Scenario: Owner of a failed Start connects again

- **WHEN** `SyncClient.Connect` or the edge app's FIMP sync client setup fails in `Start`
- **THEN** it drops the stopped transport, so the next call creates and starts a new one rather than using the stopped one

#### Scenario: Attempt completes after Start failed

- **WHEN** the connection comes up after `Start` has returned an error
- **THEN** the transport reports disconnected

#### Scenario: Attempt of a failed Start completes during the next Start

- **WHEN** the connection of a failed `Start` comes up after the next `Start` began
- **THEN** the transport does not report the next `Start` connected because of it

#### Scenario: Start again while the previous attempt winds down

- **WHEN** `Start` is called again while paho is still winding down the attempt cancelled by a failed `Start`
- **THEN** `Start` retries after the backoff and connects once the broker is reachable

### Requirement: Unparsable payloads are diagnosable

When an incoming payload cannot be decoded as a FIMP message, the transport SHALL log the topic, the payload length and the payload's top-level JSON keys, without logging any values, also not inside the error text. Keys SHALL be read from the decompressed payload, logged quoted, and bounded to 16 keys of at most 64 characters each.

#### Scenario: Payload without a value type

- **WHEN** a message arrives whose JSON has `serv`, `type` and `val` but no `val_t`
- **THEN** the error log line names the topic, the payload length and the keys `serv`, `type` and `val`

#### Scenario: Compressed payload without a value type

- **WHEN** the same message arrives gzip-compressed on a `j1c1` topic
- **THEN** the error log line names the keys `serv`, `type` and `val`

#### Scenario: Hostile keys

- **WHEN** a payload has a key containing a newline, a key longer than 64 characters, and more than 16 keys
- **THEN** the error log line is a single line, the long key is cut to 64 characters, and the key list ends with `...` after 16 keys

#### Scenario: Value of the wrong type

- **WHEN** a message has `"val_t":"float"` and `"val":"user@example.com"`
- **THEN** the error log line names the error kind `Value is not a number`
- **AND** it does not contain `user@example.com`

### Requirement: Message field parse errors omit values

Parse errors of message fields SHALL name the field and the error kind but never the field's value, both in log lines and in the error `NewMessageFromBytes` returns. A wrong-typed optional field (`uid`, `corid`, `ctime`, `resp_to`, `src`, `topic`, `ver`) SHALL be logged at Debug; a missing or null optional field SHALL not be logged.

#### Scenario: Legacy producer sends a numeric version

- **WHEN** a message arrives with `"ver":4242`
- **THEN** one Debug line names the `ver` field
- **AND** no log line contains `4242`

#### Scenario: Required field of the wrong type

- **WHEN** a message arrives with `"type":4242`
- **THEN** `NewMessageFromBytes` returns an error that does not contain `4242`
