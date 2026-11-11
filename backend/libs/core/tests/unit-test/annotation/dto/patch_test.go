package dto_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	annotationdto "grovehub/libs/core/annotation/dto"
)

var _ = Describe("Reusable presence-aware PATCH field", func() {
	It("distinguishes omission, explicit null, and supplied values", func() {
		type input struct { Name annotationdto.Patch[string] `json:"name"` }
		var missing input
		Expect(json.Unmarshal([]byte(`{}`), &missing)).To(Succeed())
		Expect(missing.Name.Set).To(BeFalse())

		var cleared input
		Expect(json.Unmarshal([]byte(`{"name":null}`), &cleared)).To(Succeed())
		Expect(cleared.Name.Set).To(BeTrue())
		Expect(cleared.Name.Value).To(BeNil())

		var updated input
		Expect(json.Unmarshal([]byte(`{"name":"Mira"}`), &updated)).To(Succeed())
		Expect(updated.Name.Set).To(BeTrue())
		Expect(*updated.Name.Value).To(Equal("Mira"))
	})
	It("rejects incompatible JSON values", func() {
		type input struct { Name annotationdto.Patch[string] `json:"name"` }
		var value input
		Expect(json.Unmarshal([]byte(`{"name":42}`), &value)).NotTo(Succeed())
	})
})
