# `database`

> **Needs update:** ported from storeit. The code is unchanged apart from imports. Package comments still use storeit's asset examples.

`database` opens the shared Postgres pool (`pgxpool`) and runs transactions.

## Quick Start

Open the pool once in `main`. `Open` pings, so a wrong URL or a database that is down stops the app at startup.

```go
pool, err := database.Open(ctx, cfg.DB)
if err != nil {
    log.Fatal(err)
}
defer database.Close(pool)
```

Repositories take the pool. A method that writes does all of its work inside one `WithTx`, including any outbox event, so the change and the event commit or disappear together. Services never see a transaction.

```go
func (r *AssetRepository) Save(ctx context.Context, a *domain.Asset) error {
    return database.WithTx(ctx, r.pool, func(tx pgx.Tx) error {
        // ... write ...
        return r.outbox.Append(ctx, tx, events...)
    })
}
```

## API Reference

* `Open(ctx, cfg) (*pgxpool.Pool, error)`: Validates, applies the pool settings (overriding `pool_*` in the URL), reads `PGPASSWORD_FILE`, then pings.
* `Close(pool)`: Closes the pool.
* `WithTx(ctx, db, fn) error`: Commits when `fn` returns nil. Rolls back on an error, a panic or `runtime.Goexit`. The error from `fn` is returned intact, so domain errors keep their status. Rollback uses a non-cancelled context.
* `Beginner`: `*pgxpool.Pool`, or a `pgx.Tx` (which makes it a savepoint).

## Configuration

| Env var | Default |
| --- | --- |
| `DB_URL` | empty: pgx reads `PG*` |
| `DB_MAX_CONNS` | 25 |
| `DB_MIN_CONNS` | 0 |
| `DB_MAX_CONN_LIFETIME` | 1h |
| `DB_MAX_CONN_LIFETIME_JITTER` | 5m |
| `DB_MAX_CONN_IDLE_TIME` | 30m |
| `DB_HEALTH_CHECK_PERIOD` | 1m |
| `DB_CONNECT_TIMEOUT` | 5s |
| `PGPASSWORD_FILE` | password file, for Docker secrets |

For integration tests against a real database, see [`dbtest`](dbtest/README.md).
