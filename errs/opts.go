package errs

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

// WithTitle sets a short, human-readable summary of the problem type.
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
// It can be called multiple times to continue appending fields.
func WithFields(fields ...FieldError) Option {
	return func(e *Error) {
		e.fields = append(e.fields, fields...)
	}
}
