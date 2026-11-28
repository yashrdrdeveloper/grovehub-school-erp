package validation_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"errors"
	"strings"
	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/operations/validation"
)

var _ = Describe("Request normalization",func(){
	It("normalizes search text and enforces limits",func(){
		terms,err:=validation.NormalizeSearchTerms("  Ada   LOVELACE ",30,2,nil)
		Expect(err).NotTo(HaveOccurred())
		Expect(terms).To(Equal([]string{"ada","lovelace"}))
		for _,limit:=range []int{-1,0,1}{
			_,err:=validation.NormalizeSearchTerms("abc",limit,2,nil)
			Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		}
		_,err=validation.NormalizeSearchTerms("a b",10,1,nil)
		Expect(err).To(MatchError(contracts.ErrInvalidQuery))
		rows,err:=validation.ResolveRows(0,nil)
		Expect(err).NotTo(HaveOccurred());Expect(rows).To(Equal(25))
		_,err=validation.ResolveRows(101,nil);Expect(err).To(HaveOccurred())
	})
	It("normalizes sort and filters with resource-owned policies",func(){
		fallback:=contracts.QuerySort{Field:"name",Direction:contracts.DirectionAscending}
		allowed:=map[string]contracts.CursorValueValidator{"name":validation.ValidateNonEmptyCursor}
		sort,err:=validation.NormalizeSort([]contracts.Sort{{Field:" name "}},fallback,allowed,nil)
		Expect(err).NotTo(HaveOccurred());Expect(sort).To(Equal(fallback))
		_,err=validation.NormalizeSort([]contracts.Sort{{Field:"unknown"}},fallback,allowed,nil);Expect(err).To(HaveOccurred())
		sentinel:=errors.New("bad value")
		fields:=map[string]contracts.FieldSpec{"name":{Operators:validation.TextOperators(),Normalizer:func(_ contracts.Operator,value string)(string,error){if value=="bad"{return "",sentinel};return strings.TrimSpace(value),nil}}}
		filters,err:=validation.NormalizeFilters([]contracts.Filter{{Field:" name ",Operator:contracts.OperatorEquals,Values:[]string{" Ada "}}},fields,2,2,nil)
		Expect(err).NotTo(HaveOccurred());Expect(filters[0].Values).To(Equal([]string{"Ada"}))
		_,err=validation.NormalizeFilters([]contracts.Filter{{Field:"name",Operator:contracts.OperatorEquals,Values:[]string{"bad"}}},fields,2,2,nil);Expect(err).To(MatchError(sentinel))
		_,err=validation.NormalizeFilters([]contracts.Filter{{Field:"unknown"}},fields,2,2,nil);Expect(err).To(MatchError(contracts.ErrInvalidQuery))
	})
})
