package impl

import (
	"context"
	"errors"
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
	return service.recordSelf(ctx,"REAUTHENTICATED",func(tx context.Context) error { return service.Service.Reauthenticate(tx,password) })
}
