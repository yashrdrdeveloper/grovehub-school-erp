package helpers

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"strings"

	"github.com/google/uuid"
	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	client "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/entity"
)

func projectOccurrence(raw []byte) (entity.Activity,error) {
	var event coreevents.Envelope[client.Activity]
	if len(raw)>512*1024 || json.Unmarshal(raw,&event)!=nil || event.EventType!=client.ActivityOccurred || event.Version!=1 { return entity.Activity{},ErrInvalidEvent }
	data:=event.Data
	if !validID(data.OccurrenceID) || event.EventID!=data.OccurrenceID || !validID(data.AssetID) || event.AggregateID!=data.AssetID || !validID(data.ActorID) || data.TenantID=="" || strings.TrimSpace(data.TenantID)!=data.TenantID || len(data.TenantID)>128 || data.OccurredAt.IsZero() || !event.OccurredAt.Equal(data.OccurredAt) { return entity.Activity{},ErrInvalidEvent }
	if !registeredOccurrence(data) { return entity.Activity{},ErrInvalidEvent }
	key,_:=json.Marshal([]any{data.TenantID,data.Source,data.AssetType,data.AssetID,client.ActivityOccurred,data.Operation,data.OccurrenceID})
	id:=uuid.NewSHA1(namespace,key).String()
	semantic,_:=json.Marshal([]any{string(key),data.ActorID,data.OccurredAt.UTC()})
	digest:=sha256.Sum256(semantic)
	return entity.Activity{AuditID:id,SourceEventID:event.EventID,TenantID:data.TenantID,AssetType:data.AssetType,AssetID:data.AssetID,AuditType:client.ActivityOccurred,Operation:data.Operation,Source:data.Source,OccurrenceID:&data.OccurrenceID,ActorID:&data.ActorID,SubjectAccountID:&data.AssetID,Data:json.RawMessage("{}"),SemanticHash:hex.EncodeToString(digest[:]),OccurredAt:data.OccurredAt.UTC()},nil
}

func registeredOccurrence(data client.Activity) bool {
	if data.Source!="identity" || data.AssetType!="ACCOUNT" || data.ActorID!=data.AssetID { return false }
	switch data.Operation {
	case "LOGIN_SUCCEEDED","PASSWORD_RESET","PASSWORD_CHANGED","LOGOUT","LOGOUT_ALL","SESSION_REVOKED","REAUTHENTICATED":
		return true
	default:
		return false
	}
}
