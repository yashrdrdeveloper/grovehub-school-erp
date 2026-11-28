package helpers_test

import (
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/entities/helpers"
	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
)

var _ = Describe("Generate", func() {
	It("retains qualified types in generated accessor imports", func() {
		root := GinkgoT().TempDir()
		out := GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\nimport clock \"time\"\n// @Getter\ntype Record struct{ Created clock.Time }\n")
		Expect(helpers.Generate(root,out)).To(Succeed())
		content:=fixture.Read(filepath.Join(out,"record_gombok.go"))
		Expect(content).To(ContainSubstring(`clock "time"`))
		Expect(content).To(ContainSubstring("GetCreated() clock.Time"))
	})
})
