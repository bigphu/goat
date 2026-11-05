# `errs`

`errs` is a structured, HTTP-aware error package for Go. It is designed to represent problem details for HTTP APIs
(inspired by [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) and [RFC 9457](https://datatracker.ietf.org/doc/html/rfc9457/)),
allowing errors to carry statuses, types, titles, details, and field-specific validation contexts.

By utilizing the functional options pattern, `errs` provides a clean, unified, and extensible API for building complex
API errors while maintaining full compatibility with the standard Go `errors` package (`errors.Is`, `errors.As`, `Unwrap`).

## Features

* **HTTP-Aware:** Built-in helper functions for common HTTP statuses (400, 401, 404, etc.).
* **Minimal Call Sites:** The detail is the only positional argument; everything else is an option.
* **Sensible Titles:** The title defaults to the status text (`404` → "Not Found"), so you only write the occurrence-specific `detail`.
* **Error Wrapping:** Seamlessly wrap underlying causes to preserve logs while presenting clean errors to users.
* **Validation Context:** `Validation`, `FieldRequired` and `FieldInvalid` attach field errors that are reachable through `errors.Is` and `errors.As`.
* **Standard Library Compatible:** Implements `Error()`, `Unwrap() []error`, and `Is()`. Requires Go 1.20+.

---

## Quick Start

### Creating a Basic HTTP Error

Use the HTTP helper functions to quickly generate contextual errors. The detail is the first argument; the title is derived from the status unless you override it.

```go
package main

import (
	"fmt"
	"github.com/bigphu/goat/errs"
)

func main() {
	err := errs.NotFound("No user exists with the ID '123'.")
	
	fmt.Println(err.Error()) 
	// Output: Not Found: No user exists with the ID '123'.
}

```

### Wrapping an Existing Error

You can wrap underlying database or system errors to preserve the original stack or context, while safely controlling what gets exposed in the HTTP response.

```go
func getUser(id string) error {
    dbErr := sql.ErrNoRows // Simulated database error
    
    return errs.NotFound(
        "The requested user does not exist in the database.",
        errs.WithCause(dbErr),
    )
}

```

### Adding Field-Specific Validation Errors

If a request contains invalid input, `Validation` builds a 422 error from field errors. `FieldRequired` and `FieldInvalid` set the matching sentinel as the cause.

```go
err := errs.Validation(
    errs.FieldInvalid("email", "Must be a valid email address"),
    errs.FieldRequired("password"),
)

errors.Is(err, errs.ErrInvalid)  // true, through the email field's cause
errors.Is(err, errs.ErrRequired) // true, through the password field's cause

var fe *errs.FieldError
errors.As(err, &fe) // first field error, fe.Field() == "email"

```

For a field error with a different cause, or none, use `NewField`:

```go
errs.NewField("age", "Must be at least 18", errs.WithFieldCause(errTooYoung))
```

To attach fields to another kind of error, use the option directly:

```go
errs.BadRequest("Malformed filter.", errs.WithFields(errs.FieldInvalid("sort", "Unknown sort key")))
```

### Deriving a Variant of an Existing Error

`With` returns a modified copy and never changes the original, so it is safe to use on shared, package-level errors.

```go
var ErrUserNotFound = errs.NotFound("", errs.WithTitle("User not found"))

func getUser(id string) error {
    return ErrUserNotFound.With(errs.WithDetail("No user exists with the ID '" + id + "'."))
}

// errors.Is(getUser("123"), ErrUserNotFound) == true
```

---

## API Reference

### Core Types

* `Error`: The main struct representing a structured problem detail. Implements the standard `error` interface.
* `FieldError`: A validation failure for one field. Created with `NewField`, `FieldRequired` or `FieldInvalid`.
* `Option`: A functional option `func(*Error)` used to configure an `Error`.
* `FieldOption`: A functional option `func(*FieldError)` used to configure a `FieldError`.

### Creating Errors

| Function | Default HTTP Status | Default Type URI |
| --- | --- | --- |
| `New(detail, opts...)` | `500 Internal Server Error` | `/errors/internal` |
| `BadRequest(detail, opts...)` | `400 Bad Request` | `/errors/bad-request` |
| `Unauthorized(detail, opts...)` | `401 Unauthorized` | `/errors/unauthorized` |
| `Forbidden(detail, opts...)` | `403 Forbidden` | `/errors/forbidden` |
| `NotFound(detail, opts...)` | `404 Not Found` | `/errors/not-found` |
| `Conflict(detail, opts...)` | `409 Conflict` | `/errors/conflict` |
| `Unprocessable(detail, opts...)` | `422 Unprocessable Entity` | `/errors/unprocessable` |
| `Validation(fields...)` | `422 Unprocessable Entity` | `/errors/validation` |

The detail may be an empty string. The default title of each is the standard text of its status code, evaluated when `Title()` is called, so changing the status with `WithStatus` changes the title too.

### Deriving and Converting

* `(*Error).With(opts ...Option) *Error`: Returns a copy with the options applied. The receiver is not modified.
* `From(err error) *Error`: Returns the first `*Error` in `err`'s chain. Any other error becomes `ErrInternal` with `err` as its cause, so internal messages never reach the client. Returns `nil` for a `nil` error.
* `ErrInternal`: The generic 500 used by `From`.

### Functional Options

You can pass these options to any of the constructors above, or to `With`, to customize the error.

* `WithCause(cause error)`: Sets the underlying cause of the error (used for `Unwrap`).
* `WithType(typ string)`: Overrides the default problem type URI.
* `WithStatus(status int)`: Overrides the HTTP status code.
* `WithTitle(title string)`: Overrides the default title derived from the status.
* `WithDetail(detail string)`: Replaces the detail, mainly useful with `With`.
* `WithFields(fields ...*FieldError)`: Appends one or more field-specific validation errors. `nil` entries are skipped.

### Field Errors

* `NewField(field, detail string, opts ...FieldOption) *FieldError`
* `FieldRequired(field string)`: Detail "This field is required.", cause `ErrRequired`.
* `FieldInvalid(field, detail string)`: Cause `ErrInvalid`.
* `WithFieldDetail(detail string)`, `WithFieldCause(cause error)`: Field options.
* `ErrRequired`, `ErrInvalid`: Sentinel errors for matching a class of validation failures with `errors.Is`.

A field error's message is `field: detail`. The cause is not part of the message; it exists for `errors.Is` and `errors.As`.

### Extracting Error Data

When handling the error (e.g., in your HTTP middleware), use `From` to get an `*Error` from any error, then the getter methods to build your JSON response:

```go
e := errs.From(err)
statusCode := e.Status()    // int
typ := e.Type()             // string
title := e.Title()          // string
detail := e.Detail()        // string
fields := e.Fields()        // []*errs.FieldError, never nil

for _, f := range fields {
    _ = f.Field()           // string
    _ = f.Detail()          // string
}

```

`errors.Is` compares errors by type URI and status, so every `NotFound` matches every other `NotFound`, and every `Validation` error matches every other, unless you give them different types. `errors.As` works natively as well:

```go
var e *errs.Error
if errors.As(err, &e) {
// Write e.Status() and e.Title() to HTTP response
}

```