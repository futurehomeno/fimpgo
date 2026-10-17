## 1. Hub token reload

- [x] 1.1 Failing test: a second `LoadHubTokenFromCB` after one without a CloudBridge answer times out waiting for the answer instead of failing to subscribe
- [x] 1.2 Drop the stopped transport and sync client in `LoadHubTokenFromCB`
- [x] 1.3 Failing test: a failed subscribe also releases the transport
- [x] 1.4 Release the transport on every return
- [x] 1.5 Failing test: two concurrent `LoadHubTokenFromCB` calls both time out without a race or panic
- [x] 1.6 Serialize `LoadHubTokenFromCB` with a mutex

## 2. Start with retry count 0

- [x] 2.1 Failing test: `Start` with retry count 0 connects to a reachable broker
- [x] 2.2 Treat a retry count below 1 as 1

## 3. Private key import

- [x] 3.1 Failing test: after `ImportJsonPrivateKey`, `ImportX509PrivateKey` or `SetPrivateKey`, also over a generated key, the key verifies its own signature and exports its public key
- [x] 3.2 Loading a private key clears the stored public key; `PublicKey()` falls back to the private key's public key
