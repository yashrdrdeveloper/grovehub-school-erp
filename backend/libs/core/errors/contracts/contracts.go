package contracts

type Code string

func (code Code) String() string {
	return string(code)
}

type Definition struct {
	code    Code
	message string
}

func NewDefinition(code Code, message string) Definition {
	return Definition{code: code, message: message}
}

func (definition Definition) Code() Code {
	return definition.code
}

func (definition Definition) Message() string {
	return definition.message
}

type Error struct {
	definition Definition
	cause      error
	metadata   any
}

func NewError(definition Definition) *Error {
	return &Error{definition: definition}
}

func NewErrorWithMetadata(definition Definition, metadata any) *Error {
	return &Error{definition: definition, metadata: metadata}
}

func Wrap(definition Definition, cause error) *Error {
	return &Error{definition: definition, cause: cause}
}

func WrapWithMetadata(definition Definition, cause error, metadata any) *Error {
	return &Error{definition: definition, cause: cause, metadata: metadata}
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

func (err *Error) Metadata() any {
	return err.metadata
}
