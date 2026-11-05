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
	"strings"
)

// Error represents a structured problem detail that implements the standard error interface.
// It encapsulates standard HTTP error fields such as status, type, title, and detail,
// along with an optional underlying cause and field-specific validation errors.
type Error struct {
	cause  error
	typ    string
	status int
	title  string // explicit override; empty means derive from status, see Title
	detail string
	fields []*FieldError
}

var _ error = (*Error)(nil)

// Error implements the error interface. It returns a formatted string containing
// the title, detail, the underlying cause's message, and the field errors, each
// only if present.
func (e *Error) Error() string {
	var msg strings.Builder
	msg.WriteString(e.Title())
	if e.detail != "" {
		msg.WriteString(": " + e.detail)
	}
	if e.cause != nil {
		msg.WriteString(": " + e.cause.Error())
	}
	if len(e.fields) > 0 {
		parts := make([]string, len(e.fields))
		for i, f := range e.fields {
			parts[i] = "(" + f.Error() + ")"
		}
		msg.WriteString(": [" + strings.Join(parts, ", ") + "]")
	}
	return msg.String()
}

// Unwrap returns the cause followed by all field errors, so [errors.Is] and
// [errors.As] traverse every one of them.
func (e *Error) Unwrap() []error {
	out := make([]error, 0, len(e.fields)+1)
	if e.cause != nil {
		out = append(out, e.cause)
	}
	for _, f := range e.fields {
		out = append(out, f)
	}
	return out
}

// Is reports whether target is an [*Error] with the same type and status.
// [errors.Is] takes care of walking the chain, so only a direct comparison is needed.
func (e *Error) Is(target error) bool {
	var t *Error
	ok := errors.As(target, &t)
	return ok && e.typ == t.typ && e.status == t.status
}

// New returns a new [*Error] with the given detail, configured by opts. Without
// options it is a 500 Internal Server Error of type "/errors/internal". The detail
// may be empty.
//
// The title is not required: unless overridden with [WithTitle], it is derived
// from the status code, see [Error.Title].
func New(detail string, opts ...Option) *Error {
	e := &Error{
		status: http.StatusInternalServerError,
		typ:    "/errors/internal",
		detail: detail,
	}
	return e.With(opts...)
}

// With returns a copy of e with opts applied, leaving e untouched.
// It is safe to use on shared sentinels: ErrInternal.With(WithDetail("...")).
func (e *Error) With(opts ...Option) *Error {
	c := *e
	c.fields = append([]*FieldError(nil), e.fields...)
	for _, opt := range opts {
		opt(&c)
	}
	return &c
}

// From returns the first [*Error] in base's chain. If there is none, it returns
// [ErrInternal] with err as its cause, so internal messages never reach the client.
// It returns nil if err is nil.
func From(base error) *Error {
	if base == nil {
		return nil
	}
	var e *Error
	if errors.As(base, &e) {
		return e
	}
	return ErrInternal.With(WithCause(base))
}

// Status returns the HTTP status code associated with the error.
func (e *Error) Status() int {
	return e.status
}

// Title returns a short, human-readable summary of the problem type.
// It is the value set with [WithTitle] if any, otherwise the standard text of
// the status code (for example "Not Found"), otherwise "Error".
func (e *Error) Title() string {
	if e.title != "" {
		return e.title
	}
	if t := http.StatusText(e.status); t != "" {
		return t
	}
	return "Error"
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
// The returned slice is never nil.
func (e *Error) Fields() []*FieldError {
	res := make([]*FieldError, len(e.fields))
	copy(res, e.fields)
	return res
}

// Option defines a functional option for configuring an [*Error].
// It allows for flexible and extensible initialization of errors.
type Option func(*Error)

// WithCause sets the underlying cause of the error.
// This is useful for wrapping standard Go errors to preserve the original context
// for unwrapping and logging, while presenting a clean HTTP error to the client.
func WithCause(cause error) Option {
	return func(e *Error) {
		e.cause = cause
	}
}

// WithType sets the URI reference that identifies the problem type.
// According to RFC 7807 and RFC 9457, this should ideally be a URI that yields human-readable documentation.
func WithType(typ string) Option {
	return func(e *Error) {
		e.typ = typ
	}
}

// WithStatus sets the HTTP status code associated with the error.
func WithStatus(status int) Option {
	return func(e *Error) {
		e.status = status
	}
}

// WithTitle overrides the default title derived from the status code.
// It should not change from occurrence to occurrence of the problem, except for purposes of localization.
func WithTitle(title string) Option {
	return func(e *Error) {
		e.title = title
	}
}

// WithDetail sets a human-readable explanation specific to this occurrence of the problem.
func WithDetail(detail string) Option {
	return func(e *Error) {
		e.detail = detail
	}
}

// WithFields appends one or more field-specific validation errors to the problem details.
// It can be called multiple times to continue appending fields. Nil entries are skipped.
func WithFields(fields ...*FieldError) Option {
	return func(e *Error) {
		for _, f := range fields {
			if f != nil {
				e.fields = append(e.fields, f)
			}
		}
	}
}
