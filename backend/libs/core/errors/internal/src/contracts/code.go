package contracts

// Code is the stable machine-readable application error identifier.
type Code string

func (code Code) String() string {
	return string(code)
}
