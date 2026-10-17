## Why

`MqttTransport.Start` checks whether the transport is started and then calls `connState.Init` without holding a lock across both. Two concurrent `Start` calls on a fresh transport can both pass the check, replace each other's channels and launch two receive loops, so a later `Stop` miscounts the reader wait group. A `Start` that runs while `Stop` waits for the reader also replaces the done channel, and the busy reader picks up the new one when it loops, so `Stop` never returns.

`Stop` calls `Disconnect(0)`, and paho then returns before its disconnect goroutine has run. A `Start` right after `Stop` gets "status can only transition to connecting from disconnected" and waits out a 1 s backoff, or fails when the start retry count is 1.

## What Changes

- `Start` and `Stop` are serialized by one mutex. `Stop` closes the done channel before it takes the mutex, so a `Start` still connecting gives up at once instead of holding `Stop` for its whole retry budget.
- `Start` waits for the reader of a stopped session to exit before it starts a new one.
- `Stop` disconnects with a 250 ms quiesce. paho returns as soon as its disconnect has finished, which normally takes a few milliseconds.

## Impact

- `mqtt_transport.go`: `Start`, `Stop`.
- `mqtt_transport_defs.go`: one mutex field.
