package coreerrors

type Code string

func (code Code) String() string {
	return string(code)
}
