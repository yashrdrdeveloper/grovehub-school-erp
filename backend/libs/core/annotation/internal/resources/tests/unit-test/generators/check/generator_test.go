package check_test

import (
	"os"
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/check"
)

var _ = Describe("CheckSource", func() {
	It("Validation Comments Require Matching Executable Tags", func() {
		t := GinkgoT()
	path := filepath.Join(t.TempDir(), "dto.go")
	broken := "package dto\ntype Input struct {\n // @NotBlank\n Name string\n}\n"
	Expect(os.WriteFile(path, []byte(broken), 0600)).To(Succeed())
	{
		err := check.CheckSource(path)
		Expect(err).To(MatchError(ContainSubstring("notblank")))
	}
	valid := "package dto\ntype Input struct {\n // @NotBlank\n // @Size max=100\n Name string `validate:\"notblank,maxrunes=100\"`\n}\n"
	Expect(os.WriteFile(path, []byte(valid), 0600)).To(Succeed())
	Expect(check.CheckSource(path)).To(Succeed())

	})
})

var _ = Describe("CheckSource", func() {
	It("Strategy Annotation Accepted On Interface And Implementation", func() {
		t := GinkgoT()
    path := filepath.Join(t.TempDir(), "strategy.go")
    source := "package strategies\n// @Strategy\ntype CreateStrategy interface { Execute() }\n// @Strategy\ntype Creator struct{}\n"
    Expect(os.WriteFile(path, []byte(source), 0600)).To(Succeed())
    Expect(check.CheckSource(path)).To(Succeed())

	})
})

var _ = Describe("CheckSource", func() {
	It("Provider Annotation Validates Static Metadata", func() {
		t := GinkgoT()
    path:=filepath.Join(t.TempDir(),"provider.go")
    valid:="package examples\n// @Provider id=GMAIL category=mail\ntype GmailProvider struct{}\n"
    Expect(os.WriteFile(path,[]byte(valid),0600)).To(Succeed())
    Expect(check.CheckSource(path)).To(Succeed())
    invalid:="package examples\n// @Provider id=gmail category=mail\ntype GmailProvider struct{}\n"
    Expect(os.WriteFile(path,[]byte(invalid),0600)).To(Succeed())
    {
		err:=check.CheckSource(path)
		Expect(err).To(MatchError(ContainSubstring("uppercase")))
	}

	})
})

var _ = Describe("CheckServices", func() {
	It("Provider Annotation Rejects Duplicate IDs", func() {
		t := GinkgoT()
    root:=t.TempDir()
    for _,name:=range []string{"first","second"} {
        dir:=filepath.Join(root,"services",name)
        Expect(os.MkdirAll(dir,0700)).To(Succeed())
        content:="package "+name+"\n// @Provider id=GMAIL category=mail\ntype Provider struct{}\n"
        Expect(os.WriteFile(filepath.Join(dir,"provider.go"),[]byte(content),0600)).To(Succeed())
    }
    {
		err:=check.CheckServices(root)
		Expect(err).To(MatchError(ContainSubstring("duplicate")))
	}

	})
})
