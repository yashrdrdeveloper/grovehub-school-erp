package run_test

import (
	"context"
	"errors"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"grovehub/libs/core/database/internal/src/transaction/operations/run"
	databasemock "grovehub/libs/core/tests/mocks/database"
)

var _ = Describe("Run", func() {
	It("returns a begin failure without invoking the operation", func() {
		db, mock, cleanup, err := databasemock.NewPostgres()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func() {
			mock.ExpectClose()
			Expect(cleanup()).To(Succeed())
			Expect(mock.ExpectationsWereMet()).To(Succeed())
		})
		failure := errors.New("begin failed")
		mock.ExpectBegin().WillReturnError(failure)
		called := false

		err = run.Run(context.Background(), db, func(context.Context) error {
			called = true
			return nil
		})

		Expect(err).To(MatchError(failure))
		Expect(called).To(BeFalse())
	})
})
