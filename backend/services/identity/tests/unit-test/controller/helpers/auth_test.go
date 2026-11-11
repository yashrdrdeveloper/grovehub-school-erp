package helpers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "grovehub/libs/core/errors"
	identityerror "grovehub/services/identity/constants/error"
	controllerhelpers "grovehub/services/identity/controller/helpers"
)

var _ = Describe("Identity registration error responses", func() {
	It("uses the HTTP status defined with the semantic error", func() {
		cases := []struct {
			definition coreerrors.Definition
			status int
		}{
			{identityerror.InvalidRegistration, http.StatusBadRequest},
			{identityerror.InvalidRegistrationIdempotencyKey, http.StatusBadRequest},
			{identityerror.RegistrationAlreadyExists, http.StatusConflict},
			{identityerror.RegistrationIdempotencyConflict, http.StatusConflict},
			{identityerror.RegistrationFailed, http.StatusInternalServerError},
		}
		for _, testcase := range cases {
			writer := httptest.NewRecorder()
			controllerhelpers.WriteRegisterError(context.Background(), writer, coreerrors.NewError(testcase.definition))
			Expect(writer.Code).To(Equal(testcase.status), string(testcase.definition.Code()))
		}
	})

	It("uses a safe HTTP 500 for unexpected errors", func() {
		writer := httptest.NewRecorder()
		controllerhelpers.WriteRegisterError(context.Background(), writer, errors.New("private details"))
		Expect(writer.Code).To(Equal(http.StatusInternalServerError))
		Expect(writer.Body.String()).NotTo(ContainSubstring("private details"))
	})
})
