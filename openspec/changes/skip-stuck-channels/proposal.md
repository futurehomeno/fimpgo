## Why

After `share-channel-receive-timeout`, a channel that is never read again still costs one receive timeout (10 s) on every message it is interested in. During that wait the receiving goroutine holds the channel registration lock: other channels get the next message late, registering and unregistering block, and at more than one message per timeout the receive queue fills until the transport stops (`main msg queue stuck`). Every message also logs the same warning (2,325 hits on one tpflow site). Issue 1.

## What Changes

- A channel that misses a message after the receive timeout is marked stuck and logged once.
- A stuck channel gets a message only when it has room; otherwise the message is dropped for it without waiting.
- The first message it takes again clears the mark and logs that the channel is read again; from then on it is waited on as before.
- Unregistering a channel clears its mark.

A consumer that blocks for longer than the timeout now misses messages sent while its channel stays full, where it used to get them up to one timeout late each.

## Impact

- `mqtt_transport.go`, `mqtt_transport_defs.go`: channel delivery, a stuck mark per channel name.
