package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGenerateResourceOwnedMapper(t *testing.T) {
	backend:=t.TempDir()
	owner:=filepath.Join(backend,"services","foundation")
	files:=map[string]string{
		filepath.Join(owner,"dto","person","response.go"): `package person
type PersonErrorMetadata struct { GID string }`,
		filepath.Join(owner,"mappers","person","person.go"): `package person
import persondto "grovehub/services/foundation/dto/person"
// @Mapper PersonMapper
type PersonMapper interface {
	// @Mapping GID=gid
	ToPersonErrorMetadata(gid string) persondto.PersonErrorMetadata
}`,
	}
	for path,source:=range files {
		if err:=os.MkdirAll(filepath.Dir(path),0755);err!=nil{t.Fatal(err)}
		if err:=os.WriteFile(path,[]byte(source),0644);err!=nil{t.Fatal(err)}
	}
	if err:=generateAll(backend);err!=nil{t.Fatal(err)}
	output:=filepath.Join(owner,".generated","mappers","person","impl","zz_mappers.gen.go")
	code,err:=os.ReadFile(output)
	if err!=nil{t.Fatal(err)}
	for _,expected:=range []string{
		"func (mapper *PersonMapper) ToPersonErrorMetadata(",
		"persondto.PersonErrorMetadata{",
		"GID: gid,",
	}{
		if !strings.Contains(string(code),expected){t.Fatalf("resource mapper output missing %q:\n%s",expected,code)}
	}
	if _,err:=os.Stat(filepath.Join(owner,".generated","mappers","impl","zz_mappers.gen.go"));!os.IsNotExist(err) {
		t.Fatalf("unexpected flat mapper output: %v",err)
	}
}
