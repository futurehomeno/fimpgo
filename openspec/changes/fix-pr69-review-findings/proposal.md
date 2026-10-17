## Why

Review of the 1.18.0 release (PR #69) found three defects:

- `FhOAuth2Client.LoadHubTokenFromCB` stops its MQTT transport after each request but keeps the stopped transport and sync client. Every later call (a retry after `Init` got no answer from CloudBridge, or `postMsg` with an empty hub token) subscribes on the stopped client and fails, so the hub token can never be loaded again.
- `MqttTransport.Start` sizes its budget from the start retry count. A count of 0 gives a zero budget, so `Start` returns `ErrTimeout` at once without waiting for the connect attempt.
- `EcdsaKey.ImportJsonPrivateKey`, `ImportX509PrivateKey` and `SetPrivateKey` leave the public key unset although the private key carries it, so `PublicKey()` returns nil and `VerifyStringES256` or `ExportJsonEncodedKeys` panic on a key loaded this way.

## What Changes

- `LoadHubTokenFromCB` stops and drops its transport and sync client on every return, including a failed subscribe, so the next call connects again.
- `Start` treats a retry count below 1 as 1.
- `PublicKey()` falls back to the private key's public key when none was set.

## Impact

- `edgeapp/auth.go`: `LoadHubTokenFromCB`.
- `mqtt_transport.go`: `connect`.
- `security/ecdsa-key.go`: `PublicKey`, `ExportX509EncodedKeys`, `ExportJsonEncodedKeys`.
