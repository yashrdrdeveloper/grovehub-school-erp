package contracts_test

import (
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/headercontext/internal/src/contracts"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("Values", func() {
	It("retains metadata as a value copy", func() {
		values := contracts.Values{IdempotencyKey: "key", CorrelationID: "corr", RequestID: "request"}
		copy := values
		copy.IdempotencyKey = ""
		Expect(values.IdempotencyKey).To(Equal("key"))
		Expect(copy.CorrelationID).To(Equal("corr"))
		Expect(copy.RequestID).To(Equal("request"))
		Expect(contracts.Values{}).To(Equal(contracts.Values{IdempotencyKey: "", CorrelationID: "", RequestID: ""}))
	})
})
