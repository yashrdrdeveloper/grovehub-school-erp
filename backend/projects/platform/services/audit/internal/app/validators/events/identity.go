package events

import (
	"encoding/json"
	"strings"

	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	eventtypes "grovehub/services/audit/internal/app/constants/events"
	"grovehub/services/audit/internal/app/validators/events/constants"
	"grovehub/services/audit/internal/app/validators/events/helpers"
	dto "grovehub/services/audit/internal/app/dto/events"
)

func Identity(expectedType string, raw []byte) (coreevents.Envelope[dto.Payload], string, string, error) {
	var event coreevents.Envelope[dto.Payload]
	if len(raw) > constants.MaxPayloadBytes || json.Unmarshal(raw, &event) != nil || event.EventType != expectedType || event.Version != constants.EnvelopeVersion || !helpers.ValidID(event.EventID) || !helpers.ValidID(event.AggregateID) || event.OccurredAt.IsZero() || event.Data.New == nil || event.Data.New.Version < 1 {
		return event, "", "", ErrInvalidEvent
	}
	data := event.Data
	tenant := data.TenantID
	if tenant == "" { tenant = data.New.TenantID }
	if tenant == "" || strings.TrimSpace(tenant) != tenant || (data.New.TenantID != "" && data.New.TenantID != tenant) || (data.Old != nil && data.Old.TenantID != "" && data.Old.TenantID != tenant) || data.Operation == "" || len(data.Operation) > 80 {
		return event, "", "", ErrInvalidEvent
	}
	assetType := ""
	switch expectedType {
	case eventtypes.AccountChanged, eventtypes.AccountRegistrationChanged, eventtypes.AccountMetadataChanged, eventtypes.AccountAuthorityChanged:
		assetType = "ACCOUNT"
	case eventtypes.RoleChanged, eventtypes.RoleMetadataChanged:
		assetType = "ROLE"
	case eventtypes.ApprovalChanged:
		assetType = "APPROVAL"
	default:
		return event, "", "", ErrInvalidEvent
	}
	for _, value := range []*dto.Snapshot{data.Old,data.New} {
		if value == nil { continue }
		if (value.ID != "" && value.ID != event.AggregateID) || (value.AccountID != "" && value.AccountID != event.AggregateID) { return event, "", "", ErrInvalidEvent }
	}
	if assetType == "APPROVAL" && data.RequestID != event.AggregateID { return event, "", "", ErrInvalidEvent }
	return event, tenant, assetType, nil
}
