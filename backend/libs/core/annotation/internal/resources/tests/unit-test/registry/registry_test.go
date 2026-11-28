package registry_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/annotation/internal/src/registry"
)

var _ = Describe("ParseComment", func() {
	DescribeTable("rejects invalid directive syntax",func(line string){_,found,err:=registry.ParseComment(line);Expect(found).To(BeTrue());Expect(err).To(HaveOccurred())},
		Entry("empty name","// @"),Entry("invalid name","// @Bad-Name"),Entry("unknown name","// @UnknownDirective"),
		Entry("entity missing schema","// @Entity persons"),Entry("entity bad option","// @Entity demo.persons mode=wrong"),
		Entry("repository option","// @Repository entity=1bad"),Entry("search duplicate","// @Searchable config=cfg alias=p alias=q"),
		Entry("search config","// @Searchable config=a.b.c alias=p"),Entry("filter duplicates","// @Filterable fields=gid,gid"),
		Entry("cursor invalid","// @Cursor field=.gid"),
	)
	It("accepts valid search metadata and ignores regular comments",func(){
		d,found,err:=registry.ParseComment("// @Searchable config=constants.Config alias=p impl=Repo entity=Person")
		Expect(err).NotTo(HaveOccurred());Expect(found).To(BeTrue());Expect(d.Name).To(Equal("Searchable"))
		_,found,err=registry.ParseComment("// regular comment")
		Expect(err).NotTo(HaveOccurred());Expect(found).To(BeFalse())
	})
})
var _ = Describe("ValidateDirective", func() {
	It("rejects multiline arguments",func(){Expect(registry.ValidateDirective(registry.Directive{Name:"Entity",Arguments:"demo.persons\n"})).To(HaveOccurred())})
})
