package query_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/query"
)

var _ = Describe("Query edge cases", func() {
	It("does not add the cursor ID twice when already included in the sort", func() {
		db, _ := newDB(true)
		input := baseQuery()
		input.SecondarySorts = []contracts.QuerySort{{Field: "id", Direction: contracts.DirectionDescending}}
		input.Cursor = &contracts.Cursor[string]{Value: "Ada", SecondaryValues: []string{"id-1"}, CursorId: "id-1"}
		result, err := query.Apply(db.Model(&record{}), input, queryConfig())
		Expect(err).NotTo(HaveOccurred())
		statement := result.Find(&[]record{}).Statement
		Expect(statement.SQL.String()).To(ContainSubstring("ORDER BY name ASC,id DESC"))
		Expect(statement.Vars).To(Equal([]any{"Ada", "Ada", "id-1", 3}))
	})
	It("rejects invalid queries before executing count or data queries", func() {
		for _, number := range []int{0, 1} {
			db, mock := newDB(false)
			input := baseQuery()
			input.PageNumber = number
			input.Limit = 0
			_, err := query.Execute[string, record](db, input, queryConfig())
			Expect(err).To(MatchError(contracts.ErrInvalidQuery))
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		}
	})
})
