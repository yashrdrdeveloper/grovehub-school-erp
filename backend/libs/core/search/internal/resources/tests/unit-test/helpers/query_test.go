package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/search/internal/src/contracts"
	"grovehub/libs/core/search/internal/src/helpers"
)

var _ = Describe("SQL condition helpers",func(){
	It("parameterizes scalar operators and escapes wildcard characters",func(){
		field:=contracts.ScalarField{Expression:"p.name",Cast:"::text"}
		for _,item:=range []struct{op contracts.Operator;values []string;sql string;args []any}{
			{contracts.OperatorEquals,[]string{"a"},"p.name = ?::text",[]any{"a"}},
			{contracts.OperatorNotEquals,[]string{"a"},"p.name <> ?::text",[]any{"a"}},
			{contracts.OperatorGreaterThan,[]string{"a"},"p.name > ?::text",[]any{"a"}},
			{contracts.OperatorGTE,[]string{"a"},"p.name >= ?::text",[]any{"a"}},
			{contracts.OperatorLessThan,[]string{"a"},"p.name < ?::text",[]any{"a"}},
			{contracts.OperatorLTE,[]string{"a"},"p.name <= ?::text",[]any{"a"}},
			{contracts.OperatorIn,[]string{"a","b"},"p.name IN (?::text, ?::text)",[]any{"a","b"}},
			{contracts.OperatorNotIn,[]string{"a"},"p.name NOT IN (?::text)",[]any{"a"}},
			{contracts.OperatorExists,nil,"NULLIF(p.name, '') IS NOT NULL",nil},
			{contracts.OperatorNotExists,nil,"NULLIF(p.name, '') IS NULL",nil},
			{contracts.OperatorContains,[]string{"a%b"},"p.name LIKE ?",[]any{`%a\%b%`}},
			{contracts.OperatorStartsWith,[]string{"a_b"},"p.name LIKE ?",[]any{`a\_b%`}},
		}{
			builder:=&helpers.QueryBuilder{}
			sql,err:=helpers.BuildScalarCondition(builder,field,item.op,item.values)
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(Equal(item.sql))
			Expect(builder.Args()).To(Equal(item.args))
		}
		Expect(helpers.EscapeLike(`a\b_%`)).To(Equal(`a\\b\_\%`))
	})
	It("binds relation fixed conditions before user-supplied IN values",func(){
		field:=contracts.RelationField{Table:"contacts",Alias:"c",ParentCondition:"c.person_id = p.id",ValueExpression:"c.value",FixedConditions:[]contracts.FixedCondition{{Expression:"c.type = ",Value:"EMAIL"}}}
		for _,op:=range []contracts.Operator{contracts.OperatorIn,contracts.OperatorNotIn}{
			builder:=&helpers.QueryBuilder{}
			sql,err:=helpers.BuildRelationCondition(builder,field,op,[]string{"a","b"})
			Expect(err).NotTo(HaveOccurred())
			Expect(sql).To(ContainSubstring("c.type = ? AND c.value IN (?, ?)"))
			Expect(builder.Args()).To(Equal([]any{"EMAIL","a","b"}))
		}
		for _,op:=range []contracts.Operator{contracts.OperatorEquals,contracts.OperatorNotEquals,contracts.OperatorContains,contracts.OperatorStartsWith,contracts.OperatorExists,contracts.OperatorNotExists}{
			builder:=&helpers.QueryBuilder{}
			values:=[]string{"a"}
			if op==contracts.OperatorExists||op==contracts.OperatorNotExists {values=nil}
			_,err:=helpers.BuildRelationCondition(builder,field,op,values)
			Expect(err).NotTo(HaveOccurred())
			Expect(builder.Args()[0]).To(Equal("EMAIL"))
		}
	})
	It("rejects unsupported operators, invalid arity and unallowlisted fields",func(){
		for _,op:=range []contracts.Operator{contracts.OperatorEquals,contracts.OperatorIn,contracts.OperatorContains,"UNKNOWN"}{
			_,err:=helpers.BuildScalarCondition(&helpers.QueryBuilder{},contracts.ScalarField{},op,nil)
			Expect(err).To(MatchError(contracts.ErrUnsupportedOperator))
			_,err=helpers.BuildRelationCondition(&helpers.QueryBuilder{},contracts.RelationField{},op,nil)
			Expect(err).To(MatchError(contracts.ErrUnsupportedOperator))
		}
		for _,op:=range []contracts.Operator{contracts.OperatorExists,contracts.OperatorNotExists}{
			_,err:=helpers.BuildScalarCondition(&helpers.QueryBuilder{},contracts.ScalarField{},op,[]string{"a"})
			Expect(err).To(MatchError(contracts.ErrUnsupportedOperator))
			_,err=helpers.BuildRelationCondition(&helpers.QueryBuilder{},contracts.RelationField{},op,[]string{"a"})
			Expect(err).To(MatchError(contracts.ErrUnsupportedOperator))
		}
		_,err:=helpers.BuildFilter(&helpers.QueryBuilder{},contracts.QueryFilter{Field:"unknown"},contracts.Config{})
		Expect(err).To(MatchError(contracts.ErrInvalidField))
	})
	It("uses stable cursor ordering and tie-breakers",func(){
		b:=&helpers.QueryBuilder{}
		Expect(helpers.BuildCursorCondition(b,contracts.SortField{Expression:"p.id",Cast:"::uuid"},contracts.DirectionAscending,"a","b","p.id")).To(Equal("p.id > ?::uuid"))
		Expect(b.Args()).To(Equal([]any{"a"}))
		b=&helpers.QueryBuilder{}
		Expect(helpers.BuildCursorCondition(b,contracts.SortField{Expression:"p.name"},contracts.DirectionDescending,"a","b","p.id")).To(Equal("(p.name < ? OR (p.name = ? AND p.id < ?))"))
		Expect(b.Args()).To(Equal([]any{"a","a","b"}))
	})
})
