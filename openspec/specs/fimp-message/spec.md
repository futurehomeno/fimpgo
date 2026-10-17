# fimp-message Specification

## Purpose
How a FIMP message is built, encoded to JSON and decoded from the wire, so every app agrees on the
envelope and the typed value it carries.

## Requirements
### Requirement: A new message carries a fresh identity

A newly built message SHALL carry version `1`, a creation time in RFC 3339 with milliseconds, and a
new random UID. When it answers a request, its correlation ID SHALL be the request's UID.

#### Scenario: Plain message

- **WHEN** an app builds a message without a request
- **THEN** `ver` is `1`, `ctime` is the current time, `uid` is a new UUID and `corid` is empty

#### Scenario: Response to a request

- **WHEN** an app builds a message for a request whose `uid` is `abc`
- **THEN** its `corid` is `abc`

### Requirement: The JSON envelope uses the FIMP field names

A message SHALL encode to JSON with the fields `type`, `serv`, `val_t`, `val`, `tags`, `props`,
`ver`, `corid`, `ctime` and `uid`, and SHALL add `resp_to`, `src`, `topic` and `storage` only when
they are set. A binary value SHALL be carried base64-encoded with value type `bin`. An object value
received from the wire and not unpacked SHALL be re-encoded unchanged.

#### Scenario: Optional fields left empty

- **WHEN** a message without a response topic, source, topic or storage strategy is encoded
- **THEN** the JSON has no `resp_to`, `src`, `topic` or `storage` key

#### Scenario: Object passed through

- **WHEN** a received `object` message is published again without reading its value
- **THEN** the `val` is the original JSON object, byte for byte

### Requirement: Decoding follows the value type

Decoding SHALL read `val` as the type named by `val_t`: `string`, `int`, `float`, `bool`, their
`_array` and `_map` forms, `bin`, `object` (kept raw until the consumer unpacks it) and `null`. It
SHALL fail when `type` is missing, when `val_t` is missing or names none of those types, or when
`val` cannot be read as `val_t`; `object` takes any JSON value and `null` ignores `val`. An int, float or bool array element of the wrong type SHALL be skipped
with a warning, while an int, float or bool map entry of the wrong type SHALL fail decoding. A missing `serv` SHALL be logged as a
warning. Missing optional fields (`uid`, `corid`, `ctime`, `resp_to`, `src`, `topic`, `ver`) SHALL
be silent; only a malformed one is logged.

#### Scenario: Int map

- **WHEN** a message with `val_t` `int_map` and `val` `{"a":1,"b":2}` is decoded
- **THEN** its value is the map `a→1, b→2`

#### Scenario: Unknown value type

- **WHEN** a message with `val_t` `decimal`, or without `val_t`, is decoded
- **THEN** decoding fails

#### Scenario: Array with a stray element

- **WHEN** a message with `val_t` `int_array` and `val` `["a", 1]` is decoded
- **THEN** its value is `[1]`

#### Scenario: Value does not match its type

- **WHEN** a message with `val_t` `int` and `val` `"ten"` is decoded
- **THEN** decoding fails with an error naming the value type

#### Scenario: Message without a correlation ID

- **WHEN** a message without `corid` is decoded
- **THEN** decoding succeeds and nothing is logged for `corid`

### Requirement: Typed value access reports a mismatch

Reading the value as a given type SHALL return an error when the value holds another type, and
reading an integer SHALL fail rather than overflow. Properties SHALL be readable as string, int,
float, bool or timestamp, reporting separately whether the key exists and whether it parsed.

#### Scenario: Read a string as int

- **WHEN** a consumer reads an int from a message whose value is `"on"`
- **THEN** it gets an error

#### Scenario: Property absent

- **WHEN** a consumer reads int property `duration` from a message without it
- **THEN** it learns the property is absent, without an error

### Requirement: Timestamps accept common RFC 3339 variants

Parsing a creation time or a timestamp property SHALL accept RFC 3339 with or without fractional
seconds and a zone offset with or without a colon, and the form with a space instead of `T` and a
space before the zone offset. An unparsable timestamp SHALL yield the zero time.

#### Scenario: Offset without colon

- **WHEN** `2024-01-02T03:04:05.123+0100` is parsed
- **THEN** the result is 03:04:05.123 at UTC+1

#### Scenario: Garbage

- **WHEN** `yesterday` is parsed
- **THEN** the result is the zero time

### Requirement: Compressed payloads round-trip

A message SHALL compress to gzip and decompress back to an equal message.

#### Scenario: Round trip

- **WHEN** a message is compressed and then decompressed
- **THEN** the result equals the original message
