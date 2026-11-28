package mappers_test

import (
	"os"
	"path/filepath"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/annotation/internal/src/generators/mappers"
)

var _ = Describe("GenerateAll", func() {
	It("Generate Resource Owned Mapper", func() {
		t := GinkgoT()
	backend:=t.TempDir()
	Expect(os.WriteFile(filepath.Join(backend,"go.mod"),[]byte("module example.org/internal-platform\n"),0644)).To(Succeed())
	owner:=filepath.Join(backend,"services","foundation")
	files:=map[string]string{
		filepath.Join(owner,"dto","person","response.go"): `package person
type PersonErrorMetadata struct { GID string }`,
		filepath.Join(owner,"mappers","person","person.go"): `package person
import persondto "example.org/internal-platform/services/foundation/dto/person"
// @Mapper PersonMapper
type PersonMapper interface {
	// @Mapping GID=gid
	ToPersonErrorMetadata(gid string) persondto.PersonErrorMetadata
}`,
	}
	for path,source:=range files {
		Expect(os.MkdirAll(filepath.Dir(path),0755)).To(Succeed())
		Expect(os.WriteFile(path,[]byte(source),0644)).To(Succeed())
	}
	Expect(mappers.GenerateAll(backend)).To(Succeed())
	output:=filepath.Join(owner,".generated","mappers","person","impl","zz_mappers.gen.go")
	code,err:=os.ReadFile(output)
	Expect(err).NotTo(HaveOccurred())
	for _,expected:=range []string{
		"func (mapper *PersonMapper) ToPersonErrorMetadata(",
		"persondto.PersonErrorMetadata{",
		"GID: gid,",
	}{
		Expect(string(code)).To(ContainSubstring(expected))
	}
	{
		_,err:=os.Stat(filepath.Join(owner,".generated","mappers","impl","zz_mappers.gen.go"))
		Expect(os.IsNotExist(err)).To(BeTrue())
	}

	})
})
