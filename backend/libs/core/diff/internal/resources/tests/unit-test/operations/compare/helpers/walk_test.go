package helpers_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"grovehub/libs/core/diff/internal/src/contracts"
	"grovehub/libs/core/diff/internal/src/operations/compare/helpers"
)

var _ = Describe("Changes", func() {
	It("preserves the root path for scalar changes", func() {
		Expect(helpers.Changes(true, false)).To(Equal([]contracts.Change{{
			Field: "$", OldValue: true, NewValue: false,
		}}))
	})

	It("reports type replacements at the containing field instead of flattening mismatched types", func() {
		old := map[string]any{"data": map[string]any{"name": "before"}}
		next := map[string]any{"data": []any{"after"}}

		Expect(helpers.Changes(old, next)).To(Equal([]contracts.Change{{
			Field: "data", OldValue: old["data"], NewValue: next["data"],
		}}))
	})

	It("returns no changes for equal JSON values", func() {
		value := map[string]any{"value": json.Number("2"), "same": []any{true, nil}}

		Expect(helpers.Changes(value, value)).NotTo(BeNil())
		Expect(helpers.Changes(value, value)).To(BeEmpty())
	})
})
