package rest_test

import (
	"context"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	corerequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	rest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/internal/src/operations/rest"
)

type optionalQuery struct {
	Value *string `query:"value"`
	Ignored string `query:"-"`
	Untagged string
}

type optionalQueryRequest struct {
	corerequest.Request
	Query optionalQuery
}

var _ = Describe("Optional and ignored fields",func(){
	It("skips untagged and ignored fields in the typed binder",func(){
		w:=runQueryCase("/",false,func(_ context.Context,input rest.Input[rest.Empty,optionalQuery,rest.Empty])(rest.Result[string],error){
			Expect(input.Query.Value).To(BeNil())
			Expect(input.Query.Ignored).To(BeEmpty())
			Expect(input.Query.Untagged).To(BeEmpty())
			return rest.OK("ok"),nil
		})
		Expect(w.Code).To(Equal(200))
	})
	It("skips ignored and absent optional fields in the structured binder",func(){
		w:=runStructCase(httptest.NewRequest("GET","/",nil),true,func(_ context.Context,input optionalQueryRequest)(rest.Result[string],error){
			Expect(input.Query.Value).To(BeNil())
			Expect(input.Query.Ignored).To(BeEmpty())
			Expect(input.Query.Untagged).To(BeEmpty())
			return rest.OK("ok"),nil
		},nil)
		Expect(w.Code).To(Equal(200))
	})
})
