package contracts_test

import (
	"errors"
	"net/http"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/errors/internal/src/contracts"
)

var _ = Describe("NewError", func() {
	It("retains its definition, cause and metadata without exposing the cause as its message", func() {
		definition := contracts.NewDefinition("TEST", "Public message", http.StatusBadRequest)
		cause := errors.New("private failure")
		metadata := map[string]any{"field": "name"}

		value := contracts.NewError(definition, cause, metadata)

		Expect(value.Error()).To(Equal("Public message"))
		Expect(value.Unwrap()).To(BeIdenticalTo(cause))
		Expect(value.Definition()).To(Equal(definition))
		Expect(value.Metadata()).To(Equal(metadata))
		Expect(errors.Is(value, cause)).To(BeTrue())
	})
})
