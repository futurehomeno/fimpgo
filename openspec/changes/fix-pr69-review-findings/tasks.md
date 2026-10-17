## 1. Hub token reload

- [x] 1.1 Failing test: a second `LoadHubTokenFromCB` after one without a CloudBridge answer times out waiting for the answer instead of failing to subscribe
- [x] 1.2 Drop the stopped transport and sync client in `LoadHubTokenFromCB`

## 2. Start with retry count 0

- [x] 2.1 Failing test: `Start` with retry count 0 connects to a reachable broker
- [x] 2.2 Treat a retry count below 1 as 1

## 3. JSON private key import

- [x] 3.1 Failing test: after `ImportJsonPrivateKey` the key verifies its own signature
- [x] 3.2 Set the public key from the parsed private key
