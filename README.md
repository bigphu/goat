# goat

**Go Application Toolkit**: a collection of building blocks for Go web applications, built around Postgres, River background jobs, a transactional outbox and OpenAPI-first HTTP APIs.

> **Status: early prototype.** Only `errs` is written specifically for goat. The other packages were ported from an earlier application and are due to be re-implemented by hand. **The repo has no tests yet**, and the Postgres-backed code has not been run here. Expect breaking changes.

## Installation

Requires Go 1.26 or later.

```bash
go get github.com/bigphu/goat@v0.0.1
```

## Quick Start

A server whose errors are returned as `application/problem+json`:

```go
package main

import (
	"context"
	"log"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/bigphu/goat/errs"
	"github.com/bigphu/goat/server"
	"github.com/bigphu/goat/web"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	srv := server.New(server.Config{}, slog.Default())
	srv.Router().Get("/things/{id}", web.Handle(func(w http.ResponseWriter, r *http.Request) error {
		return errs.NotFound("No thing has that ID.")
	}))

	if err := srv.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
```

`GET /things/1` then answers 404 with:

```json
{"type":"/errors/not-found","title":"Not Found","status":404,"detail":"No thing has that ID."}
```

## Packages

Listed so that each package only depends on packages above it, which also makes this a workable order for the re-implementation. Each link opens the package's README.

| Package | Purpose | Status |
| --- | --- | --- |
| [`errs`](errs/README.md) | Structured, HTTP-aware errors | Written for goat |
| [`config`](config/README.md) | Environment variables into config structs, with validation | To re-implement |
| [`logger`](logger/README.md) | `slog` setup and request-scoped attributes | To re-implement |
| [`jwt`](jwt/README.md) | HS256 access tokens with a key ring | To re-implement |
| [`storage`](storage/README.md) | File storage behind an interface, on local disk | To re-implement |
| [`database`](database/README.md) | `pgxpool` from config and transactions | To re-implement |
| [`database/dbtest`](database/dbtest/README.md) | Real Postgres for integration tests | To re-implement |
| [`mail`](mail/README.md) | Email over SMTP, Resend or the log | To re-implement |
| [`web`](web/README.md) | problem+json errors, OpenAPI validation and JSON helpers | To re-implement |
| [`web/apicommon`](web/apicommon/README.md) | Generated types for shared API schemas | To re-implement |
| [`middleware`](middleware/README.md) | Body limit, client IP, request ID, access log, rate limit and more | To re-implement |
| [`server`](server/README.md) | `http.Server` with timeouts and graceful shutdown | To re-implement |
| [`auth`](auth/README.md) | Actor in context, bearer-token middleware and permission checks | To re-implement |
| [`jobs`](jobs/README.md) | Background jobs on River | To re-implement |
| [`events`](events/README.md) | Transactional outbox and subscriber delivery | To re-implement |

[`docs/platform.md`](docs/platform.md) is the original reference across all of the packages, with notes on what is out of date.

## Development

```bash
go build ./...
go vet ./...
```

Regenerate the generated code after changing `api/common.yaml` or the files in `events/queries` and `migrations`:

```bash
go generate ./web/apicommon
```

```bash
go run github.com/sqlc-dev/sqlc/cmd/sqlc@v1.31.1 generate
```

`migrations/` holds the schema `events` needs. `api/common.yaml` holds the API schemas that modules share.
