## Why

`postMsg` reads the hub token once before its retry loop, so a token set with `SetHubToken` while it waits between attempts is never used: every remaining attempt sends the rejected token. Before the hub token lock was added, each attempt read the current token.

## What Changes

- Each attempt of a token exchange reads the current hub token.

## Impact

- `edgeapp/auth.go`: `postMsg`.
