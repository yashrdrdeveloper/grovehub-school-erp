package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/search/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Write",func(){
	It("uses the facade transaction resolver and deterministic output",func(){
		root:=searchFixture(searchTestImpl,searchTestConfig)
		specs,err:=helpers.Discover(root,"grovehub","demo")
		Expect(err).NotTo(HaveOccurred())
		out:=filepath.Join(root,"output.txt")
		Expect(helpers.Write(out,"grovehub",specs)).To(Succeed())
		content:=fixture.Read(out)
		Expect(content).To(ContainSubstring("coredatabase.TransactionDB(ctx, repository.db)"))
		Expect(content).To(ContainSubstring(`coredatabase "grovehub/libs/core/database"`))
		Expect(helpers.Write(out,"grovehub",specs)).To(Succeed())
		Expect(fixture.Read(out)).To(Equal(content))
	})
})
