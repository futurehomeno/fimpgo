# security Specification

## Purpose
ECDSA P-256 keys, their storage, and signed FIMP messages, so a hub app can verify that a command
comes from a known user device.

## Requirements
### Requirement: Keys import and export in PEM and JSON form

A generated key pair SHALL use the P-256 curve. Keys SHALL export as PEM (`PRIVATE KEY` holding SEC 1
EC, `PUBLIC KEY` holding PKIX) and as JSON objects with hex coordinates `x`, `y` and, for the private
key, the scalar `d`. Importing JSON SHALL accept coordinates whose leading zeros were dropped. A
public PEM key that is not ECDSA SHALL be rejected.

#### Scenario: JSON round trip

- **WHEN** a generated key pair is exported as JSON and imported again
- **THEN** a signature made with the imported private key verifies with the imported public key

#### Scenario: Short coordinate

- **WHEN** a public key whose `x` coordinate has a leading zero byte is exported without it and imported
- **THEN** the import succeeds and yields the same key

### Requirement: Signed messages wrap the original

Signing SHALL serialize the original message, carry it base64-encoded as a `bin` value in a
`cmd.transport.signed` message for a command or an `evt.transport.signed` message otherwise, keep
the original service, and add the ES256 signature over the encoded value as base64url property `sig`
to the given properties, or to `user_id` alone when none are given. Verifying SHALL accept only those two interfaces, require `sig`, check it with
the given public key, and return the original message.

#### Scenario: Valid signature

- **WHEN** a signed command is verified with the signer's public key
- **THEN** the original command is returned

#### Scenario: Tampered payload

- **WHEN** the wrapped value of a signed message is altered
- **THEN** verification fails

#### Scenario: Not a signed message

- **WHEN** an `evt.presence.report` message is verified
- **THEN** verification fails with `incorrect message type`

### Requirement: The key store persists keys per user and device

The key store SHALL hold serialized keys per user, device, key type and algorithm, add or replace
them, and write the whole store to disk on each change, creating the file with mode 0600 for a
private store and 0664 otherwise.

#### Scenario: Upsert

- **WHEN** a public key is upserted twice for the same user, device, type and algorithm
- **THEN** the store holds one record with the newer key

#### Scenario: Private store

- **WHEN** a private store is saved for the first time
- **THEN** its file is readable only by its owner
