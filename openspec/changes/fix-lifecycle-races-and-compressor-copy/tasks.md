## 1. Lifecycle

- [x] 1.1 Failing test: subscribing and unsubscribing while config, auth and connection states change does not crash
- [x] 1.2 State setters and getters hold `busMux`
- [x] 1.3 Failing test: a wait returns when the state changes between its check and its subscription
- [x] 1.4 `WaitForState` subscribes before checking and always unsubscribes
- [x] 1.5 Failing tests: a wait returns when its state event was dropped, and does not remove a subscriber with the same id
- [x] 1.6 `WaitForState` rechecks the state every second and subscribes under its own unique id

## 2. Compressor

- [x] 2.1 Failing test: a compressed payload is unchanged after the next compression
- [x] 2.2 `CompressBinMsg` returns a copy
