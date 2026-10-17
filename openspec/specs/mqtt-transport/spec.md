# mqtt-transport Specification

## Purpose
How an app connects to the MQTT broker, subscribes, publishes FIMP messages and receives them, so
delivery survives reconnects and a slow consumer cannot stall the client.

## Requirements
### Requirement: Start retries the connection with growing backoff

Start SHALL try to connect up to the configured retry count (10 by default), waiting `n²` seconds
before attempt `n`, and SHALL return once the client is connected. When every attempt fails it SHALL
return the last error, and when the transport is stopped meanwhile it SHALL return a connection-lost
error. Once connected the client SHALL reconnect by itself after a lost connection.

#### Scenario: Broker comes up late

- **WHEN** the broker accepts connections only from the second attempt
- **THEN** Start returns without error, at most 4 s of backoff later

#### Scenario: Stopped while connecting

- **WHEN** the transport is stopped while Start waits between attempts
- **THEN** Start returns a connection-lost error

### Requirement: Subscriptions survive a reconnect

Subscribe SHALL apply the global prefix, wait up to 20 s for the broker to confirm, and remember the
topic. On every (re)connect the transport SHALL subscribe again to every remembered topic. A blank
topic SHALL be ignored. Unsubscribe SHALL forget the topic.

#### Scenario: Broker restart

- **WHEN** the broker restarts and the client reconnects
- **THEN** messages on topics subscribed before the restart are delivered again

#### Scenario: Unconfirmed subscription

- **WHEN** the broker does not confirm a subscription within 20 s
- **THEN** Subscribe returns a timeout error and the topic is not remembered

### Requirement: Publishing encodes by payload type

Publishing to an address SHALL encode the message as JSON for `j1` and as gzip-compressed JSON for
`j1c1`, SHALL set the transport's default source when the message has none, and SHALL apply the
global prefix. Publishing to an address with any other payload type SHALL fail, except a synchronous
publish, which SHALL warn and send JSON as `j1`. Publishing to a topic
string SHALL compress when the topic's payload type is `j1c1`. A synchronous publish SHALL wait up to
5 s for the broker and return a timeout error otherwise. Responding to a request SHALL publish to the
request's response topic and fail when it has none.

#### Scenario: Default source

- **WHEN** an app with default source `zigbee` publishes a message without a source
- **THEN** the published message has `src` `zigbee`

#### Scenario: Explicit source kept

- **WHEN** the same app publishes a message with source `vinculum`
- **THEN** the published message keeps `src` `vinculum`

#### Scenario: Request without response topic

- **WHEN** an app responds to a request that has no `resp_to`
- **THEN** responding fails and nothing is published

### Requirement: Incoming messages reach the handler and matching channels

Each received message SHALL have its global prefix removed, its topic parsed and its payload decoded
(decompressed for `j1c1`), then SHALL go to the message handler and to every registered channel whose
filter accepts it. A channel filter SHALL match its topic as an MQTT route and its service and interface
exactly, `*` matching any service or interface; a filter function SHALL take precedence over it,
and a channel without a filter SHALL get every message. A channel not read within the receive
timeout (10 s by default) SHALL miss that message with a warning. A message with an unknown payload
type or an unparsable topic or payload SHALL be dropped with a log line, and a panic in a filter
function SHALL be logged and count as no match.

#### Scenario: Filtered channel

- **WHEN** a channel filters on topic `#`, service `out_bin_switch` and interface `*`
- **THEN** it gets switch events and commands but not meter reports

#### Scenario: Stuck consumer

- **WHEN** a registered channel is not read for longer than the receive timeout
- **THEN** that message is skipped for it, a warning names the channel, and delivery continues

### Requirement: A stuck receive queue is reported

Received messages SHALL be queued (100 by default) so the MQTT client is never blocked. A message
arriving at a full queue SHALL be dropped with an error log. When more than 20 messages in a row are
dropped while connected, the transport SHALL report `main msg queue stuck` to the error handler.

#### Scenario: Consumer deadlocked

- **WHEN** the message handler blocks and 21 more messages arrive after the queue filled
- **THEN** the error handler gets `main msg queue stuck`

### Requirement: Stop ends delivery

Stop SHALL disconnect immediately and return only after the receive loop has exited, so no handler
runs after Stop returns. A stopped transport SHALL be startable again.

#### Scenario: Restart

- **WHEN** a transport is stopped and started again
- **THEN** it connects and delivers messages on its remembered subscriptions

### Requirement: TLS connections verify both sides

When a client key and certificate are configured, the transport SHALL connect over TLS 1.2 or newer,
trusting `root-ca-1.pem` from the certificate directory plus `root-ca-2.pem` and `root-ca-3.pem` when
present, and presenting the client certificate. For AWS IoT it SHALL negotiate the `x-amzn-mqtt-ca`
protocol. A missing root CA, or a key without a certificate, SHALL fail the setup.

#### Scenario: Missing root CA

- **WHEN** the certificate directory has no `root-ca-1.pem`
- **THEN** no transport is created and the error is logged
