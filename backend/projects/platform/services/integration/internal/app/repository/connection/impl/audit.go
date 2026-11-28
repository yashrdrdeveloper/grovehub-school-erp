package impl

import (
	"context"
	"time"

	"github.com/google/uuid"
	schoolgateway "grovehub/libs/platform/session"
	auditclient "grovehub/services/audit/client/producer"
	"grovehub/services/integration/internal/app/entity"
)

func (r *ConnectionRepository) recordActivity(ctx context.Context,entry entity.Connection,operation string,at time.Time) error {
	if r.audit==nil { return nil }
	actorID:=""
	if actor,ok:=schoolgateway.PrincipalFrom(ctx);ok && actor.TenantID==entry.TenantID { actorID=actor.AccountID }
	// Public OAuth callback is correlated to the stored connection, not a
	// forged account principal. Never send external IDs, tokens or test errors.
	return r.audit.Record(ctx,auditclient.Activity{OccurrenceID:uuid.NewString(),TenantID:entry.TenantID,Source:"integration",AssetType:"CONNECTION",AssetID:entry.ID,ActorID:actorID,PersonID:entry.OwnerGID,Operation:operation,OccurredAt:at})
}
