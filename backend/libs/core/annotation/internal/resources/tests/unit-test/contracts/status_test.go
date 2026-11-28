package contract_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	contract "grovehub/libs/core/annotation/internal/src/contracts"
)

var _ = Describe("Status", func() {
	It("preserves the declared metadata contract", func() {
		Expect([]contract.Status{contract.Existing, contract.Gombok, contract.Proposed, contract.Optional, contract.Discouraged, contract.Pilot}).To(Equal([]contract.Status{"EXISTING", "GOMBOK", "PROPOSED", "OPTIONAL", "DISCOURAGED", "PILOT"}))
	})
})
