## Why

Each received message goes to the registered channels one by one, under the channel registration lock, waiting up to the receive timeout (10 s) on each full channel. With N stuck consumers every message waits N × 10 s: healthy channels are delayed, registering and unregistering a channel block, and the receive queue overflows until the transport stops (`[fimpgo] Channel <flowId>_<nodeId> not read for N sec`, tpflow 2.3.0, 2,325 hits on one site).

## What Changes

- A message is offered to every interested channel without blocking first, so channels with room get it at once.
- Full channels are then waited on together against one shared deadline, in the receiving goroutine, so one message costs at most one receive timeout.
- A channel closed by its consumer during that wait is logged as a handler crash, as before. The registration lock is now released too; before, it stayed held and the next message blocked forever.

## Impact

- `mqtt_transport.go`: delivery of incoming messages to channels.
