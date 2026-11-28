package events_test

import (
	"encoding/json"
	"errors"
	"testing"
	"time"

	coreevents "grovehub/libs/core/events"
	eventtypes "grovehub/services/audit/internal/app/constants/events"
	"grovehub/services/audit/internal/app/dto/producer"
	validator "grovehub/services/audit/internal/app/validators/events"
	"grovehub/services/audit/internal/app/validators/events/helpers"
)

func TestOccurrenceVersionIdentity(t *testing.T) {
	const asset = "e93696cc-15f0-4e06-9577-0d5cfc2a83c8"
	data := producer.Activity{
		TenantID: "school", Source: "foundation", AssetType: "PERSON", AssetID: asset,
		ActorID: asset, PersonID: asset, AssetVersion: 1, Operation: "CREATED",
		OccurredAt: time.Now().UTC(),
	}
	data.OccurrenceID = helpers.VersionEventID(data)
	event := coreevents.Envelope[producer.Activity]{
		EventID: data.OccurrenceID, EventType: eventtypes.ActivityOccurred, Version: 1,
		AggregateID: asset, OccurredAt: data.OccurredAt, Data: data,
	}

	raw, err := json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := validator.Occurrence(raw); err != nil {
		t.Fatal(err)
	}

	event.EventID = asset
	event.Data.OccurrenceID = asset
	raw, err = json.Marshal(event)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := validator.Occurrence(raw); !errors.Is(err, validator.ErrInvalidEvent) {
		t.Fatalf("invalid deterministic ID accepted: %v", err)
	}

	if _, err := validator.Occurrence([]byte("bad JSON")); !errors.Is(err, validator.ErrInvalidEvent) {
		t.Fatalf("malformed payload accepted: %v", err)
	}
}
