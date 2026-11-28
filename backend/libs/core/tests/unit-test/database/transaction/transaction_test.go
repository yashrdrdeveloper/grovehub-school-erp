package transaction_test

import (
	"context"
	"errors"

	"github.com/DATA-DOG/go-sqlmock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"gorm.io/gorm"

	coredatabase "grovehub/libs/core/database"
	databasemock "grovehub/libs/core/tests/mocks/database"
)

var _ = Describe("RunTransaction", func() {
	It("rejects a missing database without executing the operation", func() {
		called := false

		err := coredatabase.RunTransaction(context.Background(), nil, func(context.Context) error {
			called = true
			return nil
		})

		Expect(err).To(MatchError(coredatabase.ErrMissingDatabase))
		Expect(called).To(BeFalse())
	})

	It("does not register callbacks outside a transaction", func() {
		ctx := context.Background()

		Expect(coredatabase.TransactionActive(ctx)).To(BeFalse())
		Expect(coredatabase.AfterCommit(ctx, func() {})).To(BeFalse())
	})
})

var _ = Describe("NewUnitOfWork", func() {
	var (
		db      *gorm.DB
		mock    sqlmock.Sqlmock
		cleanup func() error
	)

	BeforeEach(func() {
		var err error
		db, mock, cleanup, err = databasemock.NewPostgres()
		Expect(err).NotTo(HaveOccurred())

		DeferCleanup(func() {
			Expect(mock.ExpectationsWereMet()).To(Succeed())
			mock.ExpectClose()
			Expect(cleanup()).To(Succeed())
		})
	})

	It("shares a single transaction and publishes callbacks after commit in registration order", func() {
		mock.ExpectBegin()
		mock.ExpectCommit()

		ctx := context.Background()
		unit := coredatabase.NewUnitOfWork(db)
		var callbacks []string

		err := unit.Run(ctx, func(tx context.Context) error {
			Expect(coredatabase.TransactionActive(tx)).To(BeTrue())
			bound := coredatabase.TransactionDB(tx, db)
			Expect(bound.Statement.ConnPool).NotTo(BeIdenticalTo(db.Statement.ConnPool))

			Expect(coredatabase.AfterCommit(tx, func() {
				Expect(mock.ExpectationsWereMet()).To(Succeed())
				callbacks = append(callbacks, "outer")
			})).To(BeTrue())

			Expect(unit.Run(tx, func(nested context.Context) error {
				Expect(coredatabase.TransactionDB(nested, db).Statement.ConnPool).To(BeIdenticalTo(bound.Statement.ConnPool))
				Expect(coredatabase.AfterCommit(nested, func() {
					callbacks = append(callbacks, "inner")
				})).To(BeTrue())
				return nil
			})).To(Succeed())

			Expect(callbacks).To(BeEmpty())
			return nil
		})

		Expect(err).NotTo(HaveOccurred())
		Expect(callbacks).To(Equal([]string{"outer", "inner"}))
		Expect(coredatabase.TransactionActive(ctx)).To(BeFalse())
	})

	It("rolls back and discards callbacks when an operation fails", func() {
		mock.ExpectBegin()
		mock.ExpectRollback()

		failure := errors.New("operation failed")
		called := false

		err := coredatabase.NewUnitOfWork(db).Run(context.Background(), func(ctx context.Context) error {
			Expect(coredatabase.AfterCommit(ctx, func() { called = true })).To(BeTrue())
			return failure
		})

		Expect(err).To(MatchError(failure))
		Expect(called).To(BeFalse())
	})

	It("does not publish callbacks when commit fails", func() {
		failure := errors.New("commit failed")
		mock.ExpectBegin()
		mock.ExpectCommit().WillReturnError(failure)

		called := false

		err := coredatabase.RunTransaction(context.Background(), db, func(ctx context.Context) error {
			Expect(coredatabase.AfterCommit(ctx, func() { called = true })).To(BeTrue())
			return nil
		})

		Expect(err).To(MatchError(failure))
		Expect(called).To(BeFalse())
	})

	It("detaches transaction state and cancellation while preserving context values", func() {
		mock.ExpectBegin()
		mock.ExpectCommit()

		type contextKey struct{}
		parent, cancel := context.WithCancel(context.WithValue(context.Background(), contextKey{}, "trace"))
		defer cancel()

		var detached context.Context

		Expect(coredatabase.RunTransaction(parent, db, func(ctx context.Context) error {
			detached = coredatabase.DetachedContext(ctx)
			Expect(coredatabase.TransactionActive(detached)).To(BeFalse())
			Expect(coredatabase.AfterCommit(detached, func() {})).To(BeFalse())
			Expect(coredatabase.TransactionDB(detached, db).Statement.ConnPool).To(BeIdenticalTo(db.Statement.ConnPool))
			return nil
		})).To(Succeed())

		cancel()

		Expect(detached.Err()).NotTo(HaveOccurred())
		Expect(detached.Value(contextKey{})).To(Equal("trace"))
	})
})
