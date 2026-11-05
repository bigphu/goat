# `auth`

> **Needs update:** ported from storeit. The error variables use goat's `errs`, with the same types and titles as before. Package comments and permission examples still refer to storeit modules.

`auth` keeps the caller (`Actor`) in the request context, reads bearer tokens, and checks permissions.

## Quick Start

Build the middleware once, listing the operations that need no token. Attach it through oapi-codegen's `Middlewares`, not `r.Use`, because public operations are matched by chi route pattern and that is only known after routing. Mounted with `r.Use`, it fails closed and every route requires a token.

```go
public := web.MustOperations(spec, "/api/v1", "POST /api/v1/auth/login") // panics on a typo
authMW := auth.Middleware(tokens, auth.Public(public...))

api.HandlerWithOptions(strict, api.ChiServerOptions{
    BaseRouter:       r,
    BaseURL:          "/api/v1",
    // The last entry runs first, so a missing token is a 401 rather than a 422.
    Middlewares:      []api.MiddlewareFunc{web.ValidateRequests(spec, "/api/v1"), authMW},
    ErrorHandlerFunc: web.RequestError,
})
```

Check permissions in the service, so the check holds whether a handler, worker or subscriber calls it:

```go
actor, err := auth.Require(ctx, domain.PermAssetCreate) // 401 without an actor, 403 without the permission
```

Scheduled jobs run as `SystemActor`, which has every permission and no account:

```go
ctx = auth.WithActor(ctx, auth.SystemActor)
```

## API Reference

| Function | Purpose |
| --- | --- |
| `Middleware(tokens, opts...)` | Verifies the bearer token, puts the `Actor` in the context and adds `actor_id` to the log scope |
| `Public(ops...)` | Option listing operations that need no token, as `"METHOD <route pattern>"` |
| `Require(ctx, perm)` | Returns the actor, or an error if it is missing or lacks `perm` |
| `WithActor(ctx, a)` / `FromContext(ctx)` | Stores and reads the actor |
| `ActorIDFrom(ctx)` | Account ID, or `uuid.Nil` for the system actor or no actor |

`Actor` has the fields `AccountID` and `Permissions`, and the methods `Can(perm)` and `IsSystem()`.

### Errors

| Variable | Status | Type |
| --- | --- | --- |
| `ErrUnauthenticated` | 401 | `/errors/unauthenticated` |
| `ErrMissingToken` | 401 | `/errors/missing-token` |
| `ErrInvalidToken` | 401 | `/errors/invalid-token` |
| `ErrForbidden` | 403 | `/errors/forbidden` |
