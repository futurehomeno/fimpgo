## Why

Review of #76 found three defects in develop:

- `SetConfigState`, `SetAuthState` and `SetConnectionState` range the subscriber map without `busMux` while `Subscribe` and `Unsubscribe` write it, which can crash the app with a concurrent map access.
- `WaitForState` checks the current state before subscribing, so a change landing in between is never seen and the wait blocks forever.
- `MsgCompressor.CompressBinMsg` returns the internal buffer's bytes and resets it, so the next compression overwrites a payload paho may still be sending.

## What Changes

- All state setters and getters hold `busMux`.
- `WaitForState` subscribes before checking the current state and always unsubscribes.
- `CompressBinMsg` returns a copy of the compressed bytes.

## Impact

- `edgeapp/lifecycle.go`: state setters, getters, `GetAllStates`, `Publish`, `WaitForState`.
- `compressor.go`: `CompressBinMsg`.
