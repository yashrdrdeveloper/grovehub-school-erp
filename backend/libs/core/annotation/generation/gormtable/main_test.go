package main

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestGenerate(t *testing.T){
 dir:=t.TempDir()
 content:="package sample\n// @Entity foundation.person_contacts\ntype Contact struct { Value string }\n"
 if err:=os.WriteFile(filepath.Join(dir,"contact.go"),[]byte(content),0644);err!=nil {t.Fatal(err)}
 if err:=generate(dir);err!=nil {t.Fatal(err)}
 path:=filepath.Join(dir,"zz_tables.gen.go")
 first,err:=os.ReadFile(path);if err!=nil {t.Fatal(err)}
 if !strings.Contains(string(first),`func (Contact) TableName() string`)||!strings.Contains(string(first),`"foundation.person_contacts"`){t.Fatalf("bad output: %s",first)}
 if err:=generate(dir);err!=nil {t.Fatal(err)}
 again,_:=os.ReadFile(path);if string(first)!=string(again){t.Fatal("not idempotent")}
}

func TestRejectInvalidTable(t *testing.T){
 dir:=t.TempDir()
 if err:=os.WriteFile(filepath.Join(dir,"contact.go"),[]byte("package sample\n// @Entity foundation.person_contacts;DROP\ntype Contact struct{}\n"),0644);err!=nil {t.Fatal(err)}
 if err:=generate(dir);err==nil{t.Fatal("expected invalid table error")}
}

func TestRejectDuplicateTableName(t *testing.T){
 dir:=t.TempDir()
 source:="package sample\n// @Entity foundation.person_contacts\ntype Contact struct{}\nfunc (Contact) TableName() string{return \"foundation.person_contacts\"}\n"
 if err:=os.WriteFile(filepath.Join(dir,"contact.go"),[]byte(source),0644);err!=nil {t.Fatal(err)}
 if err:=generate(dir);err==nil{t.Fatal("expected duplicate TableName error")}
}


func TestGenerateTableOnlyEntity(t *testing.T) {
    source:=t.TempDir()
    output:=filepath.Join(t.TempDir(),"entity","zz_tables.gen.go")
    code:="package entity\n// @Entity identity.password_credentials mode=table\ntype PasswordCredential struct {PasswordHash string}\n"
    if err:=os.WriteFile(filepath.Join(source,"password.go"),[]byte(code),0644);err!=nil{t.Fatal(err)}
    if err:=generateInto(source,output);err!=nil{t.Fatal(err)}
    generated,err:=os.ReadFile(output);if err!=nil{t.Fatal(err)}
    if !strings.Contains(string(generated),`func (PasswordCredential) TableName() string`)||
       !strings.Contains(string(generated),`"identity.password_credentials"`) {t.Fatalf("wrong table-only mapping: %s",generated)}
}
