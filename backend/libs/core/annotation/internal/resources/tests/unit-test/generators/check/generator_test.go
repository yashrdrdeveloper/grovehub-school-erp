package check_test

import (
	"os"
	"path/filepath"
	"strings"
	. "github.com/onsi/ginkgo/v2"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/check"
)

var _ = Describe("CheckSource", func() {
	It("Validation Comments Require Matching Executable Tags", func() {
		t := GinkgoT()
	path := filepath.Join(t.TempDir(), "dto.go")
	broken := "package dto\ntype Input struct {\n // @NotBlank\n Name string\n}\n"
	if err := os.WriteFile(path, []byte(broken), 0600); err != nil { t.Fatal(err) }
	if err := check.CheckSource(path); err == nil || !strings.Contains(err.Error(), "notblank") {
		t.Fatalf("expected missing tag to fail, got %v", err)
	}
	valid := "package dto\ntype Input struct {\n // @NotBlank\n // @Size max=100\n Name string `validate:\"notblank,maxrunes=100\"`\n}\n"
	if err := os.WriteFile(path, []byte(valid), 0600); err != nil { t.Fatal(err) }
	if err := check.CheckSource(path); err != nil { t.Fatalf("expected matching validation tags: %v", err) }

	})
})

var _ = Describe("CheckSource", func() {
	It("Strategy Annotation Accepted On Interface And Implementation", func() {
		t := GinkgoT()
    path := filepath.Join(t.TempDir(), "strategy.go")
    source := "package strategies\n// @Strategy\ntype CreateStrategy interface { Execute() }\n// @Strategy\ntype Creator struct{}\n"
    if err := os.WriteFile(path, []byte(source), 0600); err != nil { t.Fatal(err) }
    if err := check.CheckSource(path); err != nil { t.Fatalf("strategy interface and implementation should be accepted: %v", err) }

	})
})

var _ = Describe("CheckSource", func() {
	It("Provider Annotation Validates Static Metadata", func() {
		t := GinkgoT()
    path:=filepath.Join(t.TempDir(),"provider.go")
    valid:="package examples\n// @Provider id=GMAIL category=mail\ntype GmailProvider struct{}\n"
    if err:=os.WriteFile(path,[]byte(valid),0600);err!=nil{t.Fatal(err)}
    if err:=check.CheckSource(path);err!=nil{t.Fatalf("expected valid provider: %v",err)}
    invalid:="package examples\n// @Provider id=gmail category=mail\ntype GmailProvider struct{}\n"
    if err:=os.WriteFile(path,[]byte(invalid),0600);err!=nil{t.Fatal(err)}
    if err:=check.CheckSource(path);err==nil||!strings.Contains(err.Error(),"uppercase"){t.Fatalf("expected invalid provider ID: %v",err)}

	})
})

var _ = Describe("CheckServices", func() {
	It("Provider Annotation Rejects Duplicate IDs", func() {
		t := GinkgoT()
    root:=t.TempDir()
    for _,name:=range []string{"first","second"} {
        dir:=filepath.Join(root,"services",name)
        if err:=os.MkdirAll(dir,0700);err!=nil{t.Fatal(err)}
        content:="package "+name+"\n// @Provider id=GMAIL category=mail\ntype Provider struct{}\n"
        if err:=os.WriteFile(filepath.Join(dir,"provider.go"),[]byte(content),0600);err!=nil{t.Fatal(err)}
    }
    if err:=check.CheckServices(root);err==nil||!strings.Contains(err.Error(),"duplicate"){t.Fatalf("duplicate @Provider IDs must be rejected: %v",err)}

	})
})
