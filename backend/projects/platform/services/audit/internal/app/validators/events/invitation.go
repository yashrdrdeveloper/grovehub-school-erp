package events

import (
	"encoding/json"
	"strings"

	coreevents "grovehub/libs/core/events"
	eventtypes "grovehub/services/audit/internal/app/constants/events"
	client "grovehub/services/audit/internal/app/dto/producer"
	"grovehub/services/audit/internal/app/validators/events/constants"
	"grovehub/services/audit/internal/app/validators/events/helpers"
)

func Invitation(raw []byte) (coreevents.Envelope[client.InvitationEvent], error) {
	var event coreevents.Envelope[client.InvitationEvent]
	if len(raw) > constants.MaxPayloadBytes || json.Unmarshal(raw, &event) != nil ||
		event.EventType != eventtypes.InvitationDelivery || event.Version != constants.EnvelopeVersion ||
		!helpers.ValidID(event.EventID) || !helpers.ValidID(event.AggregateID) || event.OccurredAt.IsZero() {
		return event, ErrInvalidEvent
	}

	data := event.Data
	if data.TenantID == "" || strings.TrimSpace(data.TenantID) != data.TenantID ||
		!helpers.ValidID(data.AccountID) || data.Attempt < 0 ||
		(data.Purpose != "ACCOUNT_SETUP" && data.Purpose != "PASSWORD_RESET" && data.Purpose != "UNKNOWN") {
		return event, ErrInvalidEvent
	}

	switch data.Operation {
	case "QUEUED":
		if data.Attempt != 0 || data.Reason != "" || data.NextAttemptAt != nil {
			return event, ErrInvalidEvent
		}
	case "CLAIMED", "DELIVERED":
		if data.Attempt < 1 || data.Reason != "" || data.NextAttemptAt != nil {
			return event, ErrInvalidEvent
		}
	case "SUPPRESSED":
		if data.Attempt < 1 || data.NextAttemptAt != nil ||
			(data.Reason != "INVALID_PROOF" && data.Reason != "ACCOUNT_INELIGIBLE" && data.Reason != "DELIVERY_UNCERTAIN") {
			return event, ErrInvalidEvent
		}
	case "RETRY_SCHEDULED":
		if data.Attempt < 1 || data.Reason != "DELIVERY_FAILED" || data.NextAttemptAt == nil || data.NextAttemptAt.IsZero() {
			return event, ErrInvalidEvent
		}
	default:
		return event, ErrInvalidEvent
	}

	return event, nil
}
