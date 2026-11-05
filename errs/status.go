package errs

import (
	"net/http"
)

// withKind builds an [*Error] of the given status and type slug. The defaults come
// first so that options passed by the caller, such as WithType, override them.
func withKind(status int, slug, detail string, opts []Option) *Error {
	base := []Option{WithStatus(status), WithType("/errors/" + slug)}
	return New(detail, append(base, opts...)...)
}

// BadRequest returns an [*Error] configured with a 400 Bad Request status.
// It sets a default problem type, which can be overridden via opts.
func BadRequest(detail string, opts ...Option) *Error {
	return withKind(http.StatusBadRequest, "bad-request", detail, opts)
}

// Unauthorized returns an [*Error] configured with a 401 Unauthorized status.
func Unauthorized(detail string, opts ...Option) *Error {
	return withKind(http.StatusUnauthorized, "unauthorized", detail, opts)
}

// Forbidden returns an [*Error] configured with a 403 Forbidden status.
func Forbidden(detail string, opts ...Option) *Error {
	return withKind(http.StatusForbidden, "forbidden", detail, opts)
}

// NotFound returns an [*Error] configured with a 404 Not Found status.
func NotFound(detail string, opts ...Option) *Error {
	return withKind(http.StatusNotFound, "not-found", detail, opts)
}

// Conflict returns an [*Error] configured with a 409 Conflict status.
func Conflict(detail string, opts ...Option) *Error {
	return withKind(http.StatusConflict, "conflict", detail, opts)
}

// Unprocessable returns an [*Error] configured with a 422 Unprocessable Entity status.
func Unprocessable(detail string, opts ...Option) *Error {
	return withKind(http.StatusUnprocessableEntity, "unprocessable", detail, opts)
}

// Validation returns a 422 [*Error] of type "/errors/validation" carrying the given
// field errors. Derive variants with [Error.With], for example to change the detail.
func Validation(fields ...*FieldError) *Error {
	return Unprocessable("The request contains invalid fields.",
		WithType("/errors/validation"),
		WithFields(fields...),
	)
}

// ErrInternal is the generic 500 used by [From] for errors that are not an [*Error].
// Derive variants with [Error.With] rather than modifying it.
var ErrInternal = New("An unexpected error occurred.")
