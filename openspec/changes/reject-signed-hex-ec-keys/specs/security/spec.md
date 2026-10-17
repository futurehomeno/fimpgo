## ADDED Requirements

### Requirement: JSON EC keys use unsigned hex

JSON EC key import SHALL accept only unsigned hexadecimal values for `x`, `y` and `d`, and SHALL return an error for any value with a sign.

#### Scenario: Signed private scalar

- **WHEN** a JSON private key with `d` set to `-1` is imported
- **THEN** the import returns an error

#### Scenario: Signed public coordinate

- **WHEN** a JSON public key whose `x` is a valid coordinate prefixed with `-` is imported
- **THEN** the import returns an error
