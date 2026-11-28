package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/validation"
)

var _ = Describe("ValidateQuery",func(){
	It("requires valid limits, configured expressions and allowed ordering",func(){
		config:=contracts.Config{CursorExpression:"p.id",SortFields:map[string]contracts.SortField{"id":{Expression:"p.id"}}}
		valid:=contracts.Query[string]{Limit:25,Sort:contracts.QuerySort{Field:"id",Direction:contracts.DirectionAscending}}
		field,err:=validation.ValidateQuery(valid,config)
		Expect(err).NotTo(HaveOccurred())
		Expect(field.Expression).To(Equal("p.id"))
		for _,mutate:=range []func(*contracts.Query[string]){
			func(q *contracts.Query[string]){q.Limit=0},
			func(q *contracts.Query[string]){q.Offset=-1},
			func(q *contracts.Query[string]){q.Cursor=&contracts.Cursor[string]{};q.Offset=1},
			func(q *contracts.Query[string]){q.SearchTerms=[]string{"term"}},
			func(q *contracts.Query[string]){q.Sort.Field="unknown"},
			func(q *contracts.Query[string]){q.Sort.Direction="invalid"},
		}{q:=valid;mutate(&q);_,err:=validation.ValidateQuery(q,config);Expect(err).To(MatchError(contracts.ErrInvalidQuery))}
	})
})
