package contracts

// Error carries an application definition, an optional cause and metadata.
type Error struct {
	definition Definition
	cause      error
	metadata   any
}

// NewError is the common construction path for creation and wrapping.
func NewError(definition Definition, cause error, metadata any) *Error {
	return &Error{definition: definition, cause: cause, metadata: metadata}
}

// Error returns only the declared message, never the underlying cause.
func (err *Error) Error() string {
	return err.definition.Message()
}

func (err *Error) Unwrap() error {
	return err.cause
}

func (err *Error) Definition() Definition {
	return err.definition
}

func (err *Error) Metadata() any {
	return err.metadata
}
