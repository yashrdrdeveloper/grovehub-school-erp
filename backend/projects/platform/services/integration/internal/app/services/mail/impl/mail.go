package impl

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"
	coredb "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/database"
	"grovehub/services/identity/client/invitation"
	receipt "grovehub/services/integration/internal/app/repository/maildelivery/impl"
	"grovehub/services/integration/internal/app/services/mail"
	reconciliation "grovehub/services/integration/internal/app/services/mail/operations/reconciliation/impl"
)

func New(db *gorm.DB, jobs invitation.DeliveryReader) mail.Service {
	return reconciliation.New(jobs,receipt.NewReceiptReader(db),func(ctx context.Context, read func(context.Context,time.Time)error)error{
		if db==nil{return errors.New("reconciliation database unavailable")}
		// This diagnostic never joins a caller's mutation. Both owners use the
		// same context-bound snapshot in the current shared PostgreSQL deployment.
		return coredb.RunTransaction(coredb.DetachedContext(ctx),db,func(txCtx context.Context)error{
			tx:=coredb.TransactionDB(txCtx,db)
			if err:=tx.Exec("SET TRANSACTION ISOLATION LEVEL REPEATABLE READ, READ ONLY").Error;err!=nil{return err}
			var now time.Time
			if err:=tx.Raw("SELECT CURRENT_TIMESTAMP").Row().Scan(&now);err!=nil{return err}
			return read(txCtx,now)
		})
	})
}
