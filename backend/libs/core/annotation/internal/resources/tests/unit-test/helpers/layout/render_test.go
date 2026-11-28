package layout_test

import (
	"go/ast"
	"go/token"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
)

var _=Describe("Render",func(){
	It("formats AST nodes and reports invalid node errors",func(){
		value,err:=layout.Render(token.NewFileSet(),ast.NewIdent("Record"))
		Expect(err).NotTo(HaveOccurred());Expect(value).To(Equal("Record"))
		_,err=layout.Render(token.NewFileSet(),nil)
		Expect(err).To(HaveOccurred())
	})
})
