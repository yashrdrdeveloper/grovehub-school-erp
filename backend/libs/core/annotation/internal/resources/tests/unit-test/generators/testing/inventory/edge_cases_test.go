package inventory_test

import (
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/testing/inventory"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Source inventory",func(){
	It("counts methods and missing mirrors deterministically in an application",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"internal/app/z.go","package app\nfunc Last(){}")
		fixture.Write(root,"internal/app/a.go","package app\ntype Record struct{}\nfunc (Record) Z(){}\nfunc (Record) A(){}\nfunc First(){}")
		fixture.Write(root,"internal/app/ignored_test.go","invalid source")
		fixture.Write(root,"internal/app/README.md","ignored")
		fixture.Write(root,"internal/tests/unit-test/z_test.go","package app")
		report,err:=inventory.Inspect(root)
		Expect(err).NotTo(HaveOccurred())
		Expect(report.SourceFiles).To(Equal(2));Expect(report.TestFiles).To(Equal(1))
		Expect(report.Functions).To(Equal(4));Expect(report.MirroredFunctions).To(Equal(1))
		Expect(report.Missing).To(HaveLen(3))
		Expect(report.Missing[0].Name).To(Equal("(receiver).A"))
		Expect(report.Missing[1].Name).To(Equal("(receiver).Z"))
		Expect(report.Missing[2].Name).To(Equal("First"))
	})
	It("rejects malformed Go source and missing owners",func(){
		root:=GinkgoT().TempDir()
		_,err:=inventory.Inspect(filepath.Join(root,"missing"))
		Expect(err).To(HaveOccurred())
		fixture.Write(root,"broken.go","package core\nfunc broken(")
		_,err=inventory.Inspect(root)
		Expect(err).To(HaveOccurred())
	})
})
