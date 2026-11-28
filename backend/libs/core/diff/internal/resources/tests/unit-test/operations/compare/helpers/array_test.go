package helpers_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/contracts"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/operations/compare/helpers"
)

var _ = Describe("Changes", func() {
	It("compares array positions and reports newly appended values", func() {
		old := []any{"first", "second"}
		next := []any{"second", "first", "third"}

		Expect(helpers.Changes(old, next)).To(Equal([]contracts.Change{
			{Field: "[0]", OldValue: "first", NewValue: "second"},
			{Field: "[1]", OldValue: "second", NewValue: "first"},
			{Field: "[2]", NewValue: "third"},
		}))
	})

	It("reports removed array items and nested object paths", func() {
		old := map[string]any{"contacts": []any{map[string]any{"name": "old"}, "gone"}}
		next := map[string]any{"contacts": []any{map[string]any{"name": "new"}}}

		Expect(helpers.Changes(old, next)).To(Equal([]contracts.Change{
			{Field: "contacts[0].name", OldValue: "old", NewValue: "new"},
			{Field: "contacts[1]", OldValue: "gone"},
		}))
	})

	It("expands arrays added to or removed from null", func() {
		value := []any{[]any{"nested"}}

		Expect(helpers.Changes(nil, value)).To(Equal([]contracts.Change{
			{Field: "[0][0]", NewValue: "nested"},
		}))
		Expect(helpers.Changes(value, nil)).To(Equal([]contracts.Change{
			{Field: "[0][0]", OldValue: "nested"},
		}))
	})
})
