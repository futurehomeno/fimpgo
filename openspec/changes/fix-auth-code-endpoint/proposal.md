## Why

`ExchangeCodeForTokens` posts the authorization code to the refresh endpoint (`/proxy/refresh`) instead of `/proxy/auth-code`, so an authorization-code login through `FhOAuth2Client` never gets tokens. Its test points both endpoints at one server, which hid it.

## What Changes

- `ExchangeCodeForTokens` posts to the auth-code endpoint.

## Impact

- `edgeapp/auth.go`: `ExchangeCodeForTokens`.
