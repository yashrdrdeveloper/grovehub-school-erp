package errors_test

import (
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "grovehub/libs/core/errors"
)

var _ = Describe("Error definition HTTP status", func() {
	It("preserves the original code and message alongside a declared HTTP status", func() {
		definition := coreerrors.NewDefinition("TEST_NOT_FOUND", "Not found", http.StatusNotFound)
		Expect(definition.Code()).To(Equal(coreerrors.Code("TEST_NOT_FOUND")))
		Expect(definition.Message()).To(Equal("Not found"))
		Expect(definition.HTTPStatus()).To(Equal(http.StatusNotFound))
		Expect(coreerrors.NewError(definition).Definition().HTTPStatus()).To(Equal(http.StatusNotFound))
	})

	It("defaults legacy two-argument and zero-value definitions to HTTP 500", func() {
		Expect(coreerrors.NewDefinition("TEST_LEGACY", "Legacy").HTTPStatus()).To(Equal(http.StatusInternalServerError))
		Expect(coreerrors.Definition{}.HTTPStatus()).To(Equal(http.StatusInternalServerError))
	})

	It("rejects invalid status values or multiple status arguments", func() {
		Expect(func() { coreerrors.NewDefinition("TEST", "Bad", http.StatusOK) }).To(Panic())
		Expect(func() { coreerrors.NewDefinition("TEST", "Bad", 600) }).To(Panic())
		Expect(func() { coreerrors.NewDefinition("TEST", "Bad", http.StatusBadRequest, http.StatusConflict) }).To(Panic())
	})
})
