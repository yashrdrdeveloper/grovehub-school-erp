package transaction_test

import (
	"context"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database/internal/src/transaction"
	databasemock "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/tests/mocks/database"
)

var _ = Describe("Internal transaction facade",func(){
	It("delegates transaction state, resolution, callbacks and unit of work",func(){
		db,mock,cleanup,err:=databasemock.NewPostgres()
		Expect(err).NotTo(HaveOccurred())
		DeferCleanup(func(){mock.ExpectClose();Expect(cleanup()).To(Succeed());Expect(mock.ExpectationsWereMet()).To(Succeed())})
		mock.ExpectBegin()
		mock.ExpectCommit()
		called:=false
		unit:=transaction.NewUnitOfWork(db)
		Expect(unit.Run(context.Background(),func(ctx context.Context)error{
			Expect(transaction.Active(ctx)).To(BeTrue())
			Expect(transaction.Database(ctx,db).Statement.ConnPool).NotTo(BeIdenticalTo(db.Statement.ConnPool))
			Expect(transaction.AfterCommit(ctx,func(){called=true})).To(BeTrue())
			Expect(transaction.Run(ctx,db,func(nested context.Context)error{
				Expect(transaction.Active(nested)).To(BeTrue())
				return nil
			})).To(Succeed())
			Expect(called).To(BeFalse())
			return nil
		})).To(Succeed())
		Expect(called).To(BeTrue())
		Expect(transaction.Active(context.Background())).To(BeFalse())
	})
})
