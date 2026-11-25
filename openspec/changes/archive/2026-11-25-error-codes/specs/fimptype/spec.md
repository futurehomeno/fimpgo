## ADDED Requirements

### Requirement: Error codes

`fimptype` SHALL declare `ErrorCodeT`, a string type with a `Str()` method like `ValueTypeT`, and
one constant for each error code of an `evt.error.report` value, valued with the exact wire name:
`ErrorCodeBadRequest` `bad_request`, `ErrorCodeOutOfRange` `out_of_range`,
`ErrorCodeNotSupported` `not_supported`, `ErrorCodeNotFound` `not_found`, `ErrorCodeBusy` `busy`,
`ErrorCodeUnauthorized` `unauthorized`, `ErrorCodeDeviceRefused` `device_refused`,
`ErrorCodeUnavailable` `unavailable`, `ErrorCodeRateLimited` `rate_limited` (too many requests for
now; temporary, the caller may try again later) and `ErrorCodeInternal` `internal`. Every wire name SHALL be
lower_snake. A producer MAY send a lower_snake code of its own as an `ErrorCodeT`.

#### Scenario: Code on the wire

- **WHEN** a producer sends `fimptype.ErrorCodeOutOfRange` as the value of an `evt.error.report`
- **THEN** the value on the wire is `out_of_range`

#### Scenario: Internal code

- **WHEN** `fimptype.ErrorCodeInternal.Str()` is asked
- **THEN** it is `internal`

#### Scenario: Producer's own code

- **WHEN** a producer sends `fimptype.ErrorCodeT("eco_mode")`
- **THEN** the value on the wire is `eco_mode`
