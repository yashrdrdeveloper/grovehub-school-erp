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
        Version:     accountevents.CreatedVersion,
        AggregateID: string(account.AccountId),
        OccurredAt:  account.CreatedAt,
    }, accountevents.Created{
        AccountID: string(account.AccountId),
        PersonGID: string(account.PersonGid),
        Status:    string(account.Status),
        Version:   account.Version,
    })
}

func CreatedMetadata(eventID string) map[string]string {
    return map[string]string{
        "event-id":   eventID,
        "event-type": accountevents.CreatedType,
        "source":     accountevents.Source,
    }
}
