# `middleware`

> **Needs update:** ported from storeit. `ErrRateLimited` uses goat's `errs`, with the same status, type, title and detail as before. This page has not otherwise been reviewed against goat's needs.

`middleware` holds the HTTP middleware the toolkit uses. `server.New` installs the shared chain, so most apps only add per-operation overrides.

## Quick Start

Settings for one OpenAPI operation go through `ForOperations` in the module's `Middlewares`, not `r.With`:

```go
upload := middleware.ForOperations(web.MustOperations(spec, "/api/v1", "POST /api/v1/imports"),
    middleware.BodyLimit(50<<20), middleware.ReadTimeout(10*time.Minute))
```

## API Reference

| Middleware | Purpose |
| --- | --- |
| `BodyLimit(n)` | Limits the request body. `n <= 0` means no limit. A later `BodyLimit` replaces the limit rather than stacking |
| `ReadTimeout(d)` | Sets the read deadline. `d <= 0` removes it |
| `ForOperations(ops, mws...)` | Applies middleware to the listed operations only |
| `ClientIP(trusted)` | Finds the client IP: the rightmost untrusted `X-Forwarded-For` entry, believed only from trusted proxies. Read it with `ClientIPFrom(ctx)` |
| `RequestID` | Keeps an `X-Request-Id` matching `[A-Za-z0-9._-]{1,64}`, otherwise generates a UUID, and echoes it in the response |
| `RequestLogger(log)` | Puts the logger in the context and writes one access log line: method, path, client_ip, route, status, bytes, dur, aborted |
| `Recoverer()` | Turns a panic into a logged 500 problem+json. If the response has started, it only logs. `http.ErrAbortHandler` is passed on to `net/http` |
| `NoSniff` | Sets `X-Content-Type-Options: nosniff` |
| `NoStore` | Sets `Cache-Control: no-store`, for responses carrying secrets such as tokens |
| `RateLimit(every, burst)` | Per client IP token bucket, see below |

### Rate Limiting

`RateLimit` holds `burst` requests per client IP and refills one per `every`. Over the limit, it answers `ErrRateLimited` (429, `/errors/rate-limited`) with `Retry-After` in seconds, rounded up. Rejected requests don't use up tokens.

Each call has its own counters, so attach one per operation with `ForOperations`. Counters live in process memory, so with several replicas each one counts separately.

## Rules

* Operation strings are `"METHOD <chi route pattern>"`. Wrap lists in `web.MustOperations` so a typo panics at startup.
* `Recoverer` logs through the context logger, so it goes after `RequestLogger`.
