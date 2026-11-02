package helpers

import (
	coreevents "grovehub/libs/core/events"
	"grovehub/services/identity/entity"
	accountevents "grovehub/services/identity/events/accounts"
)

func BuildCreated(eventID string, account entity.Account) coreevents.Envelope[accountevents.Created] {
	return coreevents.Envelope[accountevents.Created]{
		EventID:    eventID,
		EventType:  accountevents.CreatedType,
		Source:     accountevents.Source,
		OccurredAt: account.GetCreatedAt(),
		Data: accountevents.Created{
			AccountID: string(account.GetAccountId()),
			PersonGID: string(account.GetPersonGid()),
			Status:    string(account.GetStatus()),
			Version:   account.GetVersion(),
		},
	}
}

func CreatedMetadata(eventID string) map[string]string {
	return map[string]string{
		"event-id":   eventID,
		"event-type": accountevents.CreatedType,
		"source":     accountevents.Source,
	}
}
