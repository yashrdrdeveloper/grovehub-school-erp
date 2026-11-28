package output_test

import (
	"os"
	"path/filepath"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/helpers/output"
)

var _ = Describe("WriteGo", func() {
	It("formats source and does not rewrite identical output", func() {
		path := filepath.Join(GinkgoT().TempDir(), "generated", "source.go")
		Expect(output.WriteGo(path, []byte("package sample\nfunc Answer()int{return 42}"), "generate")).To(Succeed())
		data, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("package sample\n\nfunc Answer() int { return 42 }\n"))
		stamp := time.Unix(1000,0)
		Expect(os.Chtimes(path, stamp, stamp)).To(Succeed())
		Expect(output.WriteGo(path, data, "generate")).To(Succeed())
		info, err := os.Stat(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(info.ModTime()).To(Equal(stamp))
		Expect(output.WriteGo(path, []byte("package sample\nfunc Answer()int{return 7}"), "")).To(Succeed())
		data, err = os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(ContainSubstring("return 7"))
	})
	It("reports invalid Go syntax without modifying an existing file", func() {
		path := filepath.Join(GinkgoT().TempDir(), "source.go")
		Expect(os.WriteFile(path, []byte("maintained"), 0644)).To(Succeed())
		Expect(output.WriteGo(path, []byte("package sample\nfunc broken("), "format client")).To(MatchError(ContainSubstring("format client:")))
		Expect(output.WriteGo(path, []byte("package sample\nfunc broken("), "")).NotTo(Succeed())
		data, err := os.ReadFile(path)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal("maintained"))
	})
})

var _ = Describe("WriteIfChanged",func(){
	It("rejects non-file output and blocked paths",func(){
		root:=GinkgoT().TempDir()
		blocked:=filepath.Join(root,"blocked")
		Expect(os.WriteFile(blocked,[]byte("file"),0644)).To(Succeed())
		Expect(output.WriteIfChanged(filepath.Join(blocked,"source.go"),[]byte("new"))).NotTo(Succeed())
		Expect(output.WriteIfChanged(root,[]byte("new"))).NotTo(Succeed())
		link:=filepath.Join(root,"dangling")
		Expect(os.Symlink(filepath.Join(root,"missing","source.go"),link)).To(Succeed())
		Expect(output.WriteIfChanged(link,[]byte("new"))).NotTo(Succeed())
	})
})
