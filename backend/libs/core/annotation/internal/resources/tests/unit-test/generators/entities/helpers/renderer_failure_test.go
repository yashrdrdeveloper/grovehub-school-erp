package helpers_test

import (
	"errors"
	"go/ast"
	"go/token"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/entities/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Accessor syntax formatting failures",func(){
	It("reports failed field formatting without writing a partial accessor",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"source.go","package entity\n// @Builder\ntype Record struct{ Name string; Other string }")
		out:=filepath.Join(root,"generated")
		failure:=errors.New("renderer failed")
		calls:=0
		Expect(helpers.GenerateWithRenderer(root,out,func(set *token.FileSet,node ast.Node)(string,error){
			calls++
			if calls==2{return "",failure}
			return layout.Render(set,node)
		})).To(MatchError(failure))
		Expect(out).NotTo(BeAnExistingFile())
	})
})
