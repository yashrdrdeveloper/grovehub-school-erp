package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go/parser"
	"go/token"
	"grovehub/libs/core/annotation/internal/src/parser/helpers"
	contract "grovehub/libs/core/annotation/internal/src/contracts"
)

var _ = Describe("CollectTargets", func() {
	It("distinguishes package, type, interface, method, field and function docs", func() {
		source:="// package docs\npackage demo\n// type docs\ntype Record struct {\n// field docs\nName string\n}\n// interface docs\ntype Service interface {\n// method docs\nRun()\n}\n// function docs\nfunc Run(){}\n// receiver docs\nfunc (Record) Run(){}\n"
		file,err:=parser.ParseFile(token.NewFileSet(),"demo.go",source,parser.ParseComments)
		Expect(err).NotTo(HaveOccurred())
		targets:=helpers.CollectTargets(file)
		var result []contract.Target
		for _,group:=range file.Comments {result=append(result,targets.For(group))}
		Expect(result).To(Equal([]contract.Target{contract.PackageTarget,contract.TypeTarget,contract.FieldTarget,contract.InterfaceTarget,contract.MethodTarget,contract.FunctionTarget,contract.MethodTarget}))
		Expect(targets.For(nil)).To(Equal(contract.FileTarget))
	})
})
