# Error codes for evt.error.report

## Why

The value of `evt.error.report` is free text today: cliffhanger's router always sends "failed to
process incoming message", zaptec sends its own words, zigbee-ad, ecollector and CloudBridge send
whatever error they had. A client cannot tell a refused value from a cloud outage without parsing
the `msg` property. Cliffhanger change `error-report-policy` (draft PR, under review) proposes that
the value is one of nine error codes. The codes are wire values shared by every FIMP producer and
consumer, so they belong in `fimptype` next to the value types, not in cliffhanger.

## What Changes

- `fimptype.ErrorCodeT`, a string type with `Str()`, modelled on `fimptype.ValueTypeT`, and nine
  constants valued `bad_request`, `out_of_range`, `not_supported`, `not_found`, `busy`,
  `unauthorized`, `device_refused`, `unavailable` and `internal`.
- Additive only. Nothing else moves: no message, transport or edge app code uses the type, and
  cliffhanger keeps its `router.EvtErrorReport` and `Property*` constants.

## Impact

- New file `fimptype/error_code.go`; spec `fimptype`.
- No consumer has to change. Cliffhanger uses the type once this is released; a producer may also
  send a lower_snake code of its own (zaptec: `eco_mode`), as a plain `ErrorCodeT` value.
