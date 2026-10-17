## Why

Review of the 1.18.0 release (PR #69) found three defects:

- `FhOAuth2Client.LoadHubTokenFromCB` stops its MQTT transport after each request but keeps the stopped transport and sync client. Every later call (a retry after `Init` got no answer from CloudBridge, or `postMsg` with an empty hub token) subscribes on the stopped client and fails, so the hub token can never be loaded again.
- `MqttTransport.Start` sizes its budget from the start retry count. A count of 0 gives a zero budget, so `Start` returns `ErrTimeout` at once without waiting for the connect attempt.
- `EcdsaKey.ImportJsonPrivateKey` leaves the public key unset although the parsed private key carries it, so `PublicKey()` returns nil and `VerifyES256` or `ExportJsonEncodedKeys` fail on a key imported this way.

## What Changes

- `LoadHubTokenFromCB` drops the stopped transport and sync client, so the next call connects again.
- `Start` treats a retry count below 1 as 1.
- `ImportJsonPrivateKey` also sets the public key.

## Impact

- `edgeapp/auth.go`: `LoadHubTokenFromCB`.
- `mqtt_transport.go`: `connect`.
- `security/ecdsa-key.go`: `ImportJsonPrivateKey`.
