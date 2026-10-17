# topic-address Specification

## Purpose
How FIMP MQTT topics are built from and parsed into addresses, including the global prefix that
routes a topic to one hub through the cloud.

## Requirements
### Requirement: Topics follow the FIMP layout per resource type

An address SHALL serialize to `pt:<payload>/mt:<msg type>/rt:<resource type>/…` with the payload type
defaulting to `j1`. Adapter, app and cloud resources SHALL add `rn:<name>/ad:<address>`; devices SHALL
add `rn:<name>/ad:<address>/sv:<service>/ad:<service address>`; discovery SHALL stop after `rt:`. A
`+` or `#` component SHALL be written bare, without its key, so the address can be used as a
subscription.

#### Scenario: Device service topic

- **WHEN** an event address for the `out_bin_switch` service `1_1` of device `zigbee` address `1` is serialized
- **THEN** the topic is `pt:j1/mt:evt/rt:dev/rn:zigbee/ad:1/sv:out_bin_switch/ad:1_1`

#### Scenario: Wildcard

- **WHEN** a device address with service address `+` is serialized
- **THEN** the last segment is `+`

#### Scenario: Discovery

- **WHEN** a discovery command address is serialized
- **THEN** the topic is `pt:j1/mt:cmd/rt:discovery`

### Requirement: Topics parse back into addresses

Parsing a topic SHALL fill each field from its `key:value` segment, assigning the first `ad:` to the
resource and an `ad:` after `sv:` to the service. A leading segment without a key SHALL be taken as
the global prefix. Any other segment without exactly one `key:value` pair SHALL fail the parse.

#### Scenario: Round trip

- **WHEN** a serialized device topic is parsed
- **THEN** every field equals the original address

#### Scenario: Prefixed topic

- **WHEN** `site-1/pt:j1/mt:cmd/rt:app/rn:vinculum/ad:1` is parsed
- **THEN** the global prefix is `site-1` and the resource name is `vinculum`

#### Scenario: Malformed segment

- **WHEN** `pt:j1/mt:evt/rt:dev/rn:zigbee/foo` is parsed
- **THEN** parsing fails

### Requirement: The global prefix is added and stripped transparently

Adding a prefix SHALL join it with `/` unless the topic already starts with `/`, and SHALL leave the
topic unchanged when the prefix is blank. Detaching SHALL split a topic at its `pt:` segment into
prefix and FIMP topic.

#### Scenario: Add prefix

- **WHEN** prefix `site-1` is added to `pt:j1/mt:evt/rt:app/rn:x/ad:1`
- **THEN** the topic is `site-1/pt:j1/mt:evt/rt:app/rn:x/ad:1`

#### Scenario: Detach multi-level prefix

- **WHEN** `a/b/pt:j1/mt:evt/rt:app/rn:x/ad:1` is detached
- **THEN** the prefix is `a/b` and the topic is `pt:j1/mt:evt/rt:app/rn:x/ad:1`

### Requirement: Route matching follows MQTT wildcards

A route SHALL match a topic segment by segment, `+` matching exactly one segment and `#` matching the
rest of the topic, including nothing.

#### Scenario: Single-level wildcard

- **WHEN** route `pt:j1/mt:evt/rt:dev/rn:zigbee/ad:1/sv:+/ad:+` is matched against a zigbee device event
- **THEN** it matches

#### Scenario: Different length

- **WHEN** route `pt:j1/mt:evt` is matched against `pt:j1/mt:evt/rt:app`
- **THEN** it does not match
