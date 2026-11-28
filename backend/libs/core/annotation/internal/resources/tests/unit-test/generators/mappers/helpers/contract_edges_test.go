package helpers_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/mappers/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

func mapperEdgeFixture(source,destination string)(string,string){
	root:=GinkgoT().TempDir()
	dir:=filepath.Join(root,"services","demo","mappers")
	fixture.Write(root,"services/demo/mappers/record.go","package mappers\nimport dto \"example.com/project/services/demo/dto\"\n"+source)
	fixture.Write(root,"services/demo/dto/record.go","package dto\n"+destination)
	return root,dir
}

var _=Describe("Mapper contract validation",func(){
	It("rejects invalid mapper annotations and method contracts",func(){
		for _,test:=range []struct{source,reason string}{
			{"// @Mapper\ntype Mapper interface{}","requires one identifier"},
			{"// @Mapper a\n// @Mapper b\ntype Mapper interface{}","duplicate @Mapper"},
			{"// @Mapper a-b\ntype Mapper interface{}","invalid @Mapper receiver"},
			{"// @Mapper A\ntype A interface{}\n// @Mapper B\ntype B interface{}","one @Mapper interface"},
			{"// @Mapper impl\ntype Mapper interface{ Other }","embedded mapper interfaces"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Name\nMap(name string) dto.Record }","expected @Mapping"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping bad-name=name\nMap(name string) dto.Record }","invalid destination field"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=name\n// @Mapping Name=name\nMap(name string) dto.Record }","duplicate mapping"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=\nMap(name string) dto.Record }","missing expression"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=(\nMap(name string) dto.Record }","invalid Go mapping expression"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=name\nMap(string) dto.Record }","requires one named input"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=name\nMap(name string) Record }","must be an imported DTO struct"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=name\nMap(name string) other.Record }","not a service DTO import"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Other=name\nMap(name string) dto.Record }","unmapped destination field"},
			{"// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=name\n// @Mapping Other=name\nMap(name string) dto.Record }","mapped 2 fields"},
		}{
			root,dir:=mapperEdgeFixture(test.source,"type Record struct{ Name string }")
			_,_,err:=helpers.Discover(root,"example.com/project",dir)
			Expect(err).To(MatchError(ContainSubstring(test.reason)),test.source)
		}
	})
	It("rejects invalid destinations and leaves unannotated mapper methods maintained",func(){
		source:="// @Mapper impl\ntype Mapper interface{\n// @Mapping Name=name\nMap(name string) dto.Record }"
		for _,test:=range []struct{destination,reason string}{
			{"type Record string","not a struct"},
			{"type Record struct{ Embedded }","embedded fields are forbidden"},
			{"type Other struct{}","destination struct Record not found"},
			{"func broken(","expected"},
		}{
			root,dir:=mapperEdgeFixture(source,test.destination)
			_,_,err:=helpers.Discover(root,"example.com/project",dir)
			Expect(err).To(MatchError(ContainSubstring(test.reason)))
		}
		root,dir:=mapperEdgeFixture("type Plain struct{}\ntype Other interface{}\n// @Mapper impl\ntype Mapper interface{ Maintained() }","type Record struct{ Name string }")
		fixture.Write(root,"services/demo/mappers/ignored_test.go","invalid source")
		group,found,err:=helpers.Discover(root,"example.com/project",dir)
		Expect(err).NotTo(HaveOccurred());Expect(found).To(BeTrue());Expect(helpers.HasMethods(group)).To(BeFalse())
	})
})
