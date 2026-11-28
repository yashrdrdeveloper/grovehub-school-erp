package helpers

import (
	"go/ast"
)

type mappingMethod struct {
	name, signature, output string
	fields []fieldMapping
}

type fieldMapping struct {name, expression string}

type mappingGroup struct {
	receiver string
	imports []*ast.ImportSpec
	methods []mappingMethod
}

// HasMethods reports whether discovery found executable mapping methods.
func HasMethods(group mappingGroup) bool {
	return len(group.methods) > 0
}
