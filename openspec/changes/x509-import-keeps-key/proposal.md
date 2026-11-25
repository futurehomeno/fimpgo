## Why

`EcdsaKey.ImportX509PublicKey` stores nil as the public key when the PEM holds a non-ECDSA key, and `ImportX509PrivateKey` stores nil as the private key when parsing fails. A later `Sign`/`Verify` with that key then panics. `ImportX509PrivateKey` also panics on its own when the input is not PEM, because `pem.Decode` returns a nil block.

## What Changes

- Both X509 imports parse into a local and assign it only on success, so a failed import returns an error and keeps the previously loaded key.
- `ImportX509PrivateKey` returns an error for non-PEM input.

## Impact

- `security/ecdsa-key.go`: `ImportX509PublicKey`, `ImportX509PrivateKey`.
