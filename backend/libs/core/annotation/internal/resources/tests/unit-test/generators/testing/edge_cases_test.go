package testgen_test

import (
	"os"
	"path/filepath"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	testgen "grovehub/libs/core/annotation/internal/src/generators/testing"
)

func writeTestSource(path, source string) {
	Expect(os.MkdirAll(filepath.Dir(path),0755)).To(Succeed())
	Expect(os.WriteFile(path,[]byte(source),0644)).To(Succeed())
}

var _=Describe("Test generator edge cases",func(){
	It("registers suite hooks and specs in stable order, skips generated and vendor files, and cleans stale outputs",func(){
		root:=GinkgoT().TempDir()
		dir:=filepath.Join(root,"tests","unit-test","example")
		source:=filepath.Join(dir,"example_test.go")
		writeTestSource(source,"package example_test\n// @BeforeAll\nfunc setup() {}\n// @BeforeEach\nfunc reset() {}\n// @Testing suite=Example\nfunc first() {}\n// @Testing suite=Example\nfunc second() {}\n")
		writeTestSource(filepath.Join(dir,"zz_invalid.go"),"ignored invalid source")
		writeTestSource(filepath.Join(dir,"vendor","invalid.go"),"ignored invalid source")
		writeTestSource(filepath.Join(dir,".generated","invalid.go"),"ignored invalid source")
		Expect(testgen.GenerateAll(root,root)).To(Succeed())
		output:=filepath.Join(root,".generated","tests","unit-test","example","zz_suite_test.go")
		generated,err:=os.ReadFile(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(generated)).To(ContainSubstring("BeforeSuite(setup)"))
		Expect(string(generated)).To(ContainSubstring("BeforeEach(reset)"))
		Expect(string(generated)).To(ContainSubstring("first()"))
		Expect(string(generated)).To(ContainSubstring("second()"))
		Expect(testgen.GenerateAll(root,root)).To(Succeed())
		Expect(os.ReadFile(output)).To(Equal(generated))
		Expect(os.Remove(source)).To(Succeed())
		Expect(testgen.GenerateAll(root,root)).To(Succeed())
		Expect(output).NotTo(BeAnExistingFile())
	})
	It("discovers tests in separate backend services and ignores files at the service root",func(){
		backend:=GinkgoT().TempDir()
		core:=GinkgoT().TempDir()
		service:=filepath.Join(backend,"services","sample")
		writeTestSource(filepath.Join(service,"internal","resources","tests","unit-test","sample_test.go"),"package sample_test\n// @Testing suite=Sample\nfunc specs() {}\n")
		writeTestSource(filepath.Join(backend,"services","README.md"),"not a service")
		writeTestSource(filepath.Join(core,"go.mod"),"module example.com/core\n")
		Expect(testgen.GenerateAll(backend,core)).To(Succeed())
		Expect(filepath.Join(service,".generated","internal","resources","tests","unit-test","zz_suite_test.go")).To(BeAnExistingFile())
		Expect(testgen.GenerateAll(GinkgoT().TempDir(),core)).NotTo(Succeed())
		Expect(testgen.GenerateAll(backend,filepath.Join(core,"missing"))).NotTo(Succeed())
	})
	DescribeTable("rejects invalid test annotations",
		func(filename,source,message string){
			root:=GinkgoT().TempDir()
			writeTestSource(filepath.Join(root,"tests","unit-test",filename),"package example_test\n"+source+"\n")
			Expect(testgen.GenerateAll(root,root)).To(MatchError(ContainSubstring(message)))
		},
		Entry("syntax","example_test.go","func broken(","expected"),
		Entry("normal source","example.go","// @Testing suite=Demo\nfunc spec() {}","must decorate a *_test.go function"),
		Entry("parameters","example_test.go","// @Testing suite=Demo\nfunc spec(v int) {}","must take no parameters"),
		Entry("results","example_test.go","// @Testing suite=Demo\nfunc spec() int {return 1}","must take no parameters"),
		Entry("generic","example_test.go","// @Testing suite=Demo\nfunc spec[T any]() {}","must take no parameters"),
		Entry("suite missing","example_test.go","// @Testing\nfunc spec() {}","requires exactly suite=Name"),
		Entry("suite invalid","example_test.go","// @Testing suite=Demo-Name\nfunc spec() {}","suite name must contain only letters and digits"),
		Entry("suite digit","example_test.go","// @Testing suite=1Demo\nfunc spec() {}","suite name must start with a letter"),
		Entry("suite conflict","example_test.go","// @Testing suite=A\nfunc first() {}\n// @Testing suite=B\nfunc second() {}","multiple @Testing suite names"),
		Entry("duplicate directive","example_test.go","// @Testing suite=Demo\n// @Testing suite=Demo\nfunc spec() {}","duplicate @Testing"),
		Entry("before all args","example_test.go","// @BeforeAll arg\nfunc setup() {}","@BeforeAll does not accept arguments"),
		Entry("before each args","example_test.go","// @BeforeEach arg\nfunc setup() {}","@BeforeEach does not accept arguments"),
		Entry("duplicate before all","example_test.go","// @BeforeAll\nfunc one() {}\n// @BeforeAll\nfunc two() {}","multiple @BeforeAll hooks"),
		Entry("orphan hook","example_test.go","// @BeforeEach\nfunc setup() {}","testing hook requires an @Testing suite"),
		Entry("unknown directive","example_test.go","// @NotARealDirective\nfunc setup() {}","unknown annotation directive"),
		Entry("mock args","example.go","// @Mock args\ntype Store interface{}","@Mock does not accept arguments"),
		Entry("mock struct","example.go","// @Mock\ntype Store struct{}","@Mock requires interface"),
		Entry("mock test file","example_test.go","// @Mock\ntype Store interface{}","@Mock interface must be in a normal Go file"),
		Entry("mock ownership","example.go","// @Mock\ntype Store interface{}","@Mock is limited to test-owned mocks"),
	)
	It("rejects inconsistent packages and multiple mock source files",func(){
		root:=GinkgoT().TempDir()
		dir:=filepath.Join(root,"tests","unit-test")
		writeTestSource(filepath.Join(dir,"a_test.go"),"package a_test\n// @Testing suite=Demo\nfunc one() {}")
		writeTestSource(filepath.Join(dir,"b_test.go"),"package b_test\n// @Testing suite=Demo\nfunc two() {}")
		Expect(testgen.GenerateAll(root,root)).To(MatchError(ContainSubstring("inconsistent test package")))
		root=GinkgoT().TempDir()
		dir=filepath.Join(root,"tests","mocks")
		writeTestSource(filepath.Join(dir,"a.go"),"package mocks\n// @Mock\ntype A interface{}")
		writeTestSource(filepath.Join(dir,"b.go"),"package mocks\n// @Mock\ntype B interface{}")
		Expect(testgen.GenerateAll(root,root)).To(MatchError(ContainSubstring("at most one @Mock source file")))
	})
	It("rejects mock sources without an owning Go module",func(){
		root:=GinkgoT().TempDir()
		writeTestSource(filepath.Join(root,"tests","mocks","source.go"),"package mocks\n// @Mock\ntype Store interface{}")
		Expect(testgen.GenerateAll(root,root)).To(MatchError(ContainSubstring("no owning Go module")))
	})
})
