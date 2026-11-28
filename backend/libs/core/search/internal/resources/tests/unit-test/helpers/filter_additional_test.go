package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search/internal/src/helpers"
)

var _ = Describe("BuildFilter",func(){
	It("dispatches relation fields using the configured relation",func(){
		config:=contracts.Config{RelationFields:map[string]contracts.RelationField{"email":{Table:"contacts",Alias:"c",ParentCondition:"c.person_id = p.id",ValueExpression:"c.value"}}}
		builder:=&helpers.QueryBuilder{}
		sql,err:=helpers.BuildFilter(builder,contracts.QueryFilter{Field:"email",Operator:contracts.OperatorEquals,Values:[]string{"a@example.com"}},config)
		Expect(err).NotTo(HaveOccurred())
		Expect(sql).To(ContainSubstring("EXISTS (SELECT 1 FROM contacts AS c"))
		Expect(builder.Args()).To(Equal([]any{"a@example.com"}))
	})
})
