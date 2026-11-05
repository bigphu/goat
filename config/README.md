# `config`

> **Needs update:** ported from an earlier application and due to be re-implemented by hand. The ported tests were removed, so this package currently has no tests.

`config` reads environment variables into a binary's config struct at startup and stops if anything is missing or invalid.

## Quick Start

Each package owns its config block: env var names, defaults and `Validate()`. A binary lists only the blocks it uses, with **no `envPrefix`**.

```go
type serverConfig struct {
    Log  logger.Config
    HTTP server.Config
    DB   database.Config
    JWT  jwt.Config
}

var cfg serverConfig
if err := config.Load(&cfg); err != nil {
    log.Fatal(err)
}
```

## API Reference

* `Load(cfg any) error`: Fills `cfg` from the environment, then validates every nested block.
* `Validator`: An interface with `Validate() error`, for a block with its own constraints.

## Rules

* `Load` validates children first, then the parent, and joins all errors with their field paths (`config: DB: db: ...`).
* A `Validate()` checks only its own fields and never calls its children's.
* A binary writes its own `Validate` only for constraints across blocks.
* In the ported blocks, zero means "use the default" (configs built by hand in tests) and negative values are rejected.

## Env Vars by Package

| Block | Documented in |
| --- | --- |
| `logger.Config` | [`logger`](../logger/README.md) |
| `server.Config` | [`server`](../server/README.md) |
| `database.Config` | [`database`](../database/README.md) |
| `jwt.Config` | [`jwt`](../jwt/README.md) |
| `storage.Config` | [`storage`](../storage/README.md) |
| `jobs.Config` | [`jobs`](../jobs/README.md) |
| `mail.Config` | [`mail`](../mail/README.md) |
