# `storage`

> **Needs update:** ported from storeit and due to be re-implemented by hand. The ported tests were removed, so this package currently has no tests. The code is unchanged apart from imports. The key examples are storeit's, and the package comment mentions an S3 implementation that doesn't exist.

`storage` stores files behind the `Store` interface. `Local`, on disk, is the only implementation.

## Quick Start

Services take a `Store`, and `main` chooses the implementation. Keys are built by code from IDs, never from a filename a user sent.

```go
store, err := storage.NewLocal(cfg.Storage)
if err != nil {
    log.Fatal(err)
}
defer store.Close()

key := "imports/" + id.String() + "/source.xlsx"
err = store.Put(ctx, key, file)
```

## API Reference

| Method | Behaviour |
| --- | --- |
| `Put(ctx, key, r)` | Writes all of `r` under `key`, replacing any existing file. It writes a temp file, fsyncs and renames, so a failure leaves the key as it was. It stops when `ctx` is cancelled |
| `Get(ctx, key)` | Opens the file, and the caller closes it. Returns `ErrNotFound` if it is missing |
| `Delete(ctx, key)` | Removes the file. A missing file is not an error, so it is safe to retry |

Invalid keys return `ErrInvalidKey`. `Local` also has `Close()`. All of its access goes through `os.Root`, so no key can leave the root directory, even through a symlink.

## Configuration

| Env var | Default |
| --- | --- |
| `STORAGE_DIR` | ./tmp/storage |
