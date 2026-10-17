## ADDED Requirements

### Requirement: JSON EC keys use unsigned hex

JSON EC key import SHALL accept only unsigned hexadecimal values for `x`, `y` and `d`, and SHALL return an error for any value with a sign.

#### Scenario: Signed private scalar

- **WHEN** a JSON private key with `d` set to `-1` or `+1` is imported
- **THEN** the import returns an error

#### Scenario: Signed public coordinate

- **WHEN** a JSON public key whose `x` is a valid coordinate prefixed with `-` is imported
- **THEN** the import returns an error

### Requirement: Failed JSON EC key import keeps the previous key

A JSON EC key import that returns an error SHALL leave the key it would have replaced unchanged.

#### Scenario: Failed public import

- **WHEN** a valid JSON public key is imported and then a JSON public key whose point is not on the curve is imported
- **THEN** the second import returns an error
- **AND** the first public key still verifies signatures

#### Scenario: Failed private import

- **WHEN** a valid JSON private key is imported and then a JSON private key with `d` set to `0` is imported
- **THEN** the second import returns an error
- **AND** the first private key still signs
