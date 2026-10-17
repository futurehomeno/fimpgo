# fimptype Specification

## Purpose
The shared FIMP vocabulary (message, resource and value types, service and thing definitions), so
every consumer writes the same wire values.

## Requirements
### Requirement: Wire values for message, resource and value types

The library SHALL define the message types `cmd`, `evt` and `rsp`; the resource types `dev`, `app`,
`ad`, `cloud`, `discovery` and `loc`; and the value types `string`, `int`, `float`, `bool`,
`str_map`, `int_map`, `float_map`, `bool_map`, `str_array`, `int_array`, `float_array`,
`bool_array`, `object`, `base64`, `bin` and `null`. Exported names SHALL only be added; renaming or
removing one SHALL go through a deprecation.

#### Scenario: Adapter resource type

- **WHEN** a consumer addresses an adapter
- **THEN** the topic segment is `rt:ad`

### Requirement: An interface name determines its message type

The message type of an interface SHALL be taken from its prefix, ignoring surrounding spaces, and
SHALL be unknown when the prefix is none of `cmd`, `evt` or `rsp`.

#### Scenario: Command

- **WHEN** the message type of ` cmd.binary.set` is asked
- **THEN** it is `cmd`

#### Scenario: Other prefix

- **WHEN** the message type of `report.binary` is asked
- **THEN** it is unknown

### Requirement: Service definitions expose typed properties

A service definition SHALL return a property as string, float, int or bool, a list of strings,
floats or ints, or a caller's object, converting a JSON-decoded value to the requested type, and SHALL return the zero
value (or report absence) for a missing or unconvertible property. Ensuring interfaces SHALL add only
those not already present.

#### Scenario: Decoded number list

- **WHEN** a service decoded from JSON has `sup_setpoints` `[18, 21]` and its ints are asked
- **THEN** the result is `[18, 21]`

#### Scenario: Interface already present

- **WHEN** an interface the service already lists is ensured again
- **THEN** the interface list is unchanged

### Requirement: Thing inclusion and exclusion reports

The thing inclusion report SHALL carry the thing's address, services, groups, product name, hash and
ID, manufacturer, device ID, hardware and software version, communication technology, power source,
wake-up interval, security level, technology-specific properties and property sets under their FIMP
JSON keys. The exclusion report SHALL carry the address.

#### Scenario: Inclusion report keys

- **WHEN** an inclusion report is encoded
- **THEN** it has the keys `address`, `services`, `product_hash`, `comm_tech`, `power_source` and `wakeup_interval`
