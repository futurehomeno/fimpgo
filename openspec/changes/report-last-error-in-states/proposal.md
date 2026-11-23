## Why

Review of the 1.18.0 release (PR #69) found that `Lifecycle.GetAllStates` always returns an empty `LastErrorText`: the text stored with `SetLastError` never reaches the app state report.

## What Changes

- `GetAllStates` returns the last error text in `LastErrorText`. `LastErrorCode` stays empty: nothing sets a code.

## Impact

- `edgeapp/lifecycle.go`: `GetAllStates`. An app that sets a last error now shows it in its state report until it sets another one.
