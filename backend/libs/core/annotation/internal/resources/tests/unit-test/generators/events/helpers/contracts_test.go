package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/events/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Discover", func() {
	It("retains an explicit type constant and version",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"event.go","package events\n// @Event topic=person type=person.created.v12 source=school typeConst=PersonCreated\ntype Created struct{}\n")
		defs,pkg,err:=helpers.Discover(root)
		Expect(err).NotTo(HaveOccurred())
		Expect(defs).To(HaveLen(1))
		Expect(pkg).To(Equal("events"))
		out:=filepath.Join(root,"out.txt")
		Expect(helpers.Write(out,pkg,defs)).To(Succeed())
		content:=fixture.Read(out)
		Expect(content).To(MatchRegexp(`PersonCreated\s*=\s*"person.created.v12"`))
		Expect(content).To(MatchRegexp(`CreatedVersion\s*=\s*12`))
	})
})
