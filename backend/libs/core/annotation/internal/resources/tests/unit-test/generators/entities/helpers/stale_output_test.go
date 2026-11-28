package helpers_test

import (
	"os"
	"path/filepath"
	

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/entities/helpers"
	fixture "grovehub/libs/core/annotation/internal/resources/tests/helpers"
)

var _=Describe("Accessor output lifecycle",func(){
	It("ignores generic receivers and handles empty output roots",func(){
		root:=GinkgoT().TempDir()
		fixture.Write(root,"generic.go","package entity\ntype Generic[T any] struct{}\nfunc (g *Generic[T]) GetValue() T {var value T;return value}\n")
		out:=filepath.Join(root,"missing")
		Expect(helpers.Generate(root,out)).To(Succeed())
		Expect(out).NotTo(BeAnExistingFile())
		blocked:=fixture.Write(root,"blocked","file")
		Expect(helpers.Generate(root,blocked)).NotTo(Succeed())
	})
	It("reports stale generated output that cannot be removed",func(){
		root,out:=GinkgoT().TempDir(),GinkgoT().TempDir()
		fixture.Write(root,"record.go","package entity\ntype Record struct{}")
		stale:=filepath.Join(out,"record_gombok.go")
		Expect(os.MkdirAll(filepath.Join(stale,"child"),0755)).To(Succeed())
		Expect(helpers.Generate(root,out)).NotTo(Succeed())
		Expect(stale).To(BeADirectory())
	})
	It("preserves an orphaned generated file when its directory is not writable",func(){
		if os.Geteuid()==0{Skip("requires unprivileged filesystem permissions")}
		root,out:=GinkgoT().TempDir(),GinkgoT().TempDir()
		stale:=fixture.Write(out,"removed_gombok.go","generated")
		Expect(os.Chmod(out,0555)).To(Succeed())
		defer os.Chmod(out,0755)
		Expect(helpers.Generate(root,out)).NotTo(Succeed())
		Expect(stale).To(BeAnExistingFile())
	})
})
