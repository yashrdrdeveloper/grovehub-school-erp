package helpers

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	client "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/entity"
)

func projectInvitation(_ context.Context, raw []byte) (entity.Activity, error) {
	var event coreevents.Envelope[client.InvitationEvent]
	if len(raw)>512*1024 || json.Unmarshal(raw,&event)!=nil || event.EventType!=client.InvitationDelivery || event.Version!=1 || !validID(event.EventID) || !validID(event.AggregateID) || event.OccurredAt.IsZero() { return entity.Activity{},ErrInvalidEvent }
	data:=event.Data
	if data.TenantID=="" || strings.TrimSpace(data.TenantID)!=data.TenantID || !validID(data.AccountID) || data.Attempt<0 || (data.Purpose!="ACCOUNT_SETUP" && data.Purpose!="PASSWORD_RESET" && data.Purpose!="UNKNOWN") { return entity.Activity{},ErrInvalidEvent }
	switch data.Operation {
	case "QUEUED":
		if data.Attempt!=0 || data.Reason!="" || data.NextAttemptAt!=nil { return entity.Activity{},ErrInvalidEvent }
	case "CLAIMED","DELIVERED":
		if data.Attempt<1 || data.Reason!="" || data.NextAttemptAt!=nil { return entity.Activity{},ErrInvalidEvent }
	case "SUPPRESSED":
		if data.Attempt<1 || data.NextAttemptAt!=nil || (data.Reason!="INVALID_PROOF" && data.Reason!="ACCOUNT_INELIGIBLE" && data.Reason!="DELIVERY_UNCERTAIN") { return entity.Activity{},ErrInvalidEvent }
	case "RETRY_SCHEDULED":
		if data.Attempt<1 || data.Reason!="DELIVERY_FAILED" || data.NextAttemptAt==nil || data.NextAttemptAt.IsZero() { return entity.Activity{},ErrInvalidEvent }
	default:
		return entity.Activity{},ErrInvalidEvent
	}
	version:=int64(data.Attempt)+1
	key,_:=json.Marshal([]any{data.TenantID,"INVITATION",event.AggregateID,event.EventType,data.Operation,version})
	id:=uuid.NewSHA1(namespace,key).String()
	account:=data.AccountID
	tenant:=data.TenantID
	data.TenantID=""
	projected,err:=json.Marshal(data)
	if err!=nil { return entity.Activity{},ErrInvalidEvent }
	semantic,_:=json.Marshal([]any{string(key),json.RawMessage(projected)})
	digest:=sha256.Sum256(semantic)
	return entity.Activity{AuditID:id,SourceEventID:event.EventID,TenantID:tenant,AssetType:"INVITATION",AssetID:event.AggregateID,AuditType:event.EventType,Operation:data.Operation,Source:"identity",AssetVersion:version,SubjectAccountID:&account,Data:projected,SemanticHash:hex.EncodeToString(digest[:]),OccurredAt:event.OccurredAt.UTC()},nil
}
