package transaction_test

import (
 "context"
 . "github.com/onsi/ginkgo/v2"
 . "github.com/onsi/gomega"
 coretransaction "grovehub/libs/core/database/transaction"
)

var _ = Describe("Run", func() {
 It("rejects missing database", func() {
  called := false
  err := coretransaction.Run(context.Background(), nil, func(context.Context) error { called = true; return nil })
  Expect(err).To(MatchError(coretransaction.ErrMissingDatabase))
  Expect(called).To(BeFalse())
 })
 It("does not register callbacks outside transaction", func() {
  Expect(coretransaction.AfterCommit(context.Background(), func() {})).To(BeFalse())
 })
})
