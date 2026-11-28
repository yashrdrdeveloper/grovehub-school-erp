package helpers_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	databasemock "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tests/mocks/database"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/helpers"
)

var _ = Describe("Bind", func() {
	It("binds shared state to a child context without mutating its parent", func() {
		db, mock, cleanup, err := databasemock.NewPostgres()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			mock.ExpectClose()
			Expect(cleanup()).To(Succeed())
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
		ctx := context.Background()

		child, state := helpers.Bind(ctx, db)
		found, ok := helpers.FromContext(child)
		Expect(ok).To(BeTrue())
		Expect(found).To(BeIdenticalTo(state))
		Expect(state.Database(child).Statement.Context).To(Equal(child))

		_, ok = helpers.FromContext(ctx)
		Expect(ok).To(BeFalse())
	})
})

var _ = Describe("PublishCallbacks", func() {
	It("executes callbacks in registration order", func() {
		_, state := helpers.Bind(context.Background(), nil)
		var order []int

		state.AddCallback(func() { order = append(order, 1) })
		state.AddCallback(func() { order = append(order, 2) })
		Expect(order).To(BeEmpty())

		state.PublishCallbacks()

		Expect(order).To(Equal([]int{1, 2}))
	})
})
