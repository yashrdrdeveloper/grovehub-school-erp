package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"

	fixture "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/resources/tests/helpers"
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

func searchFixture(implementation, configuration string) string {
	root:=GinkgoT().TempDir()
	fixture.Write(root,"person.go",searchTestInterface)
	fixture.Write(root,"impl/person.go",implementation)
	fixture.Write(root,"helpers/search/person.go",configuration)

	return root
}
