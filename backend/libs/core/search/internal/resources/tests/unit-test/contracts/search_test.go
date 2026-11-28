package contracts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"encoding/json"
	"grovehub/libs/core/search/internal/src/contracts"
)

var _ = Describe("Search contracts",func(){
	It("preserves JSON field names and operator tokens",func(){
		payload,err:=json.Marshal(contracts.Request{Filters:[]contracts.Filter{{Field:"status",Operator:contracts.OperatorEquals,Values:[]string{"ACTIVE"}}},PageInfo:contracts.PageRequest{PageNumber:2,RowSize:25}})
		Expect(err).NotTo(HaveOccurred())
		Expect(string(payload)).To(MatchJSON(`{"filters":[{"field":"status","operator":"EQUALS","values":["ACTIVE"]}],"pageInfo":{"pageNumber":2,"rowSize":25}}`))
	})
})
