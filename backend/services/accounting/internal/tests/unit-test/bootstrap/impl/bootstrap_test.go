package impl_test

import (
 "testing"

 "github.com/onsi/ginkgo/v2"
 "github.com/onsi/gomega"
 bootstrapimpl "grovehub/services/accounting/internal/app/bootstrap/impl"
)

func TestBootstrap(t *testing.T) {
 gomega.RegisterFailHandler(ginkgo.Fail)
 ginkgo.RunSpecs(t, "Accounting Bootstrap Suite")
}

var _ = ginkgo.Describe("Build", func() {
 ginkgo.It("rejects missing dependencies", func() {
  runtime, err := bootstrapimpl.New(nil, nil, nil).Build()
  gomega.Expect(err).To(gomega.HaveOccurred())
  gomega.Expect(runtime).To(gomega.BeNil())
 })
})
