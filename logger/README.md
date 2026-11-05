# `logger`

> **Needs update:** ported from an earlier application and due to be re-implemented by hand. The ported tests were removed, so this package currently has no tests.

`logger` builds a `*slog.Logger` that writes JSON in production and a pretty format in development, and carries request-scoped attributes in the context.

## Quick Start

Call `New` once in `main`, then use `slog` directly. Attributes added with `With` appear on every log written with that context.

```go
slog.SetDefault(logger.New(os.Stdout, cfg.Log))

ctx = logger.With(ctx, slog.String("request_id", id))
slog.InfoContext(ctx, "order placed", "order_id", o.ID)
```

Attributes that are only known deep inside a request also reach logs written with outer contexts, such as the access log, when added with `AddToScope`:

```go
logger.AddToScope(r.Context(), slog.String("actor_id", id))
```

## API Reference

* `New(w, cfg) *slog.Logger`: Builds the logger.
* `With(ctx, attrs...)`: Adds attributes to every log written with the returned context.
* `WithScope(ctx)`, `AddToScope(ctx, attrs...)`: Open a scope, then add attributes to it from anywhere inside.
* `NewContext(ctx, log)`, `FromContext(ctx)`: Store and read the request's logger. `FromContext` falls back to `slog.Default()`.
* `NewPrettyHandler(w, opts)`: The handler behind the pretty format.

Context attributes always stay at the top level, even inside `log.WithGroup(...)`.

## Configuration

| Env var | Default | Values |
| --- | --- | --- |
| `LOG_LEVEL` | info | debug, info, warn, error |
| `LOG_FORMAT` | json | json, pretty |
