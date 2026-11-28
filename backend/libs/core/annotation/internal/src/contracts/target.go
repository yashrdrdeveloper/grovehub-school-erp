package contract

// Target identifies the structural source location of a directive.
type Target string

const (
	TypeTarget      Target = "type"
	InterfaceTarget Target = "interface"
	MethodTarget    Target = "method"
	FunctionTarget  Target = "function"
	FieldTarget     Target = "field"
	PackageTarget   Target = "package"
	FileTarget      Target = "file"
	ParameterTarget Target = "parameter"
)
