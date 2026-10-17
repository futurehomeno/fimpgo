## ADDED Requirements

### Requirement: Failed X509 key import keeps the loaded key

When an X509 key import fails, `EcdsaKey` SHALL return an error and SHALL keep the key it held before the import, so later signing and verification do not panic.

#### Scenario: X509 public key of another type

- **WHEN** `ImportX509PublicKey` is given a PEM-encoded public key that is not ECDSA
- **THEN** it returns an error
- **AND** the previously loaded public key still verifies signatures

#### Scenario: Non-PEM X509 private key

- **WHEN** `ImportX509PrivateKey` is given input that is not PEM
- **THEN** it returns an error without panicking
- **AND** the previously loaded private key still signs

#### Scenario: Unparsable X509 private key

- **WHEN** `ImportX509PrivateKey` is given a PEM block whose bytes are not an EC private key
- **THEN** it returns an error
- **AND** the previously loaded private key still signs
