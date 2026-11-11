package helpers_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "grovehub/libs/core/errors"
	personerror "grovehub/services/foundation/constants/error"
	controllerhelpers "grovehub/services/foundation/controller/protected/helpers"
)

var _ = Describe("Person error responses", func() {
	It("uses HTTP status directly from each application error definition", func() {
		cases := []struct {
			definition coreerrors.Definition
			status int
		}{
			{personerror.InvalidPersonGID, http.StatusBadRequest},
			{personerror.InvalidPersonSearch, http.StatusBadRequest},
			{personerror.PersonNotFound, http.StatusNotFound},
			{personerror.PersonContactConflict, http.StatusConflict},
			{personerror.PersonConcurrentModification, http.StatusConflict},
			{personerror.PersonCreateFailed, http.StatusInternalServerError},
		}
		for _, testcase := range cases {
			writer := httptest.NewRecorder()
			controllerhelpers.WritePersonError(context.Background(), writer, coreerrors.NewError(testcase.definition), personerror.PersonReadFailed)
			Expect(writer.Code).To(Equal(testcase.status), string(testcase.definition.Code()))
		}
	})

	It("preserves application error metadata with the definition's status", func() {
		writer := httptest.NewRecorder()
		err := coreerrors.NewErrorWithMetadata(personerror.PersonContactConflict, map[string]string{"field": "email"})
		controllerhelpers.WritePersonError(context.Background(), writer, err, personerror.PersonReadFailed)
		Expect(writer.Code).To(Equal(http.StatusConflict))
		Expect(writer.Body.String()).To(ContainSubstring("FOUNDATION_PERSON_CONTACT_CONFLICT"))
		Expect(writer.Body.String()).To(ContainSubstring("email"))
	})

	It("treats unknown failures as internal errors regardless of fallback status", func() {
		writer := httptest.NewRecorder()
		controllerhelpers.WritePersonError(context.Background(), writer, errors.New("database detail"), personerror.InvalidPersonGID)
		Expect(writer.Code).To(Equal(http.StatusInternalServerError))
		Expect(writer.Body.String()).NotTo(ContainSubstring("database detail"))
	})
})
