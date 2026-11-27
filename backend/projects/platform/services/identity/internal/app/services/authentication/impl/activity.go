package impl

import (
	"context"
	"errors"
	"log/slog"
	"grovehub/services/identity/internal/app/entity"
	"grovehub/services/identity/internal/app/services/authentication"
	"time"

	"github.com/google/uuid"
	auditclient "grovehub/services/audit/client/producer"
	"grovehub/services/identity/internal/app/services/security/principal"
)

func (service *RecoveryGuard) record(ctx context.Context,tenant string,accountID string,operation string) error {
	if service.audit==nil { return errors.New("authentication audit producer unavailable") }
	return service.audit.Record(ctx,auditclient.Activity{OccurrenceID:uuid.NewString(),TenantID:tenant,Source:"identity",AssetType:"ACCOUNT",AssetID:accountID,ActorID:accountID,Operation:operation,OccurredAt:time.Now().UTC()})
}

func (service *RecoveryGuard) recordSelf(ctx context.Context,operation string,action func(context.Context) error) error {
	actor,err:=principal.Require(ctx)
	if err!=nil { return err }
	if actor.TenantID=="" { return action(ctx) }
	return service.locker.Run(ctx,actor.TenantID,func(tx context.Context) error {
		if err:=action(tx);err!=nil { return err }
		return service.record(tx,actor.TenantID,actor.AccountID,operation)
	})
}

func (service *RecoveryGuard) Logout(ctx context.Context) error {
	return service.recordSelf(ctx,"LOGOUT",service.Service.Logout)
}

func (service *RecoveryGuard) LogoutAll(ctx context.Context) (int64,error) {
	var count int64
	err:=service.recordSelf(ctx,"LOGOUT_ALL",func(tx context.Context) error {
		var err error
		count,err=service.Service.LogoutAll(tx)
		return err
	})
	if err!=nil { return 0,err }
	return count,nil
}

func (service *RecoveryGuard) RevokeSession(ctx context.Context,sessionID string) error {
	return service.recordSelf(ctx,"SESSION_REVOKED",func(tx context.Context) error { return service.Service.RevokeSession(tx,sessionID) })
}

func (service *RecoveryGuard) Reauthenticate(ctx context.Context,password string) error {
	err:=service.recordSelf(ctx,"REAUTHENTICATED",func(tx context.Context) error { return service.Service.Reauthenticate(tx,password) })
	if errors.Is(err,authentication.ErrInvalidCredentials) { service.recordRejectedReauthentication(ctx) }
	return err
}

// Rejected credentials must not create a session or claim the target account
// authenticated. Persist separately, after the failed business transaction
// rolls back; keep the public error identical even when auditing is unavailable.
func (service *RecoveryGuard) recordRejectedLogin(ctx context.Context,account entity.Account,identifier string) {
	if account.TenantID==nil || *account.TenantID=="" { return }
	tenant:=*account.TenantID
	bounded,cancel:=context.WithTimeout(ctx,2*time.Second)
	defer cancel()
	err:=service.locker.Run(bounded,tenant,func(tx context.Context) error {
		current,err:=service.accounts.FindByID(tx,account.AccountID)
		if err!=nil { return err }
		if current.TenantID==nil || *current.TenantID!=tenant || current.Email!=identifier { return nil }
		if service.audit==nil { return errors.New("authentication audit producer unavailable") }
		return service.audit.Record(tx,auditclient.Activity{OccurrenceID:uuid.NewString(),TenantID:tenant,Source:"identity",AssetType:"ACCOUNT",AssetID:current.AccountID,Operation:"LOGIN_FAILED",OccurredAt:time.Now().UTC()})
	})
	if err!=nil { slog.WarnContext(ctx,"authentication failure audit unavailable","operation","LOGIN_FAILED") }
}

func (service *RecoveryGuard) recordRejectedReauthentication(ctx context.Context) {
	actor,err:=principal.Require(ctx)
	if err!=nil || actor.TenantID=="" { return }
	bounded,cancel:=context.WithTimeout(ctx,2*time.Second)
	defer cancel()
	err=service.locker.Run(bounded,actor.TenantID,func(tx context.Context) error {
		account,err:=service.accounts.FindByID(tx,actor.AccountID)
		if err!=nil { return err }
		if !eligibleRecoveryAccount(account,actor.TenantID) { return nil }
		return service.record(tx,actor.TenantID,actor.AccountID,"REAUTHENTICATION_FAILED")
	})
	if err!=nil { slog.WarnContext(ctx,"authentication failure audit unavailable","operation","REAUTHENTICATION_FAILED") }
}
