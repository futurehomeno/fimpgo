## Why

JSON EC key import parses hex with `big.Int.SetString`, which accepts a leading sign, and `FillBytes` then drops it. A private key `"-1"` is imported as the scalar `1`, so a different key is accepted than the one written.

A failed JSON import also overwrote the stored key with nil, so a later sign or verify panicked on the nil key instead of using the previous valid one.

## What Changes

- JSON EC key import rejects hex values with a sign.
- A failed JSON EC key import keeps the previous key.

## Impact

- `security/ecdsa-key.go`: `hexToFixed`, `ImportJsonPublicKey`, `ImportJsonPrivateKey`.
