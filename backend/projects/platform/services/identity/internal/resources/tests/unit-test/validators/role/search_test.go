package role_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coresearch "grovehub/libs/core/search"
	roledto "grovehub/services/identity/internal/app/dto/role"
	rolevalidator "grovehub/services/identity/internal/app/validators/role"
)

var _ = Describe("NormalizeSearch", func() {
	It("normalizes numbered pages and ordered filters", func() {
		query, err := rolevalidator.NormalizeSearch(coresearch.Request{
			PageInfo: coresearch.PageRequest{PageNumber: 2, RowSize: 10},
			Sort: []coresearch.Sort{{Field: roledto.FieldName, Direction: coresearch.DirectionAscending}, {Field: roledto.FieldCreatedAt, Direction: coresearch.DirectionDescending}},
			Filters: []coresearch.Filter{{Field: roledto.FieldStatus, Operator: coresearch.OperatorEquals, Values: []string{"inactive"}}},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(query.Offset).To(Equal(10))
		Expect(query.Limit).To(Equal(10))
		Expect(query.SecondarySorts).To(HaveLen(1))
		Expect(query.Filters[0].Values).To(Equal([]string{"INACTIVE"}))
	})
	It("rejects invalid filters, duplicate sorts and oversized pages", func() {
		for _, request := range []coresearch.Request{
			{Filters: []coresearch.Filter{{Field: roledto.FieldPersonas, Operator: coresearch.OperatorEquals, Values: []string{"UNKNOWN"}}}},
			{Sort: []coresearch.Sort{{Field: roledto.FieldName, Direction: coresearch.DirectionAscending}, {Field: roledto.FieldName, Direction: coresearch.DirectionDescending}}},
			{PageInfo: coresearch.PageRequest{RowSize: 101}},
			{PageInfo: coresearch.PageRequest{Cursor: "bad"}},
		} {
			_, err := rolevalidator.NormalizeSearch(request)
			Expect(err).To(MatchError(rolevalidator.ErrInvalidSearch))
		}
	})
})
