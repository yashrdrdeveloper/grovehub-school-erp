package contracts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/utils/internal/src/contracts"
)

var _ = Describe("Normalization errors", func() {
	It("retains the original error messages", func() {
		Expect(contracts.ErrInvalidEmail).To(MatchError("invalid email"))
		Expect(contracts.ErrInvalidPhone).To(MatchError("invalid phone"))
		Expect(contracts.ErrInvalidUUIDv7).To(MatchError("invalid UUIDv7"))
		Expect(contracts.ErrInvalidTime).To(MatchError("invalid RFC3339 time"))
	})
})
