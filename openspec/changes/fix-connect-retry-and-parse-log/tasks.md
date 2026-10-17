## 1. Connect retry

- [x] 1.1 Failing test: an unreachable broker leaves no "status can only transition" warning, and Start returns a timeout with the client disconnected
- [x] 1.2 Start reuses the pending connect attempt and disconnects it when giving up
- [x] 1.3 Failing test: the pending attempt connecting during the backoff wait makes Start return at once
- [x] 1.4 The backoff wait also watches the pending attempt
- [x] 1.5 Failing test: the pending attempt failing during the backoff wait still leads to a fresh connect attempt
- [x] 1.6 Drop a failed pending attempt in the backoff wait so the round connects again
- [x] 1.7 Failing tests: a second Start right after a failed one connects once the broker is back; a connection completing after a failed Start leaves the transport disconnected
- [x] 1.8 Replace the round loop with one overall budget: Connect again only after a finished attempt, and Stop the transport on any Start failure
- [x] 1.9 Ignore OnConnect once the transport is stopped

## 2. Parse error log

- [x] 2.1 Failing test: an unparsable payload logs its length and top-level keys
- [x] 2.2 Add length and keys to the error line
- [x] 2.3 Failing test: compressed payloads log their keys; hostile keys are quoted and bounded
- [x] 2.4 Read keys from the decompressed payload; quote and bound them
- [x] 2.5 Failing test: a wrong-typed value appears neither in the transport error line nor in NewMessageFromBytes errors and logs; wrong-typed optional fields log at Debug
- [x] 2.6 Strip jsonparser's appended value from field parse errors; log optional-field errors at Debug
