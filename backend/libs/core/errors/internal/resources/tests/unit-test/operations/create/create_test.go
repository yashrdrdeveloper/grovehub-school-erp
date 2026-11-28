package create_test

import (
	"net/http"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors/internal/src/operations/create"
)

var _ = Describe("NewError", func() {
	It("creates an error without a cause or metadata", func() {
		definition := contracts.NewDefinition("TEST", "Message", http.StatusNotFound)

		value := create.NewError(definition)

		Expect(value.Definition()).To(Equal(definition))
		Expect(value.Unwrap()).To(BeNil())
		Expect(value.Metadata()).To(BeNil())
	})
})

var _ = Describe("NewErrorWithMetadata", func() {
	It("retains metadata without introducing a cause", func() {
		definition := contracts.NewDefinition("TEST", "Message", http.StatusBadRequest)
		metadata := []string{"first", "second"}

		value := create.NewErrorWithMetadata(definition, metadata)

		Expect(value.Metadata()).To(Equal(metadata))
		Expect(value.Unwrap()).To(BeNil())
		Expect(value.Definition()).To(Equal(definition))
	})
})
