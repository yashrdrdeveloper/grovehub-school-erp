package coverage_test

import (
	"os"
	"path/filepath"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/testing/coverage"
)

// A subprocess fixture stands in only for the child Go command. It tests
// runner orchestration and error propagation, not the actual CI coverage.
func commandFixture(root, script string) string {
	bin:=filepath.Join(root,"bin")
	path:=filepath.Join(bin,"go")
	writeFixture(path,"#!/bin/sh\nset -eu\n"+script)
	Expect(os.Chmod(path,0755)).To(Succeed())
	GinkgoT().Setenv("PATH",bin+string(os.PathListSeparator)+os.Getenv("PATH"))
	return path
}

var _=Describe("Coverage subprocess boundary",func(){
	It("runs inventory for the selected owner and reports child failure",func(){
		root,_,_,args:=coverageFixture()
		log:=filepath.Join(root,"args.log")
		GinkgoT().Setenv("COVERAGE_FIXTURE_LOG",log)
		commandFixture(root,"printf '%s\\n' \"$@\" > \"$COVERAGE_FIXTURE_LOG\"\nexit 0\n")
		inventory:=append(args,"-mode","--inventory")
		Expect(coverage.Run(inventory)).To(Succeed())
		data,err:=os.ReadFile(log)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(ContainSubstring("inventory\n"+filepath.Join(root,"services","sample")+"\n"))
		commandFixture(root,"exit 17\n")
		Expect(coverage.Run(inventory)).To(MatchError(ContainSubstring("exit status 17")))
	})
	It("rejects failed and empty package discovery",func(){
		root,_,_,args:=coverageFixture()
		commandFixture(root,"exit 17\n")
		Expect(coverage.Run(args)).To(MatchError("one or more selected unit coverage runs failed"))
		commandFixture(root,"printf '%s\\n' example.com/coveragefixture/tests/unit-test example.com/coveragefixture/testdata\n")
		Expect(coverage.Run(args)).To(MatchError("one or more selected unit coverage runs failed"))
	})
	It("rejects a child that fails to create its requested profile",func(){
		root,_,_,args:=coverageFixture()
		commandFixture(root,"case \"$1\" in\nlist) echo example.com/coveragefixture/services/sample ;;\ntest) exit 0 ;;\n*) exit 1 ;;\nesac\n")
		Expect(coverage.Run(args)).To(MatchError("one or more selected unit coverage runs failed"))
	})
	It("keeps valid coverage results when optional HTML rendering fails",func(){
		root,_,_,args:=coverageFixture()
		commandFixture(root,"case \"$1\" in\nlist) echo example.com/coveragefixture/services/sample ;;\ntest) for arg in \"$@\"; do case \"$arg\" in -coverprofile=*) printf 'mode: atomic\\nfixture.go:1.1,1.2 1 1\\n' > \"${arg#-coverprofile=}\" ;; esac; done ;;\ntool) exit 19 ;;\n*) exit 1 ;;\nesac\n")
		Expect(coverage.Run(args)).To(Succeed())
		profile:=filepath.Join(root,".generated","unit-coverage","sample.out")
		contents,err:=os.ReadFile(profile)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(contents)).To(ContainSubstring("fixture.go:1.1,1.2 1 1"))
	})
	It("refuses to replace a dangling destination and leaves source intact",func(){
		root:=GinkgoT().TempDir()
		overlay,manifest:=filepath.Join(root,"overlay.json"),filepath.Join(root,"manifest")
		source:=filepath.Join(root,"source.go")
		writeFixture(source,"package example")
		destination:=filepath.Join(root,"zz_generated.go")
		Expect(os.Symlink(filepath.Join(root,"missing.go"),destination)).To(Succeed())
		writeOverlay(overlay,map[string]string{destination:source})
		_,err:=coverage.Materialize(overlay,manifest)
		Expect(err).To(HaveOccurred())
		_,err=os.Lstat(destination)
		Expect(err).NotTo(HaveOccurred())
		Expect(manifest).NotTo(BeAnExistingFile())
		data,err:=os.ReadFile(source);Expect(err).NotTo(HaveOccurred());Expect(string(data)).To(Equal("package example"))
	})
	It("preserves a generated file when its parent cannot be modified",func(){
		if os.Geteuid()==0 { Skip("permission error requires an unprivileged process") }
		root:=GinkgoT().TempDir()
		dir:=filepath.Join(root,"readonly")
		file:=filepath.Join(dir,"zz_generated.go")
		writeFixture(file,"// Code generated.\npackage example")
		manifest:=filepath.Join(root,"manifest")
		writeFixture(manifest,file+"\n")
		Expect(os.Chmod(dir,0555)).To(Succeed())
		defer os.Chmod(dir,0755)
		Expect(coverage.RecoverGenerated(manifest,root)).NotTo(Succeed())
		Expect(file).To(BeAnExistingFile())
		Expect(manifest).To(BeAnExistingFile())
	})
	It("rejects invalid current-directory resolution for both roots",func(){
		if os.PathSeparator!='/' { Skip("Unix deleted-directory semantics") }
		original,err:=os.Getwd();Expect(err).NotTo(HaveOccurred())
		root:=GinkgoT().TempDir()
		deleted:=filepath.Join(root,"deleted")
		Expect(os.Mkdir(deleted,0755)).To(Succeed())
		Expect(os.Chdir(deleted)).To(Succeed())
		defer func(){Expect(os.Chdir(original)).To(Succeed())}()
		Expect(os.Remove(deleted)).To(Succeed())
		Expect(coverage.Run([]string{"-backend","."})).NotTo(Succeed())
		Expect(coverage.Run([]string{"-backend",root,"-core",".","-generator-script","unused","-reports",root})).NotTo(Succeed())
	})
	It("exposes the exact list of selected capability test directories",func(){
		root,generator,reports,_:=coverageFixture()
		writeFixture(filepath.Join(root,"README.md"),"not a package")
		writeFixture(filepath.Join(root,"utils","utils.go"),"package utils\nfunc Answer()int{return 42}")
		writeFixture(filepath.Join(root,"utils","internal","resources","tests","unit-test","utils_test.go"),"package utils_test\nimport (\"testing\"; \"example.com/coveragefixture/utils\")\nfunc TestAnswer(t *testing.T){if utils.Answer()!=42{t.Fatal(\"bad answer\")}}")
		log:=filepath.Join(root,"child-args.log")
		GinkgoT().Setenv("COVERAGE_FIXTURE_LOG",log)
		commandFixture(root,"case \"$1\" in\nlist) printf '%s\\n' example.com/coveragefixture/utils ;;\ntest) printf '%s\\n' \"$@\" > \"$COVERAGE_FIXTURE_LOG\"; for arg in \"$@\"; do case \"$arg\" in -coverprofile=*) printf 'mode: atomic\\nfixture.go:1.1,1.2 1 1\\n' > \"${arg#-coverprofile=}\" ;; esac; done ;;\ntool) exit 0 ;;\nesac\n")
		Expect(coverage.Run([]string{"-backend",root,"-core",root,"-generator-script",generator,"-reports",reports,"-package","core"})).To(Succeed())
		data,err:=os.ReadFile(log);Expect(err).NotTo(HaveOccurred())
		Expect(strings.Split(string(data),"\n")).To(ContainElement("./utils/internal/resources/tests/unit-test/..."))
	})
})
