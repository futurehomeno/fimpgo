## Why

A second `MqttTransport.Start` on a started transport hangs for good. paho accepts `Connect` on a connected client at once without calling the connect handler, so `Start` times out waiting for a connection, because `connState.Init` replaced the connected and done channels. Its cleanup `Stop` then closes only the new done channel and waits forever for the first reader, which still waits on the old one.

## What Changes

- `Start` on a started, not stopped transport keeps the running session and only waits up to the timeout for it to be connected.

## Impact

- `mqtt_transport.go`: `Start`.
