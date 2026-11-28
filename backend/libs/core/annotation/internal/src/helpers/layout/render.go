package layout

import (
	"bytes"
	"go/ast"
	"go/format"
	"go/token"
)

// Renderer is the Go syntax formatting boundary for annotation generators.
type Renderer func(*token.FileSet,ast.Node)(string,error)

func Render(set *token.FileSet,node ast.Node)(string,error){
	var buffer bytes.Buffer
	err:=format.Node(&buffer,set,node)
	return buffer.String(),err
}
