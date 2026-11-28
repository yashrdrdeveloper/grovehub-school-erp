// Package impl composes independent Messaging operations, like Foundation PersonService.
package impl

import (
    coreerrors "grovehub/libs/core/errors"
    client "grovehub/services/infrastructure/client/messaging"
    errorcode "grovehub/services/infrastructure/internal/app/constants/error"
    messagingservice "grovehub/services/infrastructure/internal/app/services/messaging"
    consume "grovehub/services/infrastructure/internal/app/services/messaging/operations/consume"
    consumeimpl "grovehub/services/infrastructure/internal/app/services/messaging/operations/consume/strategies/impl"
    health "grovehub/services/infrastructure/internal/app/services/messaging/operations/health"
    healthimpl "grovehub/services/infrastructure/internal/app/services/messaging/operations/health/strategies/impl"
    publish "grovehub/services/infrastructure/internal/app/services/messaging/operations/publish"
    publishimpl "grovehub/services/infrastructure/internal/app/services/messaging/operations/publish/strategies/impl"
    "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
)

// @Service
type MessagingService struct {
    *publish.PublishOperation
    *consume.ConsumeOperation
    *health.HealthOperation
    provider strategies.Strategy
}

var _ messagingservice.Service = (*MessagingService)(nil)

func NewMessagingService(provider string, registry *strategies.Registry[strategies.Strategy], store client.EventStore) (messagingservice.Service, error) {
    transport, err := registry.Resolve(provider)
    if err != nil { return nil, coreerrors.Wrap(errorcode.MessagingProviderUnavailable, err) }
    return &MessagingService{
        PublishOperation: publish.New(publishimpl.New(transport, store)),
        ConsumeOperation: consume.New(consumeimpl.New(transport)),
        HealthOperation: health.New(healthimpl.New(transport)),
        provider: transport,
    }, nil
}

func (service *MessagingService) Close() error { return service.provider.Close() }
