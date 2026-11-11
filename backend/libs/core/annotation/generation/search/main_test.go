package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
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

func searchFixture(t *testing.T, iface, impl string) (string, string) {
	t.Helper()
	backend := t.TempDir()
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

func TestConcreteSearchAnnotationGeneratesPrivateHelper(t *testing.T) {
	root, output := searchFixture(t, searchTestInterface, searchTestImpl)
	if err := generateAll(root); err != nil { t.Fatal(err) }
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
	if err := generateAll(root); err != nil { t.Fatal(err) }
	repeated, err := os.ReadFile(output)
	if err != nil || string(repeated) != string(compiled) { t.Fatalf("generation is not idempotent: %v", err) }
}

func TestSearchAnnotationMustBeOnImplementation(t *testing.T) {
	annotated := strings.Replace(searchTestInterface, "\tSearch(ctx context.Context", "\t// @Searchable config=search.PersonConfig alias=p\n\tSearch(ctx context.Context", 1)
	root, _ := searchFixture(t, annotated, searchTestImpl)
	if err := generateAll(root); err == nil || !strings.Contains(err.Error(), "move @Searchable") {
		t.Fatalf("expected misplaced annotation error, got: %v", err)
	}
}

func TestSearchAnnotationRequiresGeneratedHelperDelegation(t *testing.T) {
	manual := strings.Replace(searchTestImpl, "return repository.searchGenerated(ctx, query)", "return coresearch.QueryResult[entity.Person]{}, nil", 1)
	root, _ := searchFixture(t, searchTestInterface, manual)
	if err := generateAll(root); err == nil || !strings.Contains(err.Error(), "searchGenerated") {
		t.Fatalf("expected delegation error, got: %v", err)
	}
}

func TestSearchAnnotationRejectsEntityMismatch(t *testing.T) {
	wrong := strings.Replace(searchTestImpl, "config=search.PersonConfig alias=p", "config=search.PersonConfig alias=p entity=Contact", 1)
	root, _ := searchFixture(t, searchTestInterface, wrong)
	if err := generateAll(root); err == nil || !strings.Contains(err.Error(), "must match") {
		t.Fatalf("expected entity mismatch, got: %v", err)
	}
}

func TestSearchAnnotationRejectsGeneratedNameCollision(t *testing.T) {
	collision := searchTestImpl + `
func (repository *PersonRepository) searchGenerated(ctx context.Context, query coresearch.Query[string]) (coresearch.QueryResult[entity.Person], error) {
	return coresearch.QueryResult[entity.Person]{}, nil
}
`
	root, _ := searchFixture(t, searchTestInterface, collision)
	if err := generateAll(root); err == nil || !strings.Contains(err.Error(), "already implements searchGenerated") {
		t.Fatalf("expected helper collision error, got: %v", err)
	}
}

func TestSearchGeneratorRemovesStaleOutput(t *testing.T) {
	root, output := searchFixture(t, searchTestInterface, searchTestImpl)
	if err := generateAll(root); err != nil { t.Fatal(err) }
	source := strings.Replace(searchTestImpl, "// @Searchable config=search.PersonConfig alias=p", "// unannotated", 1)
	file := filepath.Join(root, "services", "demo", "repository", "impl", "person.go")
	if err := os.WriteFile(file, []byte(source), 0644); err != nil { t.Fatal(err) }
	if err := generateAll(root); err != nil { t.Fatal(err) }
	if _, err := os.Stat(output); !os.IsNotExist(err) { t.Fatalf("stale generator output retained: %v", err) }
}
