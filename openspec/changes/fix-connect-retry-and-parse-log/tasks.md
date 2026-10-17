## 1. Connect retry

- [x] 1.1 Failing test: an unreachable broker leaves no "status can only transition" warning, and Start returns a timeout with the client disconnected
- [x] 1.2 Start reuses the pending connect attempt and disconnects it when giving up
- [x] 1.3 Failing test: the pending attempt connecting during the backoff wait makes Start return at once
- [x] 1.4 The backoff wait also watches the pending attempt
- [x] 1.5 Failing test: the pending attempt failing during the backoff wait still leads to a fresh connect attempt
- [x] 1.6 Drop a failed pending attempt in the backoff wait so the round connects again
- [x] 1.7 Failing test: the failed pending attempt's error is logged
- [x] 1.8 Check for a failed pending attempt before reusing it, whichever wait branch fired, and log its error
- [x] 1.9 Failing test: a failed pending attempt does not cut the backoff short
- [x] 1.10 Keep waiting out the backoff after the pending attempt fails
- [x] 1.11 Failing tests: a second Start right after a failed one connects once the broker is back; a connection completing after a failed Start leaves the transport disconnected
- [x] 1.12 Replace the round loop with one overall budget: Connect again only after a finished attempt, and Stop the transport on any Start failure
- [x] 1.13 Ignore OnConnect once the transport is stopped
- [x] 1.14 Failing test: a broker refusing with `not Authorized` makes Start return that error, not a timeout
- [x] 1.15 Return the last attempt's error when the budget runs out during the backoff
- [x] 1.16 SyncClient.Connect and ConfigureFimpSyncClient drop the transport of a failed Start (no test: Start's budget is fixed at minutes there)
- [x] 1.17 The failed attempt's error is still logged and the backoff still elapses before the next Connect
- [x] 1.18 Failing test: a connect callback of a failed Start that completes after the next Start began leaves the transport disconnected
- [x] 1.19 OnConnect takes the done channel of the Start its connection belongs to
- [x] 1.20 SetStartAutoRetryCount documents that the count sizes the Start budget (no test: doc only)
- [x] 1.21 Failing test: two attempts refused with `not Authorized` end Start within count times the timeout plus the one backoff between them
- [x] 1.22 Sum the backoffs between the attempts (1² to (count-1)²) into the budget
- [x] 1.23 Document that Start waits up to one more timeout for the subscriptions after connecting within the budget (no new test: behaviour unchanged, covered by TestMqttTransport_FailedStartStaysDisconnected)
- [x] 1.24 Failing test: Start with a retry count of 100000 connects to a reachable broker instead of timing out on an overflowed budget
- [x] 1.25 Cap the retry count at 1000 when sizing the budget

## 2. Parse error log

- [x] 2.1 Failing test: an unparsable payload logs its length and top-level keys
- [x] 2.2 Add length and keys to the error line
- [x] 2.3 Failing test: compressed payloads log their keys; hostile keys are quoted and bounded
- [x] 2.4 Read keys from the decompressed payload; quote and bound them
- [x] 2.5 Failing test: a wrong-typed value appears neither in the transport error line nor in NewMessageFromBytes errors and logs; wrong-typed optional fields log at Debug
- [x] 2.6 Strip jsonparser's appended value from field parse errors; log optional-field errors at Debug
