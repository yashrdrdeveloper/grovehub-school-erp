package main

import (
 "os"
 "path/filepath"
 "strings"
 "testing"
)

func TestGenerateAndIdempotence(t *testing.T) {
 dir := t.TempDir()
 source := "package entity\nimport \"time\"\ntype ID string\ntype Credential struct { ID ID; ChangedAt time.Time; Version int64 }\n"
 if err := os.WriteFile(filepath.Join(dir,"credential.go"),[]byte(source),0644);err!=nil {t.Fatal(err)}
 if err:=generate(dir,"Credential","identity.credentials","");err!=nil {t.Fatal(err)}
 path:=filepath.Join(dir,"zz_credential_entity.gen.go")
 first,err:=os.ReadFile(path);if err!=nil {t.Fatal(err)}
 for _,want:=range []string{"TableName()", "NewCredential(", "GetVersion()", "SetChangedAt(", "identity.credentials", "\"time\""} {if !strings.Contains(string(first),want) {t.Errorf("missing %q",want)}}
 if err:=generate(dir,"Credential","identity.credentials","");err!=nil {t.Fatal(err)}
 second,err:=os.ReadFile(path);if err!=nil {t.Fatal(err)}
 if string(first)!=string(second) {t.Fatal("not idempotent")}
}

func TestRejectEmbeddedFields(t *testing.T) {
 dir:=t.TempDir()
 if err:=os.WriteFile(filepath.Join(dir,"entity.go"),[]byte("package entity\ntype Base struct{}\ntype Thing struct{ Base }"),0644);err!=nil {t.Fatal(err)}
 if err:=generate(dir,"Thing","things","");err==nil {t.Fatal("expected error")}
}
