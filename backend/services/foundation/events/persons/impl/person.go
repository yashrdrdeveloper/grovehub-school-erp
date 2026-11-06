package impl

import (
    "context"
    "encoding/json"
    "log/slog"

    "github.com/google/uuid"

 coretransaction "grovehub/libs/core/database/transaction"

    coresearch "grovehub/libs/core/search"
    foundationevents "grovehub/services/foundation/constants/events"
    foundationrequest "grovehub/services/foundation/dto/request"
    foundationresponse "grovehub/services/foundation/dto/response"
    personevents "grovehub/services/foundation/events/persons"
    eventhelpers "grovehub/services/foundation/events/persons/helpers"
    foundationservices "grovehub/services/foundation/services"
    infrastructureclient "grovehub/services/infrastructure/client"
)

type PersonService struct {
    delegate foundationservices.PersonService
    messagingClient infrastructureclient.MessagingClient
    logger *slog.Logger
}

func NewPersonService(delegate foundationservices.PersonService, messagingClient infrastructureclient.MessagingClient, logger *slog.Logger) foundationservices.PersonService {
    return &PersonService{delegate: delegate, messagingClient: messagingClient, logger: logger}
}

func (service *PersonService) Create(ctx context.Context, request foundationrequest.CreatePerson) (foundationresponse.Person, bool, error) {
    person, created, err := service.delegate.Create(ctx, request)
    if err != nil || !created { return person, created, err }
    service.publishChanged(ctx, foundationevents.Created, nil, &person)
    return person, true, nil
}

func (service *PersonService) GetByGID(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error) {
    return service.delegate.GetByGID(ctx, request)
}

func (service *PersonService) Update(ctx context.Context, personRequest foundationrequest.PersonByGID, request foundationrequest.UpdatePerson) (foundationresponse.Person, error) {
    old, err := service.delegate.GetByGID(ctx, personRequest)
    if err != nil { return foundationresponse.Person{}, err }
    next, err := service.delegate.Update(ctx, personRequest, request)
    if err != nil { return next, err }
    if next.Version != old.Version {
        service.publishChanged(ctx, foundationevents.Updated, &old, &next)
    }
    return next, nil
}

func (service *PersonService) Search(ctx context.Context, request coresearch.Request) (coresearch.Result[foundationresponse.Person], error) {
    return service.delegate.Search(ctx, request)
}

func (service *PersonService) Deactivate(ctx context.Context, request foundationrequest.PersonByGID) (foundationresponse.Person, error) {
    old, err := service.delegate.GetByGID(ctx, request)
    if err != nil { return foundationresponse.Person{}, err }
    next, err := service.delegate.Deactivate(ctx, request)
    if err != nil { return next, err }
    if next.Version != old.Version {
        service.publishChanged(ctx, foundationevents.Deactivated, &old, &next)
    }
    return next, nil
}

func (service *PersonService) publishChanged(ctx context.Context, operation foundationevents.Operation, old, next *foundationresponse.Person) {
    if coretransaction.AfterCommit(ctx, func() {
        service.publishChanged(coretransaction.DetachedContext(ctx), operation, old, next)
    }) { return }
    eventID := uuid.NewString()
    envelope := eventhelpers.BuildChanged(ctx, eventID, operation, old, next)
    payload, err := json.Marshal(envelope)
    if err != nil {
        service.logger.ErrorContext(ctx, "failed to marshal foundation person changed event", "operation", operation, "error", err)
        return
    }
    if err := service.messagingClient.Publish(ctx, personevents.ChangedTopic, payload, eventhelpers.ChangedMetadata(eventID)); err != nil {
        // The database operation has already succeeded; publishing is best-effort
        // until a transactional outbox is introduced.
        service.logger.ErrorContext(ctx, "failed to publish foundation person changed event", "eventId", eventID, "operation", operation, "error", err)
    }
}
