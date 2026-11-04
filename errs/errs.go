// Package errs provides structured, HTTP-aware error handling for the broader toolkit.
//
// As a core utility within the toolkit, this package ensures consistent error
// representation across all interconnected services and APIs. It is designed to model
// problem details (inspired by RFC 7807 and RFC 9457), allowing errors to securely carry HTTP
// status codes, domain-specific types, human-readable titles, and field-level
// validation context across package boundaries.
package errs

import (
	"errors"
	"net/http"
)

// Error represents a structured problem detail that implements the standard error interface.
// It encapsulates standard HTTP error fields such as status, type, title, and detail,
// along with an optional underlying cause and field-specific validation errors.
type Error struct {
	cause error

	typ    string
	status int
	title  string
	detail string
	fields []FieldError
}

// FieldError provides additional, field-specific context for an [Error].
// This is typically used to convey validation failures for specific inputs
// or parameters in a request.
type FieldError struct {
	field  string
	detail string
}

var _ error = (*Error)(nil)

// Error implements the error interface. It returns a formatted string containing
// the title, detail, and the underlying cause's error message if present.
func (e *Error) Error() string {
	msg := e.title
	if e.detail != "" {
		msg += ": " + e.detail
	}
	if e.cause != nil {
		msg += ": " + e.cause.Error()
	}
	return msg
}

// Unwrap returns the underlying cause of the error, allowing compatibility
// with [errors.Is] and [errors.As].
func (e *Error) Unwrap() error {
	return e.cause
}

// Is reports whether any error in e's chain matches the target.
// It returns true if the target is an [*Error] with the same type and status.
func (e *Error) Is(target error) bool {
	var t *Error
	if !errors.As(target, &t) {
		return false
	}
	return e.typ == t.typ && e.status == t.status
}

// As finds the first error in e's chain that matches target, and if one is found, sets
// target to that error value and returns true. Otherwise, it returns false.
func (e *Error) As(target any) bool {
	if !errors.As(e, &target) {
		return false
	}
	return true
}

// New returns a new [*Error] with the given title and configuration options.
// If no options are provided, it defaults to a 500 Internal Server Error
// with a generic "Something went wrong" detail.
func New(title string, opts ...Option) *Error {
	e := &Error{
		status: http.StatusInternalServerError,
		typ:    "errors/internal",
		title:  title,
		detail: "Something went wrong",
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// NewFrom returns an [*Error] derived from an existing base error.
// If the base error is already an [*Error], the provided options are applied to it.
// Otherwise, it creates a new [*Error] using the base error's message as the title.
// It returns nil if the base error is nil.
func NewFrom(base error, opts ...Option) *Error {
	if base == nil {
		return nil
	}
	var target *Error
	if !errors.As(base, &target) {
		target = New(base.Error(), opts...)
	} else {
		for _, opt := range opts {
			opt(target)
		}
	}
	return target
}

// Status returns the HTTP status code associated with the error.
func (e *Error) Status() int {
	return e.status
}

// Title returns a short, human-readable summary of the problem type.
func (e *Error) Title() string {
	return e.title
}

// Type returns a URI reference that identifies the problem type.
func (e *Error) Type() string {
	return e.typ
}

// Detail returns a human-readable explanation specific to this occurrence of the problem.
func (e *Error) Detail() string {
	return e.detail
}

// Fields returns a copy of the field-specific validation errors associated with the problem.
func (e *Error) Fields() []FieldError {
	return append([]FieldError{}, e.fields...) // copy fields
}
