# `dbtest`

> **Needs update:** ported from storeit. It now migrates from `goat/migrations`, which holds only the two `platform.events` migrations. The test database is still named `storeit`, and comments still refer to storeit's `compose.yml` and `cmd/migrate`. The tests that use it have not been run in goat yet.

`dbtest` starts a real Postgres for integration tests with testcontainers, with the goose and River migrations applied.

## Quick Start

```go
func TestSave(t *testing.T) {
    pool := dbtest.Pool(t)
    // ...
}
```

## API Reference

* `Pool(t *testing.T) *pgxpool.Pool`: Returns the pool for the package's shared test database.

## Rules

* There is one container per test binary, shared between tests. Don't assume an empty database: create data with fresh IDs and check only your own data.
* Without Docker the tests **skip**. With `CI=true` they fail, so database tests can't be silently dropped.
* Several packages starting containers at once can fail on Windows, so run them one at a time and keep Docker running:

```bash
CI=true go test -p 1 ./...
```
