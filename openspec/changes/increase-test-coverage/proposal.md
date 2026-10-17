## Why

Total statement coverage is 38%. The largest untested code is the Prime FIMP API client and state/site model (`fimptype/primefimp`, ~600 uncovered statements), `edgeapp` auth/manifest/configs, the MQTT transport and message helpers in the root package, `fimptype/service.go`, `transport/streams.go` and `security/tls_config.go`. Regressions in these paths go unnoticed until they reach an adapter.

## What Changes

- Unit tests for the uncovered public behaviour of `fimptype/primefimp`, `edgeapp`, the root `fimpgo` package, `fimptype`, `transport`, `security` and `utils`.
- HTTP and MQTT peers are faked with `httptest` and the local test broker; no new test reaches an external host.

## Impact

- `*_test.go` files only, no library code.
