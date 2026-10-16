## Why

Two transport problems show up in hub logs:

- When the first connect attempt in `Start` times out, paho is still connecting (connect retry is on), but every later attempt calls `Connect` again. paho rejects that, so boot logs show `MQTT connect failed 2/10 err: status can only transition to connecting from disconnected` for every remaining attempt. If all attempts run out, `Start` returns while paho keeps connecting in the background. The transport then connects, but nothing reads its incoming messages.
- `Processing payload from topic=... err: Key path not found` / `Unknown value type` logs only the topic, so the producer of a malformed message can't be identified from field logs.

## What Changes

- `Start` keeps waiting on a connect attempt that is still pending, and starts a new one only after the previous attempt has finished with an error. When all attempts are used up, it abandons the pending attempt, so the next `Start` begins clean.
- A payload that fails to parse is logged with its length and its top-level JSON keys (no values).

## Impact

- `mqtt_transport.go`: `Start`, incoming-payload error log.
