package search_test

import (
	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	coresearch "grovehub/libs/core/search"
)

var _ = Describe("Remaining public search operations", func() {
	It("executes a typed query and builds relational and cursor predicates", func() {
		connection, mock, err := sqlmock.New()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() { _ = connection.Close() })
		db, err := gorm.Open(postgres.New(postgres.Config{Conn: connection}), &gorm.Config{DisableAutomaticPing: true})
		Expect(err).NotTo(HaveOccurred())
		type record struct{ ID string }
		input := coresearch.Query[string]{Limit: 1, Sort: coresearch.QuerySort{Field: "id", Direction: coresearch.DirectionAscending}}
		config := coresearch.Config{CursorExpression: "id", SortFields: map[string]coresearch.SortField{"id": {Expression: "id"}}}
		mock.ExpectQuery("(?i)SELECT .*FROM").WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow("a").AddRow("b"))
		result, err := coresearch.Execute[string, record](db, input, config)
		Expect(err).NotTo(HaveOccurred())
		Expect(result.Items).To(Equal([]record{{ID: "a"}}))
		Expect(result.HasNext).To(BeTrue())
		Expect(mock.ExpectationsWereMet()).To(Succeed())
		b := &coresearch.QueryBuilder{}
		sql, err := coresearch.BuildRelationCondition(b, coresearch.RelationField{Table: "contacts", Alias: "c", ParentCondition: "c.person_id = p.id", ValueExpression: "c.value"}, coresearch.OperatorEquals, []string{"Ada"})
		Expect(err).NotTo(HaveOccurred())
		Expect(sql).To(ContainSubstring("c.value = ?"))
		Expect(b.Args()).To(Equal([]any{"Ada"}))
		b = &coresearch.QueryBuilder{}
		Expect(coresearch.BuildCursorCondition(b, coresearch.SortField{Expression: "id"}, coresearch.DirectionAscending, "a", "a", "id")).To(Equal("id > ?"))
		Expect(b.Args()).To(Equal([]any{"a"}))
	})
	It("exposes operator policies, filter normalization and cursor binding", func() {
		Expect(coresearch.OperatorSet(coresearch.OperatorEquals)).To(HaveLen(1))
		Expect(coresearch.EqualityOperators()).To(HaveLen(4))
		Expect(coresearch.TextOperators()).To(HaveLen(8))
		Expect(coresearch.TimeOperators()).To(HaveLen(5))
		Expect(coresearch.IsPartialOperator(coresearch.OperatorContains)).To(BeTrue())
		Expect(coresearch.ValidValueCount(coresearch.OperatorExists, 0)).To(BeTrue())
		rows, err := coresearch.ResolveRows(0, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(rows).To(Equal(25))
		fallback := coresearch.QuerySort{Field: "name", Direction: coresearch.DirectionAscending}
		allowed := map[string]coresearch.CursorValueValidator{"name": coresearch.ValidateNonEmptyCursor}
		sort, err := coresearch.NormalizeSort(nil, fallback, allowed, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(sort).To(Equal(fallback))
		fields := map[string]coresearch.FieldSpec{"name": {Operators: coresearch.TextOperators(), Normalizer: func(_ coresearch.Operator, value string) (string, error) { return value, nil }}}
		filters, err := coresearch.NormalizeFilters([]coresearch.Filter{{Field: "name", Operator: coresearch.OperatorEquals, Values: []string{"Ada"}}}, fields, 1, 1, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(filters).To(HaveLen(1))
		token, err := coresearch.EncodeCursor(coresearch.CursorPayload{Field: "name", Direction: coresearch.DirectionAscending, Value: "Ada", CursorId: "id"})
		Expect(err).NotTo(HaveOccurred())
		cursor, _, err := coresearch.NormalizeCursor(token, sort, coresearch.ValidateNonEmptyCursor, coresearch.ValidateNonEmptyCursor, nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(cursor.CursorId).To(Equal("id"))
		Expect(coresearch.NextPageNumber(coresearch.PageRequest{PageNumber: 2}, 2)).To(Equal(3))
	})
})
