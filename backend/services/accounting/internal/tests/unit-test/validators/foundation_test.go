package validators_test

import (
 "testing"
 "strings"

 "github.com/onsi/ginkgo/v2"
 "github.com/onsi/gomega"
 "grovehub/services/accounting/internal/app/validators"
)

func TestValidators(t *testing.T) {
 gomega.RegisterFailHandler(ginkgo.Fail)
 ginkgo.RunSpecs(t, "Accounting Validators Suite")
}

var _ = ginkgo.Describe("ValidateIdentity", func() {
 ginkgo.It("accepts valid identity", func() {
  gomega.Expect(validators.ValidateIdentity("tenant", "MAIN", "Main book")).To(gomega.Succeed())
 })
 ginkgo.It("rejects missing scope and invalid fields", func() {
  gomega.Expect(validators.ValidateIdentity("", "MAIN", "Main")).NotTo(gomega.Succeed())
  gomega.Expect(validators.ValidateIdentity("tenant", "", "Main")).NotTo(gomega.Succeed())
  gomega.Expect(validators.ValidateIdentity("tenant", "MAIN", strings.Repeat("a", 256))).NotTo(gomega.Succeed())
 })
})

var _ = ginkgo.Describe("ValidateID", func() {
 ginkgo.It("requires a canonical nonzero UUID", func() {
  gomega.Expect(validators.ValidateID("550e8400-e29b-41d4-a716-446655440000")).To(gomega.Succeed())
  gomega.Expect(validators.ValidateID("00000000-0000-0000-0000-000000000000")).NotTo(gomega.Succeed())
  gomega.Expect(validators.ValidateID("invalid")).NotTo(gomega.Succeed())
 })
})
