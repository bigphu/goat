# `errs`

`errs` is a structured, HTTP-aware error package for Go. It is designed to represent problem details for HTTP APIs 
(inspired by [RFC 7807](https://datatracker.ietf.org/doc/html/rfc7807) and [RFC 9457](https://datatracker.ietf.org/doc/html/rfc9457/)), 
allowing errors to carry statuses, types, titles, details, and field-specific validation contexts.

By utilizing the functional options pattern, `errs` provides a clean, unified, and extensible API for building complex 
API errors while maintaining full compatibility with the standard Go `errors` package (`errors.Is`, `errors.As`, `Unwrap`).

## Features

* **HTTP-Aware:** Built-in helper functions for common HTTP statuses (400, 401, 404, etc.).
* **Functional Options:** Clean and flexible error construction without bloated parameter lists.
* **Error Wrapping:** Seamlessly wrap underlying causes to preserve logs while presenting clean errors to users.
* **Validation Context:** Support for attaching field-specific errors (`FieldError`).
* **Standard Library Compatible:** Implements `Error()`, `Unwrap()`, `Is()`, and `As()`.

---

## Quick Start

### Creating a Basic HTTP Error

Use the HTTP helper functions to quickly generate contextual errors. Only the `title` is required; everything else can be configured via options.

```go
package main

import (
	"fmt"
	"github.com/bigphu/goat/errs"
)

func main() {
	err := errs.NotFound(
		"User not found",
		errs.WithDetail("No user exists with the ID '123'."),
	)
	
	fmt.Println(err.Error()) 
	// Output: User not found: No user exists with the ID '123'.
}

```

### Wrapping an Existing Error

You can wrap underlying database or system errors to preserve the original stack or context, while safely controlling what gets exposed in the HTTP response.

```go
func getUser(id string) error {
    dbErr := sql.ErrNoRows // Simulated database error
    
    return errs.NotFound(
        "User not found",
        errs.WithCause(dbErr),
        errs.WithDetail("The requested user does not exist in the database."),
    )
}

```

### Adding Field-Specific Validation Errors

If a request contains invalid input, you can attach field-level details to the error.

```go
err := errs.Invalid(
    "Validation Failed",
    errs.WithDetail("The provided request payload contains invalid fields."),
    errs.WithFields(
        errs.FieldError{
            Field: "email", 
            Detail: "Must be a valid email address",
        },
        errs.FieldError{
            Field: "password", 
            Detail: "Must be at least 8 characters long",
        },
    ), 
)

```

---

## API Reference

### Core Types

* `Error`: The main struct representing a structured problem detail. Implements the standard `error` interface.
* `FieldError`: A struct `(field string, detail string)` used to convey validation failures for specific inputs.
* `Option`: A functional option `func(*Error)` used to configure an `Error`.

### Creating Errors

| Function | Default HTTP Status | Default Type URI |
| --- | --- | --- |
| `New(title, opts...)` | `500 Internal Server Error` | `/errors/internal` |
| `NewFrom(base, opts...)` | _Inherits from base or 500_ | _Inherits or_ `/errors/internal` |
| `Invalid(title, opts...)` | `400 Bad Request` | `/errors/bad-request` |
| `Unauthorized(title, opts...)` | `401 Unauthorized` | `/errors/unauthorized` |
| `Forbidden(title, opts...)` | `403 Forbidden` | `/errors/forbidden` |
| `NotFound(title, opts...)` | `404 Not Found` | `/errors/not-found` |
| `Conflict(title, opts...)` | `409 Conflict` | `/errors/conflict` |
| `Unprocessable(title, opts...)` | `422 Unprocessable Entity` | `/errors/unprocessable` |
| `Internal(title, opts...)` | `500 Internal Server Error` | `/errors/internal` |

### Functional Options

You can pass these options to any of the constructors above to customize the error.

* `WithCause(cause error)`: Sets the underlying cause of the error (used for `Unwrap`).
* `WithType(typ string)`: Overrides the default problem type URI.
* `WithStatus(status int)`: Overrides the HTTP status code.
* `WithTitle(title string)`: Overrides the problem title.
* `WithDetail(detail string)`: Sets a human-readable explanation specific to this occurrence.
* `WithFields(fields ...FieldError)`: Appends one or more field-specific validation errors.

### Extracting Error Data

When handling the error (e.g., in your HTTP middleware), you can use the getter methods to build your JSON response:

```go
if target, ok := err.(*errs.Error); ok {
    statusCode := target.Status()    // int
    typ := target.Type()             // string
    title := target.Title()          // string
    detail := target.Detail()        // string
    fields := target.Fields()        // []errs.FieldError
}

```

Because `errs.Error` implements `Unwrap()`, you can also use `errors.As` natively:

```go
var e *errs.Error
if errors.As(err, &e) {
    // Write e.Status() and e.Title() to HTTP response
}

```