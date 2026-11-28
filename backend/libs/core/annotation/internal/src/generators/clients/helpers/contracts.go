package helpers

import (
	"go/ast"
)

type endpointMethod struct {
	verb, path, resultType string
	created bool
}

type proxyMethod struct {
	name, signature, delegate, arguments string
	hasResults bool
	endpoint *endpointMethod
}

type clientProxy struct {
	receiver string
	imports []*ast.ImportSpec
	methods []proxyMethod
}
