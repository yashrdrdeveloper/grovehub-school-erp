package contracts_test

import (
	"reflect"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/request/internal/src/contracts"
)

type endpoint struct {
	contracts.Request
	Params string
}

var _ contracts.Contract = endpoint{}
var _ contracts.Contract = (*endpoint)(nil)

var _ = Describe("Request contracts", func() {
	It("provides a zero-sized marker and supports value and pointer embedding", func() {
		Expect(reflect.TypeFor[contracts.Request]().Size()).To(BeZero())
		var marker contracts.Contract = endpoint{}
		marker.IsRequest()
		marker = &endpoint{}
		marker.IsRequest()
		Expect(reflect.TypeFor[contracts.Empty]().Size()).To(BeZero())
	})

	It("preserves typed legacy parts and absent sections", func() {
		parts := contracts.Parts[string, int, contracts.Empty]{Path: "gid", Query: 3}
		Expect(parts.Path).To(Equal("gid"))
		Expect(parts.Query).To(Equal(3))
		Expect(parts.Body).To(Equal(contracts.Empty{}))
		sections := contracts.Sections{}
		Expect(sections.HasParams || sections.HasQuery || sections.HasRequestBody).To(BeFalse())
		Expect(sections.Params.IsValid()).To(BeFalse())
		Expect(sections.Query.IsValid()).To(BeFalse())
		Expect(sections.RequestBody.IsValid()).To(BeFalse())
	})
})
