package errors_test

import (
	"errors"
	"fmt"
	"net/http"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/errors"
)

var _ = Describe("NewErrorWithMetadata", func() {
	It("exposes the same public definition and metadata contract", func() {
		definition := coreerrors.NewDefinition("INVALID", "Invalid input", http.StatusBadRequest)
		value := coreerrors.NewErrorWithMetadata(definition, "field")

		Expect(value.Definition()).To(Equal(definition))
		Expect(value.Metadata()).To(Equal("field"))
		Expect(value.Unwrap()).To(BeNil())
	})
})

var _ = Describe("Wrap", func() {
	It("preserves standard Go error chaining", func() {
		cause := errors.New("internal cause")
		definition := coreerrors.NewDefinition("FAILED", "Public message")
		value := coreerrors.Wrap(definition, cause)

		Expect(errors.Is(value, cause)).To(BeTrue())
		Expect(value.Error()).To(Equal("Public message"))
	})
})

var _ = Describe("WrapWithMetadata", func() {
	It("supports the facade error type in errors.As", func() {
		definition := coreerrors.NewDefinition("FAILED", "Public message")
		cause := errors.New("cause")
		value := coreerrors.WrapWithMetadata(definition, cause, 42)

		var found *coreerrors.Error
		Expect(errors.As(fmt.Errorf("outer: %w", value), &found)).To(BeTrue())
		Expect(found.Metadata()).To(Equal(42))
		Expect(found.Definition().Code().String()).To(Equal("FAILED"))
		Expect(errors.Is(found, cause)).To(BeTrue())
	})
})
