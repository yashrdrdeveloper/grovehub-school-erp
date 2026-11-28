package impl_test

import (
	"context"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/database/internal/src/transaction/operations/unitofwork/impl"
	databasemock "grovehub/libs/core/tests/mocks/database"
)

var _ = Describe("Run", func() {
	It("executes the operation in the supplied database transaction", func() {
		db, mock, cleanup, err := databasemock.NewPostgres()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			mock.ExpectClose()
			Expect(cleanup()).To(Succeed())
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
		mock.ExpectBegin()
		mock.ExpectCommit()
		called := false

		err = impl.NewUnitOfWork(db).Run(context.Background(), func(context.Context) error {
			called = true
			return nil
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(called).To(BeTrue())
	})
})
