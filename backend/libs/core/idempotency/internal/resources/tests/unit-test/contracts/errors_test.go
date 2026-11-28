package contracts_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/idempotency/internal/src/contracts"
)

var _ = Describe("Idempotency errors", func() {
	It("preserves distinct error identities and messages", func() {
		Expect(contracts.ErrInvalidKey).To(MatchError("invalid idempotency key"))
		Expect(contracts.ErrFingerprint).To(MatchError("unable to fingerprint idempotent request"))
		Expect(contracts.ErrInvalidKey).NotTo(BeIdenticalTo(contracts.ErrFingerprint))
	})
})
