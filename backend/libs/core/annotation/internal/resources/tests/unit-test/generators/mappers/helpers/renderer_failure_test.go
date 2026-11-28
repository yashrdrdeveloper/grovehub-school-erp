package helpers_test

import (
	"errors"
	"go/ast"
	"go/token"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/mappers/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/layout"
)

var _=Describe("Mapper syntax formatting failures",func(){
	It("propagates signature and output type formatting failures",func(){
		root,dir:=mapperEdgeFixture("// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=input\nMap(input string) dto.Record }","type Record struct{Name string}")
		failure:=errors.New("renderer failed")
		for _,failAt:=range []int{1,2}{
			calls:=0
			_,found,err:=helpers.DiscoverWithRenderer(root,"example.com/project",dir,func(set *token.FileSet,node ast.Node)(string,error){
				calls++
				if calls==failAt{return "",failure}
				return layout.Render(set,node)
			})
			Expect(err).To(MatchError(failure));Expect(found).To(BeFalse())
		}
		Expect(filepath.Join(root,".generated")).NotTo(BeAnExistingFile())
	})
})
