package coverage_test

import (
	"encoding/json"
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/testing/coverage"
)

func writeFixture(path, contents string) {
	Expect(os.MkdirAll(filepath.Dir(path),0755)).To(Succeed())
	Expect(os.WriteFile(path,[]byte(contents),0644)).To(Succeed())
}

func writeOverlay(path string, replace map[string]string) {
	data,err:=json.Marshal(coverage.OverlayConfig{Replace:replace})
	Expect(err).NotTo(HaveOccurred())
	writeFixture(path,string(data))
}

var _=Describe("Coverage profile validation",func(){
	It("rejects malformed, negative and conflicting coverage records",func(){
		path:=filepath.Join(GinkgoT().TempDir(),"coverage.out")
		_,_,err:=coverage.Score(path)
		Expect(err).To(HaveOccurred())
		for _,profile:=range []string{
			"",
			"unexpected mode\n",
			"mode: atomic\ninvalid\n",
			"mode: atomic\nfile.go:1.1,2.2 nope 0\n",
			"mode: atomic\nfile.go:1.1,2.2 1 nope\n",
			"mode: atomic\nfile.go:1.1,2.2 -1 0\n",
			"mode: atomic\nfile.go:1.1,2.2 1 -1\n",
			"mode: atomic\nfile.go:1.1,2.2 1 0\nfile.go:1.1,2.2 2 1\n",
		}{
			writeFixture(path,profile)
			_,_,err:=coverage.Score(path)
			Expect(err).To(HaveOccurred(),profile)
		}
	})
	It("counts exactly the statements in a completely covered profile",func(){
		path:=filepath.Join(GinkgoT().TempDir(),"coverage.out")
		writeFixture(path,"mode: atomic\nfile.go:1.1,2.2 2 1\nfile.go:3.1,4.2 3 2\n")
		total,covered,err:=coverage.Score(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(total).To(Equal(int64(5)))
		Expect(covered).To(Equal(total))
	})
})

var _=Describe("Generated source recovery edge cases",func(){
	It("handles missing manifests and missing stale files without deleting maintained sources",func(){
		root:=GinkgoT().TempDir()
		manifest:=filepath.Join(root,"manifest")
		Expect(coverage.RecoverGenerated(manifest,root)).To(Succeed())
		writeFixture(manifest,"\n"+filepath.Join(root,"zz_missing.go")+"\n")
		Expect(coverage.RecoverGenerated(manifest,root)).To(Succeed())
		outside:=filepath.Join(GinkgoT().TempDir(),"zz_outside.go")
		writeFixture(outside,"// Code generated.\npackage example\n")
		writeFixture(manifest,outside+"\n")
		Expect(coverage.RecoverGenerated(manifest,root)).To(MatchError(ContainSubstring("unsafe stale generated path")))
		Expect(outside).To(BeAnExistingFile())
	})
	It("rejects unreadable manifest and stale entries",func(){
		root:=GinkgoT().TempDir()
		Expect(coverage.RecoverGenerated(root,root)).NotTo(Succeed())
		manifest:=filepath.Join(root,"manifest")
		entry:=filepath.Join(root,"zz_directory")
		Expect(os.Mkdir(entry,0755)).To(Succeed())
		writeFixture(manifest,entry+"\n")
		Expect(coverage.RecoverGenerated(manifest,root)).NotTo(Succeed())
	})
	It("does not overwrite maintained files or leave partial generated output",func(){
		root:=GinkgoT().TempDir()
		overlay:=filepath.Join(root,"overlay.json")
		manifest:=filepath.Join(root,"manifest")
		source:=filepath.Join(root,"source.go")
		target:=filepath.Join(root,"zz_generated.go")
		writeFixture(source,"// Code generated.\npackage example\n")
		_,err:=coverage.Materialize(overlay,manifest)
		Expect(err).To(HaveOccurred())
		writeFixture(overlay,"not-json")
		_,err=coverage.Materialize(overlay,manifest)
		Expect(err).To(HaveOccurred())
		writeOverlay(overlay,map[string]string{target:source})
		_,err=coverage.Materialize(overlay,filepath.Join(root,"missing","manifest"))
		Expect(err).To(HaveOccurred())
		writeFixture(target,"package maintained\n")
		_,err=coverage.Materialize(overlay,manifest)
		Expect(err).To(MatchError(ContainSubstring("refusing to overwrite maintained source")))
		Expect(os.ReadFile(target)).To(Equal([]byte("package maintained\n")))
		Expect(os.Remove(target)).To(Succeed())
		writeOverlay(overlay,map[string]string{target:filepath.Join(root,"missing.go")})
		_,err=coverage.Materialize(overlay,manifest)
		Expect(err).To(HaveOccurred())
		Expect(target).NotTo(BeAnExistingFile())
		Expect(manifest).NotTo(BeAnExistingFile())
	})
})

func coverageFixture() (string,string,string,[]string) {
	root:=GinkgoT().TempDir()
	generator:=filepath.Join(root,"generate.sh")
	reports:=filepath.Join(root,"reports")
	writeFixture(filepath.Join(root,"go.mod"),"module example.com/coveragefixture\n\ngo 1.20\n")
	writeFixture(filepath.Join(root,"services","sample","sample.go"),"package sample\nfunc Answer() int { return 42 }\n")
	writeFixture(filepath.Join(root,"services","sample","internal","resources","tests","unit-test","sample_test.go"),"package sample_test\nimport (\"testing\"; \"example.com/coveragefixture/services/sample\")\nfunc TestAnswer(t *testing.T) { if sample.Answer()!=42 {t.Fatal(\"wrong answer\")} }\n")
	writeFixture(generator,"#!/bin/sh\nset -eu\nmkdir -p \"$1/.generated\"\nprintf '%s' '{\"Replace\":{}}' > \"$1/.generated/overlay.json\"\n")
	return root,generator,reports,[]string{"-backend",root,"-core",root,"-generator-script",generator,"-reports",reports,"-package","sample"}
}

var _=Describe("Coverage command execution",func(){
	It("executes actual temporary Go package tests and enforces strict coverage",func(){
		root,_,reports,args:=coverageFixture()
		Expect(coverage.Run(args)).To(Succeed())
		Expect(filepath.Join(reports,"sample.html")).To(BeAnExistingFile())
		profile:=filepath.Join(root,".generated","unit-coverage","sample.out")
		total,covered,err:=coverage.Score(profile)
		Expect(err).NotTo(HaveOccurred())
		Expect(total).To(BeNumerically(">",0))
		Expect(covered).To(Equal(total))
		writeFixture(filepath.Join(root,"services","sample","sample.go"),"package sample\nfunc Answer() int { return 42 }\nfunc Unused() int { return 7 }\n")
		Expect(coverage.Run(args)).To(MatchError("one or more selected unit coverage runs failed"))
		Expect(coverage.Run(append(args,"-mode","--report-only"))).To(Succeed())
	})
	It("rejects bad arguments, generator failure, invalid overlay and invalid report directory",func(){
		root,generator,reports,args:=coverageFixture()
		Expect(coverage.Run([]string{"-unknown-option"})).NotTo(Succeed())
		Expect(coverage.Run([]string{"unexpected"})).NotTo(Succeed())
		Expect(coverage.Run(nil)).To(MatchError("core, generator-script and reports are required"))
		Expect(coverage.Run(append(args,"-mode","invalid"))).To(MatchError(ContainSubstring("unknown unit testing mode")))
		Expect(coverage.Run(append(args,"-packages","sample"))).To(MatchError("use package= or packages=, not both"))
		writeFixture(generator,"#!/bin/sh\nexit 3\n")
		Expect(coverage.Run(args)).NotTo(Succeed())
		writeFixture(generator,"#!/bin/sh\nmkdir -p \"$1/.generated\"\nprintf invalid > \"$1/.generated/overlay.json\"\n")
		Expect(coverage.Run(args)).NotTo(Succeed())
		writeFixture(generator,"#!/bin/sh\nmkdir -p \"$1/.generated\"\nprintf '%s' '{\"Replace\":{}}' > \"$1/.generated/overlay.json\"\n")
		writeFixture(reports,"not a directory")
		Expect(coverage.Run(args)).NotTo(Succeed())
		Expect(filepath.Join(root,".generated","unit-coverage",".materialized-files")).NotTo(BeAnExistingFile())
	})
})
