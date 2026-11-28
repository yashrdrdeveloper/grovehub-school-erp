package search_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coresearch "grovehub/libs/core/search"
)

var _ = Describe("Search facade",func(){
	It("exposes query construction, cursor helpers and normalization",func(){
		_,err:=coresearch.Apply[string](nil,coresearch.Query[string]{},coresearch.Config{})
		Expect(err).To(MatchError(coresearch.ErrInvalidQuery))
		b:=&coresearch.QueryBuilder{}
		sql,err:=coresearch.BuildScalarCondition(b,coresearch.ScalarField{Expression:"name"},coresearch.OperatorEquals,[]string{"Ada"})
		Expect(err).NotTo(HaveOccurred());Expect(sql).To(Equal("name = ?"));Expect(b.Args()).To(Equal([]any{"Ada"}))
		Expect(coresearch.EscapeLike("a%b")).To(Equal(`a\%b`))
		token,err:=coresearch.EncodeCursor(coresearch.CursorPayload{Field:"name",Direction:coresearch.DirectionAscending,Value:"Ada",CursorId:"id"})
		Expect(err).NotTo(HaveOccurred())
		result,err:=coresearch.DecodeCursor(token)
		Expect(err).NotTo(HaveOccurred());Expect(result.Value).To(Equal("Ada"))
		terms,err:=coresearch.NormalizeSearchTerms(" Ada ",10,2,nil)
		Expect(err).NotTo(HaveOccurred());Expect(terms).To(Equal([]string{"ada"}))
	})
})
