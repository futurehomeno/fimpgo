## 1. Baseline

- [x] 1.1 Measure per-file uncovered statements (total 38%)

## 2. Tests

- [x] 2.1 `fimptype/primefimp`: api-client cache and parsing (MQTT round trips not covered), state, site, notify, response
- [x] 2.2 `edgeapp`: auth, manifest, configs, lifecycle, utils, formatters
- [x] 2.3 root `fimpgo`: mqtt_transport, message, mqtt_conn_pool, sync_client gaps
- [x] 2.4 `fimptype`, `transport/streams`, `security`, `utils`

## 3. Verify

- [x] 3.1 `go test -p 1 -cover ./...` passes, total ≥ 60%, tree clean
- [x] 3.2 `go vet ./...`, `golangci-lint run`
