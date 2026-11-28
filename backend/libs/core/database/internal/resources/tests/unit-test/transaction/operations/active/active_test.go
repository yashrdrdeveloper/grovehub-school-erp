package active_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/database/internal/src/transaction/helpers"
	"grovehub/libs/core/database/internal/src/transaction/operations/active"
)

var _ = Describe("Active", func() {
	It("detects bound state only in the transaction context", func() {
		ctx := context.Background()
		child, _ := helpers.Bind(ctx, nil)

		Expect(active.Active(ctx)).To(BeFalse())
		Expect(active.Active(child)).To(BeTrue())
	})
})
