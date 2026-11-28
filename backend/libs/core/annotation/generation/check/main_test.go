package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestValidationCommentsRequireMatchingExecutableTags(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dto.go")
	broken := "package dto\ntype Input struct {\n // @NotBlank\n Name string\n}\n"
	if err := os.WriteFile(path, []byte(broken), 0600); err != nil { t.Fatal(err) }
	if err := CheckSource(path); err == nil || !strings.Contains(err.Error(), "notblank") {
		t.Fatalf("expected missing tag to fail, got %v", err)
	}
	valid := "package dto\ntype Input struct {\n // @NotBlank\n // @Size max=100\n Name string `validate:\"notblank,maxrunes=100\"`\n}\n"
	if err := os.WriteFile(path, []byte(valid), 0600); err != nil { t.Fatal(err) }
	if err := CheckSource(path); err != nil { t.Fatalf("expected matching validation tags: %v", err) }
}

func TestStrategyAnnotationAcceptedOnInterfaceAndImplementation(t *testing.T) {
    path := filepath.Join(t.TempDir(), "strategy.go")
    source := "package strategies\n// @Strategy\ntype CreateStrategy interface { Execute() }\n// @Strategy\ntype Creator struct{}\n"
    if err := os.WriteFile(path, []byte(source), 0600); err != nil { t.Fatal(err) }
    if err := CheckSource(path); err != nil { t.Fatalf("strategy interface and implementation should be accepted: %v", err) }
}

func TestProviderAnnotationValidatesStaticMetadata(t *testing.T) {
    path:=filepath.Join(t.TempDir(),"provider.go")
    valid:="package examples\n// @Provider id=GMAIL category=mail\ntype GmailProvider struct{}\n"
    if err:=os.WriteFile(path,[]byte(valid),0600);err!=nil{t.Fatal(err)}
    if err:=CheckSource(path);err!=nil{t.Fatalf("expected valid provider: %v",err)}
    invalid:="package examples\n// @Provider id=gmail category=mail\ntype GmailProvider struct{}\n"
    if err:=os.WriteFile(path,[]byte(invalid),0600);err!=nil{t.Fatal(err)}
    if err:=CheckSource(path);err==nil||!strings.Contains(err.Error(),"uppercase"){t.Fatalf("expected invalid provider ID: %v",err)}
}

func TestProviderAnnotationRejectsDuplicateIDs(t *testing.T) {
    root:=t.TempDir()
    for _,name:=range []string{"first","second"} {
        dir:=filepath.Join(root,"services",name)
        if err:=os.MkdirAll(dir,0700);err!=nil{t.Fatal(err)}
        content:="package "+name+"\n// @Provider id=GMAIL category=mail\ntype Provider struct{}\n"
        if err:=os.WriteFile(filepath.Join(dir,"provider.go"),[]byte(content),0600);err!=nil{t.Fatal(err)}
    }
    if err:=checkServices(root);err==nil||!strings.Contains(err.Error(),"duplicate"){t.Fatalf("duplicate @Provider IDs must be rejected: %v",err)}
}
