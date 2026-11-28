package helpers

import (
	"go/ast"
)

type entitySpec struct {
	name             string
	fields           []fieldSpec
	builder, getter, setter bool
}

type fieldSpec struct {
	name, typ string
}

type sourceFile struct {
	name, pkg string
	imports []*ast.ImportSpec
	entities []entitySpec
}

type methodIndex map[string]map[string]bool
