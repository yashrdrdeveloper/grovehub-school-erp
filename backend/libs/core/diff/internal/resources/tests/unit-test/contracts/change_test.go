package contracts_test

import (
	"encoding/json"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/diff/internal/src/contracts"
)

var _ = Describe("Change", func() {
	It("preserves the JSON change contract including explicit null values", func() {
		change := contracts.Change{Field: "name", NewValue: "new"}

		encoded, err := json.Marshal(change)
		Expect(err).NotTo(HaveOccurred())
		Expect(encoded).To(MatchJSON(`{"field":"name","oldValue":null,"newValue":"new"}`))

		var decoded contracts.Change
		Expect(json.Unmarshal(encoded, &decoded)).To(Succeed())
		Expect(decoded).To(Equal(change))
	})
})
