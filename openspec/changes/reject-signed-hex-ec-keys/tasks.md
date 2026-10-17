## 1. Signed hex

- [x] 1.1 Failing test: JSON private and public keys with a leading `-` are rejected
- [x] 1.2 `hexToFixed` rejects negative values
- [x] 1.3 Failing test: a JSON private key with a leading `+` is rejected
- [x] 1.4 `hexToFixed` rejects any sign character, including `+` and `-0`

## 2. Failed import

- [x] 2.1 Failing test: after a failed JSON import the previous private and public keys still sign and verify
- [x] 2.2 JSON imports assign the parsed key only on success
