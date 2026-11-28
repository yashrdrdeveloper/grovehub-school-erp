package query_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/query"
)

var _ = Describe("Multi-sort query", func() {
	It("preserves mixed sort directions and appends a stable cursor tie-breaker", func() {
		db, _ := newDB(true)
		input := baseQuery()
		input.SecondarySorts = []contracts.QuerySort{{Field: "status", Direction: contracts.DirectionDescending}}
		input.Cursor = &contracts.Cursor[string]{Value: "Ada", SecondaryValues: []string{"ACTIVE"}, CursorId: "id-1"}
		result, err := query.Apply(db.Model(&record{}), input, queryConfig())
		Expect(err).NotTo(HaveOccurred())
		statement := result.Find(&[]record{}).Statement
		Expect(statement.SQL.String()).To(ContainSubstring("ORDER BY name ASC,status DESC,id ASC"))
		Expect(statement.SQL.String()).To(ContainSubstring("status <"))
		Expect(statement.Vars).To(ContainElement("id-1"))
	})
	It("rejects duplicate or unknown sort fields and invalid cursor cardinality", func() {
		db, _ := newDB(true)
		for _, secondary := range []contracts.QuerySort{{Field: "name", Direction: contracts.DirectionAscending}, {Field: "unknown", Direction: contracts.DirectionAscending}, {Field: "status", Direction: "invalid"}} {
			input := baseQuery()
			input.SecondarySorts = []contracts.QuerySort{secondary}
			_, err := query.Apply(db, input, queryConfig())
			Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		}
		input := baseQuery()
		input.SecondarySorts = []contracts.QuerySort{{Field: "status", Direction: contracts.DirectionDescending}}
		input.Cursor = &contracts.Cursor[string]{Value: "Ada", CursorId: "id"}
		_, err := query.Apply(db, input, queryConfig())
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
	})
})
