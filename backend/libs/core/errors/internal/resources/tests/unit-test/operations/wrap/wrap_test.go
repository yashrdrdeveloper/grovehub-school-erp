package wrap_test

import (
	"errors"
	"fmt"
	"net/http"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/errors/internal/src/contracts"
	"grovehub/libs/core/errors/internal/src/operations/wrap"
)

var _ = Describe("Wrap", func() {
	It("keeps the underlying cause available through errors.Is", func() {
		definition := contracts.NewDefinition("TEST", "Safe message", http.StatusInternalServerError)
		cause := errors.New("private cause")

		value := wrap.Wrap(definition, cause)

		Expect(errors.Is(value, cause)).To(BeTrue())
		Expect(value.Metadata()).To(BeNil())
		Expect(value.Error()).To(Equal("Safe message"))
	})
})

var _ = Describe("WrapWithMetadata", func() {
	It("supports errors.As through nested wrapping and retains metadata", func() {
		definition := contracts.NewDefinition("TEST", "Safe message", http.StatusConflict)
		cause := errors.New("cause")
		metadata := map[string]int{"attempt": 1}
		value := wrap.WrapWithMetadata(definition, cause, metadata)
		nested := fmt.Errorf("outer: %w", value)

		var found *contracts.Error
		Expect(errors.As(nested, &found)).To(BeTrue())
		Expect(found).To(BeIdenticalTo(value))
		Expect(found.Metadata()).To(Equal(metadata))
		Expect(errors.Is(nested, cause)).To(BeTrue())
	})

	It("accepts a nil cause without changing the definition", func() {
		definition := contracts.NewDefinition("TEST", "Message", http.StatusBadRequest)
		value := wrap.WrapWithMetadata(definition, nil, nil)

		Expect(value.Unwrap()).To(BeNil())
		Expect(value.Definition()).To(Equal(definition))
	})
})
