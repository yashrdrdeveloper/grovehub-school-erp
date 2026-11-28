package inspect_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/request/internal/src/operations/inspect"
)

type allSections struct {
	contracts.Request
	Params      struct{ GID string }
	Query       struct{ Page int }
	RequestBody struct{ Name string }
}

type bodyOnly struct {
	contracts.Request
	RequestBody string
}

type privateField struct {
	contracts.Request
	secret string
}

type unsupportedSection struct {
	contracts.Request
	Params string
	Extra string
}

type missingMarker struct {
	Query string
}

func (missingMarker) IsRequest() {}

type scalarRequest string

func (scalarRequest) IsRequest() {}

var _ = Describe("Inspect", func() {
	It("accepts the empty marker by value or pointer", func() {
		for _, input := range []any{contracts.Request{}, &contracts.Request{}} {
			sections, err := inspect.Inspect(input)
			Expect(err).NotTo(HaveOccurred())
			Expect(sections).To(Equal(contracts.Sections{}))
		}
	})

	It("discovers all declared sections and keeps absent sections invalid", func() {
		sections, err := inspect.Inspect(allSections{})
		Expect(err).NotTo(HaveOccurred())
		Expect(sections.HasParams && sections.HasQuery && sections.HasRequestBody).To(BeTrue())
		Expect(sections.Params.IsValid()).To(BeTrue())
		Expect(sections.Query.IsValid()).To(BeTrue())
		Expect(sections.RequestBody.IsValid()).To(BeTrue())

		sections, err = inspect.Inspect(bodyOnly{RequestBody: "value"})
		Expect(err).NotTo(HaveOccurred())
		Expect(sections.HasParams).To(BeFalse())
		Expect(sections.HasQuery).To(BeFalse())
		Expect(sections.HasRequestBody).To(BeTrue())
		Expect(sections.Params.IsValid()).To(BeFalse())
		Expect(sections.Query.IsValid()).To(BeFalse())
		Expect(sections.RequestBody.String()).To(Equal("value"))
		Expect(sections.RequestBody.CanSet()).To(BeFalse())
	})

	It("returns writable section references for pointer inputs", func() {
		input := &allSections{}
		sections, err := inspect.Inspect(input)
		Expect(err).NotTo(HaveOccurred())
		Expect(sections.Params.CanAddr()).To(BeTrue())
		Expect(sections.Query.CanSet()).To(BeTrue())
		Expect(sections.RequestBody.CanSet()).To(BeTrue())
		sections.Params.FieldByName("GID").SetString("gid")
		sections.Query.FieldByName("Page").SetInt(2)
		sections.RequestBody.FieldByName("Name").SetString("Ada")
		Expect(input.Params.GID).To(Equal("gid"))
		Expect(input.Query.Page).To(Equal(2))
		Expect(input.RequestBody.Name).To(Equal("Ada"))
	})

	It("rejects nil, missing contracts and non-struct implementations", func() {
		var nilPointer *bodyOnly
		for _, item := range []struct {
			input any
			message string
		}{
			{nil, "nil request"},
			{struct{}{}, "request must implement common Request contract"},
			{nilPointer, "nil request pointer"},
			{scalarRequest("scalar"), "request must be a struct"},
			{missingMarker{}, "request must embed the common request.Request"},
			{privateField{}, "request includes an unexported field"},
			{unsupportedSection{}, "unsupported request section \"Extra\""},
		} {
			sections, err := inspect.Inspect(item.input)
			Expect(err).To(MatchError(item.message))
			Expect(sections).To(Equal(contracts.Sections{}))
		}
	})
})
