# `jobs`

> **Needs update:** ported from storeit. The code is unchanged apart from imports. The fixed queues and worker config are storeit's choices, and comments still use storeit's import job and `cmd/api`. Its database-backed tests have not been run in goat yet.

`jobs` puts background work on a queue (River on Postgres). Services and repositories see only the `Enqueuer` interface and never import River.

## Quick Start

### Defining a Job

Args are a JSON-serializable struct. `Kind()` must never change, because stored jobs find their worker by kind. A job a user asked for embeds `ActorArgs`.

```go
type CommitImportArgs struct {
    jobs.ActorArgs
    ImportID uuid.UUID `json:"import_id"`
}

func (CommitImportArgs) Kind() string { return "inventory.commit_import" }
```

### Creating a Client

```go
client, err := jobs.NewInsertClient(pool, log)                              // API: insert only
client, err := jobs.NewWorkerClient(pool, log, cfg.Jobs, workers, periodic) // worker
enq := jobs.NewRiver(client)
```

### Enqueuing

A job that goes with a data change is enqueued in the same transaction, so it exists only if the transaction commits:

```go
return database.WithTx(ctx, r.pool, func(tx pgx.Tx) error {
    // ... write ...
    return r.jobs.EnqueueTx(ctx, tx, CommitImportArgs{
        ActorArgs: jobs.ActorArgsFrom(ctx),
        ImportID:  id,
    }, jobs.WithUniqueArgs())
})
```

A job with no related change uses `Enqueue`, without that guarantee.

### Working

A worker restores the actor, then calls one service method. Permissions are reloaded at run time.

```go
func (w *CommitImportWorker) Work(ctx context.Context, j *river.Job[CommitImportArgs]) error {
    ctx, err := jobs.RestoreActor(ctx, j.Args.ActorID, w.loadActor)
    if err != nil {
        return err
    }
    return w.svc.CommitImport(ctx, j.Args.ImportID)
}
```

## API Reference

### Options

| Option | Purpose |
| --- | --- |
| `WithQueue(queue)` | Chooses a queue. It must be declared on the worker client, or its jobs wait forever |
| `WithPriority(n)` | 1 (highest) to 4. Anything else makes the insert fail |
| `WithSchedule(t)` | Runs the job at `t` |
| `WithMaxAttempts(n)` | Overrides the retry limit |
| `WithUniqueArgs()` | Merges with an unfinished job of the same kind and args. Off by default |

### Other Names

* `QueueDefault`, `QueueEvents`: The queues the worker client runs. A new queue must be added in `workerConfig`.
* `DefaultMaxAttempts`: 10.
* `ActorLoader`: Loads an actor for `RestoreActor`. An error wrapping `ErrActorNotFound` cancels the job, and any other error retries.

## Rules

* Workers can run a job more than once (retries, crashes), so their work must be idempotent.
* Scheduled jobs run as `auth.SystemActor`.

## Configuration

Used by the worker only.

| Env var | Default |
| --- | --- |
| `JOBS_DEFAULT_MAX_WORKERS` | 10 |
| `JOBS_EVENTS_MAX_WORKERS` | 10 |

The worker total plus a few River connections must fit within the worker's `DB_MAX_CONNS`, or jobs wait for a connection.
