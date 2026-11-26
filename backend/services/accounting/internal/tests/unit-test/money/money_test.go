package money_test

import (
 "math"
 "testing"

 "github.com/onsi/ginkgo/v2"
 "github.com/onsi/gomega"
 "grovehub/services/accounting/internal/app/services/money"
)

func TestMoney(t *testing.T) {
 gomega.RegisterFailHandler(ginkgo.Fail)
 ginkgo.RunSpecs(t, "Accounting Money Suite")
}

var _ = ginkgo.Describe("Add", func() {
 ginkgo.It("adds exact minor units", func() {
  result, err := money.Add(101, 202)
  gomega.Expect(err).NotTo(gomega.HaveOccurred())
  gomega.Expect(result).To(gomega.Equal(money.Amount(303)))
 })
 ginkgo.It("rejects overflow in both directions", func() {
  _, err := money.Add(math.MaxInt64, 1)
  gomega.Expect(err).To(gomega.MatchError(money.ErrOverflow))
  _, err = money.Add(math.MinInt64, -1)
  gomega.Expect(err).To(gomega.MatchError(money.ErrOverflow))
 })
})

var _ = ginkgo.Describe("Subtract", func() {
 ginkgo.It("supports minimum integer without negating it", func() {
  result, err := money.Subtract(-1, math.MinInt64)
  gomega.Expect(err).NotTo(gomega.HaveOccurred())
  gomega.Expect(result).To(gomega.Equal(money.Amount(math.MaxInt64)))
 })
 ginkgo.It("rejects overflow in both directions", func() {
  _, err := money.Subtract(0, math.MinInt64)
  gomega.Expect(err).To(gomega.MatchError(money.ErrOverflow))
  _, err = money.Subtract(math.MinInt64, 1)
  gomega.Expect(err).To(gomega.MatchError(money.ErrOverflow))
 })
})
