package impl

import (
    "context"
    "errors"
    "time"

    providerconstants "grovehub/services/integration/internal/app/constants/providers"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/entity"
    mailboxevents "grovehub/services/integration/internal/app/events/producers/connection/operations/mailbox"
    changedproducer "grovehub/services/integration/internal/app/events/producers/connection"
    "grovehub/services/integration/internal/app/externalService/mail/gmail"
    connectionrepo "grovehub/services/integration/internal/app/repository/connection"
    watchrepo "grovehub/services/integration/internal/app/repository/watch"
    connectionhelpers "grovehub/services/integration/internal/app/services/connection/helpers"
    oauth "grovehub/services/integration/internal/app/services/connection/operations/oauth"
    watchstrategies "grovehub/services/integration/internal/app/services/connection/operations/watch/strategies"
)

type watchStrategy struct {
    repository connectionrepo.ConnectionRepository
    watchRepository watchrepo.WatchRepository
    gmail *gmail.Client
    pushTopic string
    events changedproducer.ChangedProducer
    tokens oauth.TokenProvider
}
var _ watchstrategies.Strategy = (*watchStrategy)(nil)
func New(repository connectionrepo.ConnectionRepository, watches watchrepo.WatchRepository, gmailClient *gmail.Client, topic string, events changedproducer.ChangedProducer, tokens oauth.TokenProvider) watchstrategies.Strategy {
    return &watchStrategy{repository:repository,watchRepository:watches,gmail:gmailClient,pushTopic:topic,events:events,tokens:tokens}
}
func (s *watchStrategy) Watch(ctx context.Context, tenantID, id string) (connectiondto.Watch, error) {
	if s.gmail == nil || s.pushTopic == "" {
		return connectiondto.Watch{}, connectionhelpers.ErrPushUnavailable
	}
	c, err := s.repository.Get(ctx, tenantID, id)
	if err != nil {
		return connectiondto.Watch{}, err
	}
	if c.Provider != providerconstants.Gmail || c.Status != entity.StatusActive {
		return connectiondto.Watch{}, connectionhelpers.ErrConnectionInactive
	}
	token, err := s.tokens.ValidProviderToken(ctx, id, "GMAIL")
	if err != nil {
		return connectiondto.Watch{}, err
	}
	historyID, expiration, err := s.gmail.Watch(ctx, token, s.pushTopic)
	if err != nil {
		return connectiondto.Watch{}, err
	}
	if err = s.watchRepository.SaveWatch(ctx, id, historyID, expiration); err != nil {
		return connectiondto.Watch{}, err
	}
	return connectiondto.Watch{ConnectionID: id, ExpiresAt: expiration}, nil
}
func (s *watchStrategy) ProcessPush(ctx context.Context, notification gmail.PushNotification) error {
	if s.gmail == nil || s.pushTopic == "" {
		return connectionhelpers.ErrPushUnavailable
	}
	matches, err := s.watchRepository.GetGmailByMailbox(ctx, notification.EmailAddress)
	if err != nil {
		return err
	}
	// No message ingestion or cursor updates occur here. Communication will
	// consume the event and call Integration provider capabilities.
	for _, connection := range matches {
		if s.events == nil {
			return errors.New("integration messaging is unavailable")
		}
		notificationData := mailboxevents.Payload{
			Provider:     providerconstants.Gmail,
			ConnectionID: connection.ID,
			TenantID:     connection.TenantID,
			OwnerGID:     connection.OwnerGID,
			Mailbox:      notification.EmailAddress,
			HistoryHint:  notification.HistoryID,
		}
		if err := s.events.PublishChanged(ctx, notificationData); err != nil {
			return err
		}
	}
	return nil
}
func (s *watchStrategy) RenewExpiring(ctx context.Context) error {
	if s.gmail == nil || s.pushTopic == "" {
		return nil
	}
	rows, err := s.watchRepository.ExpiringWatches(ctx, time.Now().UTC().Add(24*time.Hour))
	if err != nil {
		return err
	}
	for _, c := range rows {
		if _, err := s.Watch(ctx, c.TenantID, c.ID); err != nil {
			return err
		}
	}
	return nil
}

