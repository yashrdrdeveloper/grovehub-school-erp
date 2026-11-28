package search_test

import (
	"os"
	"path/filepath"
	"strings"
	. "github.com/onsi/ginkgo/v2"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/search"
)

const searchTestInterface = `package repository
import (
	"context"
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
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
	coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
	"grovehub/services/demo/entity"
)
type PersonRepository struct { db *gorm.DB }

// @Searchable config=search.PersonConfig alias=p
func (repository *PersonRepository) Search(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error) {
	return repository.searchGenerated(ctx, query)
}
`

const searchTestConfig = `package search
import coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
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
	if err := os.WriteFile(filepath.Join(backend, "go.mod"), []byte("module grovehub\n"), 0644); err != nil { t.Fatal(err) }
	repo := filepath.Join(backend, "services", "demo", "repository")
	for rel, content := range map[string]string{
		"person.go": iface,
		"impl/person.go": impl,
		"helpers/search/person.go": searchTestConfig,
	} {
		filename := filepath.Join(repo, rel)
		if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil { t.Fatal(err) }
		if err := os.WriteFile(filename, []byte(content), 0644); err != nil { t.Fatal(err) }
	}
	output := filepath.Join(backend, "services", "demo", ".generated", "repository", "impl", "zz_search.gen.go")
	return backend, output
}

var _ = Describe("GenerateAll", func() {
	It("Concrete Search Annotation Generates Private Helper", func() {
		t := GinkgoT()
	root, output := searchFixture(t, searchTestInterface, searchTestImpl)
	if err := search.GenerateAll(root); err != nil { t.Fatal(err) }
	compiled, err := os.ReadFile(output)
	if err != nil { t.Fatal(err) }
	for _, fragment := range []string{
		"func (repository *PersonRepository) searchGenerated(",
		"coresearch.Execute[string, entity.Person]",
		"search.PersonConfig",
	} {
		if !strings.Contains(string(compiled), fragment) { t.Fatalf("missing %q:\n%s", fragment, compiled) }
	}
	if strings.Contains(string(compiled), "func (repository *PersonRepository) Search(") {
		t.Fatal("generator must not duplicate the public Search method")
	}
	if err := search.GenerateAll(root); err != nil { t.Fatal(err) }
	repeated, err := os.ReadFile(output)
	if err != nil || string(repeated) != string(compiled) { t.Fatalf("generation is not idempotent: %v", err) }

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Annotation Must Be On Implementation", func() {
		t := GinkgoT()
	annotated := strings.Replace(searchTestInterface, "\tSearch(ctx context.Context", "\t// @Searchable config=search.PersonConfig alias=p\n\tSearch(ctx context.Context", 1)
	root, _ := searchFixture(t, annotated, searchTestImpl)
	if err := search.GenerateAll(root); err == nil || !strings.Contains(err.Error(), "move @Searchable") {
		t.Fatalf("expected misplaced annotation error, got: %v", err)
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Annotation Requires Generated Helper Delegation", func() {
		t := GinkgoT()
	manual := strings.Replace(searchTestImpl, "return repository.searchGenerated(ctx, query)", "return coresearch.QueryResult[entity.Person]{}, nil", 1)
	root, _ := searchFixture(t, searchTestInterface, manual)
	if err := search.GenerateAll(root); err == nil || !strings.Contains(err.Error(), "searchGenerated") {
		t.Fatalf("expected delegation error, got: %v", err)
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Annotation Rejects Entity Mismatch", func() {
		t := GinkgoT()
	wrong := strings.Replace(searchTestImpl, "config=search.PersonConfig alias=p", "config=search.PersonConfig alias=p entity=Contact", 1)
	root, _ := searchFixture(t, searchTestInterface, wrong)
	if err := search.GenerateAll(root); err == nil || !strings.Contains(err.Error(), "must match") {
		t.Fatalf("expected entity mismatch, got: %v", err)
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
	if err := search.GenerateAll(root); err == nil || !strings.Contains(err.Error(), "already implements searchGenerated") {
		t.Fatalf("expected helper collision error, got: %v", err)
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Generator Removes Stale Output", func() {
		t := GinkgoT()
	root, output := searchFixture(t, searchTestInterface, searchTestImpl)
	if err := search.GenerateAll(root); err != nil { t.Fatal(err) }
	source := strings.Replace(searchTestImpl, "// @Searchable config=search.PersonConfig alias=p", "// unannotated", 1)
	file := filepath.Join(root, "services", "demo", "repository", "impl", "person.go")
	if err := os.WriteFile(file, []byte(source), 0644); err != nil { t.Fatal(err) }
	if err := search.GenerateAll(root); err != nil { t.Fatal(err) }
	if _, err := os.Stat(output); !os.IsNotExist(err) { t.Fatalf("stale generator output retained: %v", err) }

	})
})

var _ = Describe("GenerateAll", func() {
	It("Search Imports Use Target Module Path", func() {
		t := GinkgoT()
	root, output := searchFixture(t, searchTestInterface, searchTestImpl)
	if err := os.WriteFile(filepath.Join(root, "go.mod"), []byte("module example.org/internal-platform\n"), 0644); err != nil { t.Fatal(err) }
	if err := search.GenerateAll(root); err != nil { t.Fatal(err) }
	data, err := os.ReadFile(output)
	if err != nil { t.Fatal(err) }
	content := string(data)
	if !strings.Contains(content, "example.org/internal-platform/services/demo/entity") {
		t.Fatal("generated entity import must use target module")
	}
	if strings.Contains(content, "grovehub/services/") {
		t.Fatal("generated search imports contain the school module")
	}

	})
})
