package coreerrors

type Definition struct {
	code    Code
	message string
}

func NewDefinition(code Code, message string) Definition {
	return Definition{
		code:    code,
		message: message,
	}
}

func (definition Definition) Code() Code {
	return definition.code
}

func (definition Definition) Message() string {
	return definition.message
}
