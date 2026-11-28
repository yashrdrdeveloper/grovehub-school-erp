package rest_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "grovehub/libs/core/errors"
	rest "grovehub/libs/core/pipeline/internal/src/operations/rest"
)

var _ = Describe("REST result writers", func() {
	It("constructs created and existing responses without changing the payload", func() {
		Expect(rest.Created("new").StatusCode).To(Equal(http.StatusCreated))
		Expect(rest.CreatedOrOK("new",true)).To(Equal(rest.Created("new")))
		Expect(rest.CreatedOrOK("old",false)).To(Equal(rest.OK("old")))
	})

	It("does not expose unexpected service errors", func() {
		w := httptest.NewRecorder()
		rest.Run(w,httptest.NewRequest("GET","/",nil),rest.Config[rest.Empty,rest.Empty,rest.Empty,string]{
			Fallback: coreerrors.NewDefinition("INTERNAL","Internal",500),
			Service: func(context.Context,rest.Input[rest.Empty,rest.Empty,rest.Empty])(rest.Result[string],error){return rest.Result[string]{},errors.New("private detail")},
		})
		Expect(w.Code).To(Equal(500))
		Expect(w.Body.String()).NotTo(ContainSubstring("private detail"))
	})
})
