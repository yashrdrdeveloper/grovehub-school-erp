package events

import (
	"encoding/json"
	"strings"

	coreevents "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events"
	eventtypes "grovehub/services/audit/internal/app/constants/events"
	client "grovehub/services/audit/internal/app/dto/producer"
	"grovehub/services/audit/internal/app/validators/events/constants"
	"grovehub/services/audit/internal/app/validators/events/helpers"
)

func Occurrence(raw []byte) (coreevents.Envelope[client.Activity], error) {
	var event coreevents.Envelope[client.Activity]
	if len(raw) > constants.MaxPayloadBytes || json.Unmarshal(raw, &event) != nil ||
		event.EventType != eventtypes.ActivityOccurred || event.Version != constants.EnvelopeVersion {
		return event, ErrInvalidEvent
	}

	data := event.Data
	if !helpers.ValidID(data.OccurrenceID) || event.EventID != data.OccurrenceID ||
		!helpers.ValidID(data.AssetID) || event.AggregateID != data.AssetID ||
		(data.ActorID != "" && !helpers.ValidID(data.ActorID)) ||
		(data.PersonID != "" && !helpers.ValidID(data.PersonID)) ||
		data.AssetVersion < 0 || data.TenantID == "" ||
		strings.TrimSpace(data.TenantID) != data.TenantID || len(data.TenantID) > 128 ||
		data.OccurredAt.IsZero() || !event.OccurredAt.Equal(data.OccurredAt) ||
		!helpers.RegisteredOccurrence(data) {
		return event, ErrInvalidEvent
	}

	return event, nil
}
