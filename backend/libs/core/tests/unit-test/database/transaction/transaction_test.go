package transaction_test

import (
 "context"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 coredatabase "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
)

var _ = Describe("Run", func() {
 It("rejects missing database", func() {
  called := false
  err := coredatabase.RunTransaction(context.Background(), nil, func(context.Context) error { called = true; return nil })
  Expect(err).To(MatchError(coredatabase.ErrMissingDatabase))
  Expect(called).To(BeFalse())
 })
 It("does not register callbacks outside transaction", func() {
  Expect(coredatabase.AfterCommit(context.Background(), func() {})).To(BeFalse())
 })
})
