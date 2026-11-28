package helpers_test

import (
	"path/filepath"
	"strings"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/events/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Write", func() {
	It("emits a single shared Source constant and stable ordered output",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"z.go","package events\n// @Event topic=person type=person.created.v1 source=school\ntype Created struct{}\n")
		fixture.Write(root,"a.go","package events\n// @Event topic=person type=person.changed.v2 source=school\ntype Changed struct{}\n")
		defs,pkg,err:=helpers.Discover(root)
		Expect(err).NotTo(HaveOccurred())
		out:=filepath.Join(root,"out.txt")
		Expect(helpers.Write(out,pkg,defs)).To(Succeed())
		first:=fixture.Read(out)
		Expect(strings.Count(first,"Source ")).To(Equal(1))
		Expect(strings.Index(first,"ChangedTopic")).To(BeNumerically("<",strings.Index(first,"CreatedTopic")))
		Expect(helpers.Write(out,pkg,defs)).To(Succeed())
		Expect(fixture.Read(out)).To(Equal(first))
	})
})
