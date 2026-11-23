## MODIFIED Requirements

### Requirement: Incoming messages reach the handler and matching channels

Each received message SHALL have its global prefix removed, its topic parsed and its payload decoded
(decompressed for `j1c1`), then SHALL go to the message handler and to every registered channel whose
filter accepts it. A channel filter SHALL match its topic as an MQTT route and its service and interface
exactly, `*` matching any service or interface; a filter function SHALL take precedence over it,
and a channel without a filter SHALL get every message. A channel not read within the receive
timeout (10 s by default) SHALL miss that message with a warning and SHALL then count as stuck: while
stuck it SHALL get a message only when it has room, without any wait, until it takes one again or a channel is registered under its name. All
full channels SHALL share one receive timeout per message, so channels with room get the message at once
and stuck channels do not add their timeouts up. A message with an unknown payload type or an
unparsable topic or payload SHALL be dropped with a log line, and a panic in a filter function SHALL be
logged and count as no match.

#### Scenario: Filtered channel

- **WHEN** a channel filters on topic `#`, service `out_bin_switch` and interface `*`
- **THEN** it gets switch events and commands but not meter reports

#### Scenario: Stuck consumer

- **WHEN** a registered channel is not read for longer than the receive timeout
- **THEN** that message is skipped for it, a warning names the channel, and delivery continues

#### Scenario: Stuck consumer is not waited on again

- **WHEN** a channel missed a message and is still full when the next message arrives
- **THEN** that message is skipped for it at once and the other channels get it without delay

#### Scenario: Stuck consumer reads again

- **WHEN** a stuck channel has room again
- **THEN** it gets the next message, a log line says it is read again, and later messages wait for it up to the receive timeout

#### Scenario: Stuck consumer is replaced

- **WHEN** a channel is registered under the name of a stuck channel
- **THEN** the new channel is waited on up to the receive timeout like any other

#### Scenario: Several stuck consumers

- **WHEN** three registered channels are not read, one has room and one is read shortly after the message arrives
- **THEN** delivery of that message ends within one receive timeout
- **AND** the channel with room and the channel read late both get the message

#### Scenario: Channel closed while waited on

- **WHEN** a consumer closes its full channel while delivery waits on it
- **THEN** the crash is logged, the transport keeps running and channels can still be registered and unregistered
