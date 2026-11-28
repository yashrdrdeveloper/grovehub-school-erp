package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/operations/compare/helpers"
)

var _ = Describe("Changes", func() {
	It("sorts the union of object keys and emits nested field paths", func() {
		old := map[string]any{"z": "removed", "shared": map[string]any{"name": "old"}}
		next := map[string]any{"a": "added", "shared": map[string]any{"name": "new"}}

		Expect(helpers.Changes(old, next)).To(Equal([]contracts.Change{
			{Field: "a", NewValue: "added"},
			{Field: "shared.name", OldValue: "old", NewValue: "new"},
			{Field: "z", OldValue: "removed"},
		}))
	})

	It("expands added and removed nested objects into leaf changes", func() {
		value := map[string]any{"details": map[string]any{"name": "person"}}

		Expect(helpers.Changes(nil, value)).To(Equal([]contracts.Change{
			{Field: "details.name", NewValue: "person"},
		}))
		Expect(helpers.Changes(value, nil)).To(Equal([]contracts.Change{
			{Field: "details.name", OldValue: "person"},
		}))
	})

	It("preserves the historical equivalence of absent and explicit null object fields", func() {
		Expect(helpers.Changes(map[string]any{}, map[string]any{"name": nil})).To(BeEmpty())
		Expect(helpers.Changes(nil, map[string]any{})).To(BeEmpty())
	})
})
