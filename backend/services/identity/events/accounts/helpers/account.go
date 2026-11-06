package helpers

import (
    "context"

    coreevents "grovehub/libs/core/events"
    "grovehub/services/identity/entity"
    accountevents "grovehub/services/identity/events/accounts"
)

// BuildCreated wraps the Account payload in the shared versioned envelope.
func BuildCreated(ctx context.Context, eventID string, account entity.Account) coreevents.Envelope[accountevents.Created] {
    return coreevents.NewEnvelope(ctx, coreevents.Metadata{
        EventID:     eventID,
        EventType:   accountevents.CreatedType,
        Version:     1,
        AggregateID: string(account.GetAccountId()),
        OccurredAt:  account.GetCreatedAt(),
    }, accountevents.Created{
        AccountID: string(account.GetAccountId()),
        PersonGID: string(account.GetPersonGid()),
        Status:    string(account.GetStatus()),
        Version:   account.GetVersion(),
    })
}

func CreatedMetadata(eventID string) map[string]string {
    return map[string]string{
        "event-id":   eventID,
        "event-type": accountevents.CreatedType,
        "source":     accountevents.Source,
    }
}
