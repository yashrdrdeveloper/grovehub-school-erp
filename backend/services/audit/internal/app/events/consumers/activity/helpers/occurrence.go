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
	if !validID(data.OccurrenceID) || event.EventID!=data.OccurrenceID || !validID(data.AssetID) || event.AggregateID!=data.AssetID || (data.ActorID!="" && !validID(data.ActorID)) || (data.PersonID!="" && !validID(data.PersonID)) || data.AssetVersion<0 || data.TenantID=="" || strings.TrimSpace(data.TenantID)!=data.TenantID || len(data.TenantID)>128 || data.OccurredAt.IsZero() || !event.OccurredAt.Equal(data.OccurredAt) { return entity.Activity{},ErrInvalidEvent }
	if !registeredOccurrence(data) { return entity.Activity{},ErrInvalidEvent }
	identity:=any(data.OccurrenceID)
	occurrence:=&data.OccurrenceID
	if data.AssetVersion>0 {
		if data.OccurrenceID!=VersionEventID(data) { return entity.Activity{},ErrInvalidEvent }
		identity=data.AssetVersion
		occurrence=nil
	}
	key,_:=json.Marshal([]any{data.TenantID,data.Source,data.AssetType,data.AssetID,client.ActivityOccurred,data.Operation,identity})
	id:=uuid.NewSHA1(namespace,key).String()
	semantic,_:=json.Marshal([]any{string(key),data.ActorID,data.PersonID,data.OccurredAt.UTC()})
	digest:=sha256.Sum256(semantic)
	var subject *string
	if data.AssetType=="ACCOUNT" { subject=&data.AssetID }
	return entity.Activity{AuditID:id,SourceEventID:event.EventID,TenantID:data.TenantID,AssetType:data.AssetType,AssetID:data.AssetID,AuditType:client.ActivityOccurred,Operation:data.Operation,Source:data.Source,OccurrenceID:occurrence,AssetVersion:data.AssetVersion,ActorID:optionalID(data.ActorID),PersonID:optionalID(data.PersonID),SubjectAccountID:subject,Data:json.RawMessage("{}"),SemanticHash:hex.EncodeToString(digest[:]),OccurredAt:data.OccurredAt.UTC()},nil
}

// VersionEventID derives an independent stable outbox ID, not a random UUID.
func VersionEventID(data client.Activity) string {
	key,_:=json.Marshal([]any{"versioned-source",data.TenantID,data.Source,data.AssetType,data.AssetID,data.Operation,data.AssetVersion})
	return uuid.NewSHA1(namespace,key).String()
}

func registeredOccurrence(data client.Activity) bool {
	switch data.Source {
	case "identity":
		if data.AssetType!="ACCOUNT" || data.AssetVersion!=0 || data.PersonID!="" { return false }
		if data.Operation=="LOGIN_FAILED" { return data.ActorID=="" }
		if data.ActorID!=data.AssetID { return false }
		switch data.Operation {
		case "LOGIN_SUCCEEDED","PASSWORD_RESET","PASSWORD_CHANGED","LOGOUT","LOGOUT_ALL","SESSION_REVOKED","REAUTHENTICATED","REAUTHENTICATION_FAILED":
			return true
		}
	case "foundation":
		if data.AssetType!="PERSON" || data.AssetVersion<1 || data.ActorID=="" || data.PersonID!=data.AssetID { return false }
		switch data.Operation {
		case "CREATED","UPDATED","DEACTIVATED":
			return true
		}
	case "accounting":
        return (data.AssetType=="ACCOUNTING_ENTITY" || data.AssetType=="ACCOUNTING_BOOK") && data.AssetVersion==1 && data.ActorID!="" && data.PersonID=="" && data.Operation=="CREATED"
	case "integration":
		if data.AssetType!="CONNECTION" || data.AssetVersion!=0 || !validID(data.PersonID) { return false }
		switch data.Operation {
		case "CREATED","ACTIVATED","DISCONNECTED","TESTED","SYSTEM_MAIL_SENDER_CONFIGURED":
			return true
		}
	}
	return false
}
