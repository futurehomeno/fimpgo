## ADDED Requirements

### Requirement: JSON private key import sets the public key

`EcdsaKey.ImportJsonPrivateKey` SHALL set the public key derived from the imported private key.

#### Scenario: Verify with an imported private key

- **WHEN** a private key is imported with `ImportJsonPrivateKey`
- **THEN** `PublicKey()` returns its public key
- **AND** a payload signed with the key verifies with the same key
