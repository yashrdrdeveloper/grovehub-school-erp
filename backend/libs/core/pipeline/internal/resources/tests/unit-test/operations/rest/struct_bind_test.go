package rest_test

import (
	"context"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
	corerequest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request"
	rest "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline/internal/src/operations/rest"
)

type strictQuery struct {
	corerequest.Request
	Query struct { Count int `query:"count"` }
}

var _ = Describe("Struct query binding", func() {
	It("rejects malformed, unknown and duplicate scalar parameters before service execution", func() {
		for _, raw := range []string{"?count=%zz", "?extra=1", "?count=1&count=2"} {
			w := httptest.NewRecorder()
			r := httptest.NewRequest(http.MethodGet, "/items"+raw, nil)
			rest.RunStruct(w, r, rest.StructConfig[strictQuery, string]{
				StrictQuery: true,
				InvalidQuery: coreerrors.NewDefinition("INVALID_QUERY","Invalid query",400),
				Fallback: coreerrors.NewDefinition("INTERNAL","Internal",500),
				Service: func(context.Context, strictQuery)(rest.Result[string],error){ Fail("unexpected service call"); return rest.OK(""),nil },
			})
			Expect(w.Code).To(Equal(400))
		}
	})
})
