package errs

import (
	"net/http"
)

// Invalid returns an [*Error] configured with a 400 Bad Request status.
// It sets a default problem type, which can be overridden via opts.
func Invalid(title string, opts ...Option) *Error {
	baseOpts := []Option{
		WithType("/errors/bad-request"),
		WithStatus(http.StatusBadRequest),
	}
	// append(baseOpts, opts...) ensures that if the user passes WithType("..."),
	// it will overwrite the default type since it is applied last.
	return New(title, append(baseOpts, opts...)...)
}

// Unauthorized returns an [*Error] configured with a 401 Unauthorized status.
func Unauthorized(title string, opts ...Option) *Error {
	baseOpts := []Option{
		WithType("/errors/unauthorized"),
		WithStatus(http.StatusUnauthorized),
	}
	return New(title, append(baseOpts, opts...)...)
}

// Forbidden returns an [*Error] configured with a 403 Forbidden status.
func Forbidden(title string, opts ...Option) *Error {
	baseOpts := []Option{
		WithType("/errors/forbidden"),
		WithStatus(http.StatusForbidden),
	}
	return New(title, append(baseOpts, opts...)...)
}

// NotFound returns an [*Error] configured with a 404 Not Found status.
func NotFound(title string, opts ...Option) *Error {
	baseOpts := []Option{
		WithType("/errors/not-found"),
		WithStatus(http.StatusNotFound),
	}
	return New(title, append(baseOpts, opts...)...)
}

// Conflict returns an [*Error] configured with a 409 Conflict status.
func Conflict(title string, opts ...Option) *Error {
	baseOpts := []Option{
		WithType("/errors/conflict"),
		WithStatus(http.StatusConflict),
	}
	return New(title, append(baseOpts, opts...)...)
}

// Unprocessable returns an [*Error] configured with a 422 Unprocessable Entity status.
func Unprocessable(title string, opts ...Option) *Error {
	baseOpts := []Option{
		WithType("/errors/unprocessable"),
		WithStatus(http.StatusUnprocessableEntity),
	}
	return New(title, append(baseOpts, opts...)...)
}

// Internal returns an [*Error] configured with a 500 Internal Server Error status.
func Internal(title string, opts ...Option) *Error {
	baseOpts := []Option{
		WithType("/errors/internal"),
		WithStatus(http.StatusInternalServerError),
	}
	return New(title, append(baseOpts, opts...)...)
}

// ErrInternal is now cleaner to define.
var ErrInternal = Internal("Internal server error", WithDetail("An unexpected error occurred."))
