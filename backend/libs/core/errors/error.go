package coreerrors

type Error struct {
	definition Definition
	cause      error
}

func NewError(definition Definition) *Error {
	return &Error{definition: definition}
}

func Wrap(definition Definition, cause error) *Error {
	return &Error{
		definition: definition,
		cause:      cause,
	}
}

func (err *Error) Error() string {
	return err.definition.Message()
}

func (err *Error) Unwrap() error {
	return err.cause
}

func (err *Error) Definition() Definition {
	return err.definition
}
