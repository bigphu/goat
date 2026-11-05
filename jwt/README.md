# `jwt`

> **Needs update:** ported from storeit. The code is unchanged apart from imports. The `JWT_AUDIENCE` default is still `storeit-api`.

`jwt` issues and verifies HS256 access tokens, with a key ring for rotation.

## Quick Start

```go
p, err := jwt.New(cfg)

tok, err := p.Issue(accountID, perms) // send tok.Value to the client

claims, err := p.Verify(raw)
if errors.Is(err, jwt.ErrExpired) {
    // ...
}
id, err := claims.UserID()
perms := claims.Permissions
```

## API Reference

* `New(cfg) (*Provider, error)`: Builds the provider and loads the key ring.
* `(*Provider).Issue(accountID, permissions) (Token, error)`: Signs a token. `Token` has `Value`, `ID` and `ExpiresAt`.
* `(*Provider).Verify(raw) (*Claims, error)`: Accepts only HS256 with `typ: at+jwt`, matching `iss` and `aud`, and an `exp`, with 30s of clock leeway.
* `(*Claims).UserID()`, `(*Claims).TokenID()`: Parse the IDs as UUIDs. `Claims.Permissions` holds the permissions.

### Errors

Match these with `errors.Is`: `ErrMalformed`, `ErrSignature`, `ErrExpired`, `ErrNotYetValid`, `ErrClaims`, `ErrUnknownKID`, `ErrBadAlgorithm`, `ErrBadType`.

A token signed with another algorithm, including `none`, gets `ErrSignature`.

## Key Rotation

`JWT_KEYS_FILE` holds `kid:base64` pairs, such as `v1:<base64>,v2:<base64>`, separated by commas or newlines. Secrets are at least 32 bytes and kids must be unique.

To rotate, add the new key to the file and switch `JWT_ACTIVE_KID`. New tokens are signed with the active key, and old tokens still verify while their kid stays in the file.

Permissions live in the token, so a role change only takes effect when the token expires. Keep the TTL short.

## Configuration

| Env var | Default |
| --- | --- |
| `JWT_KEYS_FILE` | required |
| `JWT_ACTIVE_KID` | required |
| `JWT_ISSUER` | required |
| `JWT_AUDIENCE` | storeit-api |
| `JWT_ACCESS_TOKEN_TTL` | 15m |
