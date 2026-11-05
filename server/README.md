# `server`

> **Needs update:** ported from an earlier application and due to be re-implemented by hand. The ported tests were removed, so this package currently has no tests.

`server` runs an `http.Server` with timeouts and shuts it down gracefully when the context is cancelled. Its router comes with the shared middleware chain and returns 404 and 405 as problem+json.

## Quick Start

`Run` doesn't catch signals, so `main` turns them into a cancelled context:

```go
ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
defer stop()

srv := server.New(cfg.HTTP, slog.Default())
// Mount module routes on srv.Router().
if err := srv.Run(ctx); err != nil {
    log.Fatal(err)
}
```

## API Reference

* `New(cfg, log) *Server`: Builds the server. A nil `log` uses `slog.Default()`.
* `(*Server).Router() chi.Router`: Where modules mount their routes.
* `(*Server).Run(ctx) error`: Serves on `cfg.Addr` until `ctx` is cancelled. A busy port returns an error straight away.
* `(*Server).Serve(ctx, ln) error`: The same on an existing listener, such as `:0` in tests.
* `(*Server).Handler() http.Handler`: For `httptest`, without opening a port.

## Middleware Chain

`New` installs, in order:

```
BodyLimit(HTTP_MAX_BODY_BYTES) -> NoSniff -> ClientIP(trusted) -> RequestID
  -> RequestLogger(log) -> Recoverer() -> Compress(json) -> chi routing
```

`Compress` covers `application/json` and `application/problem+json` when the client accepts it, so the access log's `bytes` is then the compressed size. Other types pass through. No response mixes a secret with caller-supplied text, which a BREACH attack would need, so keep it that way.

chi builds the chain when the first route is added, so on a router with no routes, requests skip every middleware. This only comes up in tests.

## Configuration

| Env var | Default | Notes |
| --- | --- | --- |
| `HTTP_ADDR` | :8080 | |
| `HTTP_READ_HEADER_TIMEOUT` | 10s | |
| `HTTP_READ_TIMEOUT` | 30s | Upload routes extend it with `middleware.ReadTimeout` |
| `HTTP_IDLE_TIMEOUT` | 60s | |
| `HTTP_SHUTDOWN_TIMEOUT` | 15s | |
| `HTTP_MAX_BODY_BYTES` | 1048576 | Upload routes extend it with `middleware.BodyLimit` |
| `HTTP_TRUSTED_PROXIES` | empty | CIDRs, comma-separated. Empty means `X-Forwarded-For` is never believed |
| `HTTP_API_DOCS` | false | `server` only reads this into `Config.APIDocs`. The binary checks it and calls `web.MountDocs`. The page is public, so use it in development only |

There is deliberately no write timeout, because it would also cut long downloads.
