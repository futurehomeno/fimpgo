## 1. X509 public key import

- [x] 1.1 Failing test: importing a non-ECDSA public key returns an error and keeps the previous public key
- [x] 1.2 Assign the parsed public key only after the type check passes

## 2. X509 private key import

- [x] 2.1 Failing test: non-PEM and unparsable PEM input return an error without panicking and keep the previous private key
- [x] 2.2 Check the PEM block and assign the parsed private key only on success
