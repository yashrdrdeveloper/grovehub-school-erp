package helpers

import (
	coreevents "grovehub/libs/core/events"
	foundationresponse "grovehub/services/foundation/dto/response"
	personevents "grovehub/services/foundation/events/persons"
)

func BuildCreated(eventID string, person foundationresponse.Person) coreevents.Envelope[personevents.Created] {
	return coreevents.Envelope[personevents.Created]{
		EventID:    eventID,
		EventType:  personevents.CreatedType,
		Source:     personevents.Source,
		OccurredAt: person.CreatedAt,
		Data: personevents.Created{
			GID:     person.GID,
			Status:  person.Status,
			Version: person.Version,
		},
	}
}

func CreatedMetadata(eventID string) map[string]string {
	return map[string]string{
		"event-id":   eventID,
		"event-type": personevents.CreatedType,
		"source":     personevents.Source,
	}
}
