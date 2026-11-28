package resolve_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	databasemock "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tests/mocks/database"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/helpers"
	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction/operations/resolve"
)

var _ = Describe("Database", func() {
	It("uses the bound database instead of the fallback when state exists", func() {
		db, mock, cleanup, err := databasemock.NewPostgres()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			mock.ExpectClose()
			Expect(cleanup()).To(Succeed())
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
		ctx := context.Background()
		child, _ := helpers.Bind(ctx, db)

		Expect(resolve.Database(child, nil).Statement.Context).To(Equal(child))
		Expect(resolve.Database(ctx, db).Statement.Context).To(Equal(ctx))
	})
})
