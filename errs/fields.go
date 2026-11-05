package errs

import (
	"errors"
)

// Sentinel errors for use as the cause of a [FieldError], so callers can match
// a whole class of validation failures with [errors.Is].
var (
	ErrRequired = errors.New("required")
	ErrInvalid  = errors.New("invalid")
)

// FieldError provides additional, field-specific context for an [Error].
// This is typically used to convey validation failures for specific inputs
// or parameters in a request.
type FieldError struct {
	cause  error
	field  string
	detail string
}

var _ error = (*FieldError)(nil)

// Error implements the error interface: the field, then the detail if present.
// The cause is deliberately left out of the message; it exists for [errors.Is]
// and [errors.As] through [FieldError.Unwrap].
func (e *FieldError) Error() string {
	if e.detail == "" {
		return e.field
	}
	return e.field + ": " + e.detail
}

// NewField returns a new [*FieldError] for the given field name and detail,
// configured by opts.
func NewField(field, detail string, opts ...FieldOption) *FieldError {
	e := &FieldError{
		field:  field,
		detail: detail,
	}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// FieldRequired returns a [*FieldError] for a missing field, with [ErrRequired] as its cause.
func FieldRequired(field string) *FieldError {
	return NewField(field, "This field is required.", WithFieldCause(ErrRequired))
}

// FieldInvalid returns a [*FieldError] for a field with a bad value, with [ErrInvalid] as its cause.
func FieldInvalid(field, detail string) *FieldError {
	return NewField(field, detail, WithFieldCause(ErrInvalid))
}

// Unwrap returns the underlying cause, for [errors.Is] and [errors.As].
func (e *FieldError) Unwrap() error {
	return e.cause
}

// Is reports whether target is a [*FieldError] for the same field.
func (e *FieldError) Is(target error) bool {
	var t *FieldError
	ok := errors.As(target, &t)
	return ok && e.field == t.field
}

// Field returns the name of the field the problem applies to.
func (e *FieldError) Field() string {
	return e.field
}

// Detail returns a human-readable explanation specific to this field.
func (e *FieldError) Detail() string {
	return e.detail
}

// FieldOption defines a functional option for configuring a [*FieldError].
type FieldOption func(*FieldError)

// WithFieldDetail sets a human-readable explanation specific to this field.
func WithFieldDetail(detail string) FieldOption {
	return func(e *FieldError) { e.detail = detail }
}

// WithFieldCause sets the underlying cause of the field error, typically
// [ErrRequired] or [ErrInvalid].
func WithFieldCause(cause error) FieldOption {
	return func(e *FieldError) { e.cause = cause }
}
