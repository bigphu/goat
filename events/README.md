# `events`

> **Needs update:** ported from storeit and due to be re-implemented by hand. The ported tests were removed, so this package currently has no tests. Comments still use storeit's modules, `contract` packages and `cmd/api`. The outbox and delivery code has never run against Postgres in goat.

`events` sends events between modules through a transactional outbox, with River handling delivery.

## Quick Start

### Emitting

In a repository, inside the same transaction as the change:

```go
return database.WithTx(ctx, r.pool, func(tx pgx.Tx) error {
    // ... update ...
    e, err := events.New(contract.EventAssetCheckedOut, "asset", a.ID, payload)
    if err != nil {
        return err
    }
    return r.outbox.Append(ctx, tx, e)
})
```

`Append` fills `ActorID` from the context and inserts the event plus one `HandleEventArgs` job per subscriber.

### Subscribing

```go
registry := events.NewRegistry()
registry.OnAll("activity.record", subscriber.Record)
registry.On("notifications.email", subscriber.Email, contract.EventAssetCheckedOut)
```

### Reading a Payload

```go
p, err := events.Decode[contract.AssetCheckedOut](e)
```

### Wiring

The API and the worker must build the same registry.

```go
outbox := events.NewOutbox(registry, insertClient)              // API and worker
events.RegisterWorker(workers, pool, registry, loadActor)       // worker only
```

## API Reference

* `New(eventType, aggregateType, aggregateID, payload) (Event, error)`: Builds an event and encodes its payload to JSON.
* `Decode[T](e) (T, error)`: Decodes an event's payload.
* `NewRegistry()`, `(*Registry).On(subscriber, handler, eventTypes...)`, `(*Registry).OnAll(subscriber, handler)`, `(*Registry).For(eventType)`
* `NewOutbox(registry, client)`, `(*Outbox).Append(ctx, tx, events...)`
* `RegisterWorker(workers, pool, registry, loadActor)`: Adds the delivery worker, which restores the original actor with `jobs.RestoreActor`.
* `ListForAggregate(ctx, q, aggregateType, id, before, limit)`: Reads an aggregate's events newest first. `before` is the ID of the last event on the previous page, or nil for the first page. `q` is a pool or a transaction.

## Rules

* Handlers must be idempotent and must not rely on event order, because River can deliver an event more than once.
* Each subscriber name is registered once. It is the job key, so don't rename it while jobs are pending.
* An unknown subscriber retries, because the API may deploy before the worker. A missing event or a deleted actor cancels the job.
* Register everything at startup. Registration is not safe to run concurrently.

## Schema

The `platform.events` table comes from [`migrations/`](../migrations). It keeps every event after River cleans up the jobs, so events can be replayed for a new subscriber. `events/db` is sqlc output from `events/queries`, so don't edit it by hand. Regenerate it with `sqlc generate`, using the root `sqlc.yml`. Its `DBTX` interface can be faked to test worker logic without Postgres.
