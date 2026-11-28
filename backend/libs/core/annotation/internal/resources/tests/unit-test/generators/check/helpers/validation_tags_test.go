package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"grovehub/libs/core/annotation/internal/src/generators/check/helpers"
)

var _ = Describe("CheckValidationTags",func(){
	DescribeTable("requires an executable validation tag",func(comment,tag string,valid bool){
		path:=filepath.Join(GinkgoT().TempDir(),"source.go")
		source:="package demo\ntype Input struct {\n// "+comment+"\nName string "+tag+"\n}\n"
		Expect(os.WriteFile(path,[]byte(source),0600)).To(Succeed())
		err:=helpers.CheckValidationTags(path)
		if valid{Expect(err).NotTo(HaveOccurred());return}
		Expect(err).To(HaveOccurred())
	},Entry("size valid","@Size max=4","`validate:\"maxrunes=4\"`",true),Entry("size absent","@Size max=4","",false),Entry("size negative","@Size max=-1","",false),Entry("nonnull","@NonNull","`validate:\"nonnull\"`",true))
})
