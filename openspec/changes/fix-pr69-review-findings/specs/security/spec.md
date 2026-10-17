## ADDED Requirements

### Requirement: A private key provides its public key

`EcdsaKey.PublicKey()` SHALL return the public key of the private key when no public key was set, whether the private key came from `ImportJsonPrivateKey`, `ImportX509PrivateKey` or `SetPrivateKey`.

#### Scenario: Verify with an imported private key

- **WHEN** only a private key is loaded into an `EcdsaKey`
- **THEN** `PublicKey()` returns its public key
- **AND** a payload signed with the key verifies with the same key
- **AND** the X509 and JSON exports contain its public key
