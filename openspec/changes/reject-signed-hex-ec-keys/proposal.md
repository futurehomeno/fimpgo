## Why

JSON EC key import parses hex with `big.Int.SetString`, which accepts a leading sign, and `FillBytes` then drops it. A private key `"-1"` is imported as the scalar `1`, so a different key is accepted than the one written.

## What Changes

- JSON EC key import rejects hex values with a sign.

## Impact

- `security/ecdsa-key.go`: `hexToFixed`.
