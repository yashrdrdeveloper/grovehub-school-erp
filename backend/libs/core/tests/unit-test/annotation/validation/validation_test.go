package validation_test

import (
	"encoding/json"
	"errors"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	annotationvalidation "grovehub/libs/core/annotation/validation"
)

type CommonPerson struct {
	FirstName string `json:"firstName" normalize:"trim" validate:"notblank,maxrunes=100"`
	LastName string `json:"lastName" normalize:"trim" validate:"notblank,maxrunes=100"`
}
type createPerson struct {
	CommonPerson `validate:"valid"`
	Contacts []struct{ Value string `validate:"notblank"` } `json:"contacts" validate:"valid"`
}
type patchPerson struct {
	FirstName *string `json:"firstName,omitempty" normalize:"trim" validate:"notblank,maxrunes=100"`
	LastName *string `json:"lastName,omitempty" normalize:"trim" validate:"notblank,maxrunes=100"`
}

var _ = Describe("Nested tagged validation", func() {
	It("normalizes embedded create fields without changing wire names", func() {
		var request createPerson
		Expect(json.Unmarshal([]byte(`{"firstName":"  Arun ","lastName":" Kumar "}`), &request)).To(Succeed())
		normalized, err := annotationvalidation.Normalize(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(normalized.FirstName).To(Equal("Arun"))
		Expect(normalized.LastName).To(Equal("Kumar"))
	})

	It("preserves optional PATCH pointers and does not mutate caller data", func() {
		firstName := "  Alice  "
		request := patchPerson{FirstName: &firstName}
		normalized, err := annotationvalidation.Normalize(request)
		Expect(err).NotTo(HaveOccurred())
		Expect(*normalized.FirstName).To(Equal("Alice"))
		Expect(normalized.LastName).To(BeNil())
		Expect(firstName).To(Equal("  Alice  "))
	})

	It("fails for missing create fields and nested blank values", func() {
		_, err := annotationvalidation.Normalize(createPerson{})
		var violation annotationvalidation.Error
		Expect(errors.As(err, &violation)).To(BeTrue())
		_, err = annotationvalidation.Normalize(createPerson{
			CommonPerson: CommonPerson{FirstName:"A", LastName:"B"},
			Contacts: []struct{ Value string `validate:"notblank"` }{{Value:" "}},
		})
		Expect(errors.As(err, &violation)).To(BeTrue())
		Expect(violation.Field).To(ContainSubstring("Contacts"))
	})

	It("supports explicit non-null pointer validation", func() {
		type required struct { Name *string `validate:"nonnull"` }
		_, err := annotationvalidation.Normalize(required{})
		var violation annotationvalidation.Error
		Expect(errors.As(err, &violation)).To(BeTrue())
		Expect(violation.Rule).To(Equal("nonnull"))
	})
})
