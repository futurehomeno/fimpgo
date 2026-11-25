## ADDED Requirements

### Requirement: A private key provides its public key

`EcdsaKey.PublicKey()` SHALL return the public key of the private key loaded last by `Generate`, `ImportJsonPrivateKey`, `ImportX509PrivateKey` or `SetPrivateKey`, unless a public key was imported after it.

#### Scenario: Verify with an imported private key

- **WHEN** a private key is loaded into an `EcdsaKey`, also one that held another key before
- **THEN** `PublicKey()` returns its public key
- **AND** a payload signed with the key verifies with the same key
- **AND** the X509 and JSON exports contain its public key
