## 1. Signed hex

- [x] 1.1 Failing test: JSON private and public keys with a leading `-` are rejected
- [x] 1.2 `hexToFixed` rejects negative values

## 2. Failed import

- [x] 2.1 Failing test: after a failed JSON import the previous private and public keys still sign and verify
- [x] 2.2 JSON imports assign the parsed key only on success
