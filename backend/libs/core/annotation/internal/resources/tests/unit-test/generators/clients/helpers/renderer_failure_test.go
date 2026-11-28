package helpers_test

import (
	"errors"
	"go/ast"
	"go/token"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/clients/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Client syntax formatting failures",func(){
	It("propagates signature and response formatter failures",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"source.go","package client\n// @Client impl\ntype Client interface {\n// @RequestLine GET /records\nRead(ctx context.Context, req Request)(Result,error)\n}")
		routes:=fixture.Write(root,"routes/routes.go","package routes\nfunc register(){mux.Handle(\"GET /records\",nil)}")
		failure:=errors.New("renderer failed")
		for _,failAt:=range []int{1,2}{
			calls:=0
			proxies,err:=helpers.DiscoverWithRenderer(root,routes,func(set *token.FileSet,node ast.Node)(string,error){
				calls++
				if calls==failAt{return "",failure}
				return layout.Render(set,node)
			})
			Expect(err).To(MatchError(ContainSubstring("renderer failed")))
			Expect(proxies).To(BeNil())
		}
		Expect(filepath.Join(root,"generated.go")).NotTo(BeAnExistingFile())
	})
})
