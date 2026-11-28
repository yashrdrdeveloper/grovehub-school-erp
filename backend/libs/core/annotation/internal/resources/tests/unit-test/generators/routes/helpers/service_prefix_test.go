package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/routes/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("ServicePrefix",func(){
	It("rejects wildcard service prefixes and treats missing files as unconfigured",func(){
		root:=GinkgoT().TempDir()
		_,configured,err:=helpers.ServicePrefix(filepath.Join(root,"missing.go"))
		Expect(err).NotTo(HaveOccurred())
		Expect(configured).To(BeFalse())
		path:=fixture.Write(root,"http.go","package routes\n// @Routes /school/{id}\nfunc RegisterHTTP(){}\n")
		_,_,err=helpers.ServicePrefix(path)
		Expect(err).To(MatchError(ContainSubstring("cannot contain path variables")))
	})
})
