package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/patchdto/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Patch contract validation",func(){
	It("rejects malformed declarations and ambiguous fields",func(){
		for _,test:=range []struct{source,reason string}{
			{"// @PatchRequestDto argument\ntype Update struct{}","does not take arguments"},
			{"// @PatchRequestDto\n// @PatchRequestDto\ntype Update struct{}","duplicate @PatchRequestDto"},
			{"// @PatchRequestDto\ntype Update[T any] struct{}","requires a non-generic struct"},
			{"// @PatchRequestDto\ntype Update string","requires a non-generic struct"},
			{"// @PatchRequestDto\ntype Update struct{ private annotation.Patch[string] }","fields must be exported"},
			{"// @PatchRequestDto\ntype Update struct{ A,B annotation.Patch[string] }","fields must be exported"},
			{"// @PatchRequestDto\ntype Update struct{ Embedded }","fields must be exported"},
			{"// @PatchRequestDto\ntype Update struct{ Name Other[string] }","must use"},
			{"// @PatchRequestDto\ntype Update struct{ Name annotation.Other[string] }","must use"},
			{"// @PatchRequestDto\ntype Update struct{ Name annotation.Patch[string] }","missing json tag"},
			{"// @PatchRequestDto\ntype Update struct{ Name annotation.Patch[string] `json:\"-\"` }","invalid or duplicate json field"},
			{"// @PatchRequestDto\ntype Update struct{ A annotation.Patch[string] `json:\"x\"`; B annotation.Patch[int] `json:\"x\"` }","invalid or duplicate json field"},
			{"// @PatchRequestDto\ntype Update struct{}\n// @PatchRequestDto\ntype Update struct{}","duplicate @PatchRequestDto"},
			{"// @PatchRequestDto\ntype Update struct{}\nfunc (u *Update) MarshalJSON()([]byte,error){return nil,nil}","remove handwritten MarshalJSON"},
		}{
			root:=GinkgoT().TempDir()
			fixture.Write(root,"update.go","package dto\nimport \"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation\"\n"+test.source)
			_,_,err:=helpers.Discover(root)
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.source)
		}
	})
	It("ignores generated files and detects mixed or malformed packages",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"ignored_test.go","not Go source")
		fixture.Write(root,"ignored.gen.go","not Go source")
		fixture.Write(root,"update.go","package dto\nconst x=1\ntype Plain struct{}\nfunc Helper(){}")
		specs,pkg,err:=helpers.Discover(root)
		Expect(err).NotTo(HaveOccurred());Expect(specs).To(BeEmpty());Expect(pkg).To(Equal("dto"))
		fixture.Write(root,"other.go","package other")
		_,_,err=helpers.Discover(root)
		Expect(err).To(MatchError(ContainSubstring("mixed packages")))
		fixture.Write(root,"other.go","package dto\nfunc broken(")
		_,_,err=helpers.Discover(root)
		Expect(err).To(HaveOccurred())
	})
})
