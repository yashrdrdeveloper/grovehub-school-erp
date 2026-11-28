package search_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/search"
)

const searchTestInterface = `package repository
import (
	"context"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/demo/entity"
)
// @Repository entity=Person
type PersonRepository interface {
	Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error)
}
`

const searchTestImpl = `package impl
import (
	"context"
	"gorm.io/gorm"
	coresearch "grovehub/libs/core/search"
	"grovehub/services/demo/entity"
)
type PersonRepository struct { db *gorm.DB }

// @Searchable config=search.PersonConfig alias=p
func (repository *PersonRepository) Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error) {
	return repository.searchGenerated(ctx, query)
}
`

const searchTestConfig = `package search
import coresearch "grovehub/libs/core/search"
var PersonConfig = coresearch.Config{
	ScalarFields: map[string]coresearch.ScalarField{"gid": {Expression: "p.gid"}},
	RelationFields: map[string]coresearch.RelationField{},
	SortFields: map[string]coresearch.SortField{"gid": {Expression: "p.gid"}},
	CursorExpression: "p.gid",
}
`

func searchFixture(t GinkgoTInterface, iface, impl string) (string, string) {
	t.Helper()
	backend := t.TempDir()
	Expect(os.WriteFile(filepath.Join(backend, "go.mod"), []byte("module grovehub\n"), 0644)).To(Succeed())
	repo := filepath.Join(backend, "services", "demo", "repository")
	for rel, content := range map[string]string{
		"person.go":                iface,
		"impl/person.go":           impl,
		"helpers/search/person.go": searchTestConfig,
	} {
		filename := filepath.Join(repo, rel)
		Expect(os.MkdirAll(filepath.Dir(filename), 0755)).To(Succeed())
		Expect(os.WriteFile(filename, []byte(content), 0644)).To(Succeed())
	}
	output := filepath.Join(backend, "services", "demo", ".generated", "repository", "impl", "zz_search.gen.go")
	return backend, output
}

var _ = Describe("GenerateAll", func() {
	It("Concrete Search Annotation Generates Private Helper", func() {
		t := GinkgoT()
		root, output := searchFixture(t, searchTestInterface, searchTestImpl)
		Expect(search.GenerateAll(root)).To(Succeed())
		compiled, err := os.ReadFile(output)
		Expect(err).NotTo(HaveOccurred())
		for _, fragment := range []string{
			"func (repository *PersonRepository) searchGenerated(",
			"coresearch.Execute[string, entity.Person]",
			"coredatabase.TransactionDB(ctx, repository.db)",
			"libs/core/database\"",
			"search.PersonConfig",
		} {
			Expect(string(compiled)).To(ContainSubstring(fragment))
		}
		Expect(string(compiled)).NotTo(ContainSubstring("func (repository *PersonRepository) Search("))
		Expect(search.GenerateAll(root)).To(Succeed())
		repeated, err := os.ReadFile(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(repeated).To(Equal(compiled))

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Annotation Must Be On Implementation", func() {
		t := GinkgoT()
		annotated := strings.Replace(searchTestInterface, "\tSearch(ctx context.Context", "\t// @Searchable config=search.PersonConfig alias=p\n\tSearch(ctx context.Context", 1)
		root, _ := searchFixture(t, annotated, searchTestImpl)
		{
			err := search.GenerateAll(root)
			Expect(err).To(MatchError(ContainSubstring("move @Searchable")))
		}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Annotation Requires Generated Helper Delegation", func() {
		t := GinkgoT()
		manual := strings.Replace(searchTestImpl, "return repository.searchGenerated(ctx, query)", "return coresearch.QueryResult[entity.Person]{}, nil", 1)
		root, _ := searchFixture(t, searchTestInterface, manual)
		{
			err := search.GenerateAll(root)
			Expect(err).To(MatchError(ContainSubstring("searchGenerated")))
		}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Annotation Rejects Entity Mismatch", func() {
		t := GinkgoT()
		wrong := strings.Replace(searchTestImpl, "config=search.PersonConfig alias=p", "config=search.PersonConfig alias=p entity=Contact", 1)
		root, _ := searchFixture(t, searchTestInterface, wrong)
		{
			err := search.GenerateAll(root)
			Expect(err).To(MatchError(ContainSubstring("must match")))
		}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Annotation Rejects Generated Name Collision", func() {
		t := GinkgoT()
		collision := searchTestImpl + `
func (repository *PersonRepository) searchGenerated(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error) {
	return coresearch.QueryResult[entity.Person]{}, nil
}
`
		root, _ := searchFixture(t, searchTestInterface, collision)
		{
			err := search.GenerateAll(root)
			Expect(err).To(MatchError(ContainSubstring("already implements searchGenerated")))
		}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Generator Removes Stale Output", func() {
		t := GinkgoT()
		root, output := searchFixture(t, searchTestInterface, searchTestImpl)
		Expect(search.GenerateAll(root)).To(Succeed())
		source := strings.Replace(searchTestImpl, "// @Searchable config=search.PersonConfig alias=p", "// unannotated", 1)
		file := filepath.Join(root, "services", "demo", "repository", "impl", "person.go")
		Expect(os.WriteFile(file, []byte(source), 0644)).To(Succeed())
		Expect(search.GenerateAll(root)).To(Succeed())
		{
			_, err := os.Stat(output)
			Expect(os.IsNotExist(err)).To(BeTrue())
		}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Imports Use Target Module Path", func() {
		t := GinkgoT()
		root, output := searchFixture(t, searchTestInterface, searchTestImpl)
		Expect(os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.org/internal-platform\n"), 0644)).To(Succeed())
		Expect(search.GenerateAll(root)).To(Succeed())
		data, err := os.ReadFile(output)
		Expect(err).NotTo(HaveOccurred())
		content := string(data)
		Expect(content).To(ContainSubstring("example.org/internal-platform/services/demo/entity"))
		Expect(content).NotTo(ContainSubstring("grovehub/services/"))

	})
})
