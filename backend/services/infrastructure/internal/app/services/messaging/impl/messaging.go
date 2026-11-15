package impl

import (
    "context"
    "strings"

    "github.com/google/uuid"

    coreevents "grovehub/libs/core/events"
    coremessaging "grovehub/libs/core/events/messaging"
    messagingclient "grovehub/services/infrastructure/client/messaging"
    coreerrors "grovehub/libs/core/errors"
    errorcode "grovehub/services/infrastructure/internal/app/constants/error"
    messagingdto "grovehub/services/infrastructure/internal/app/dto/messaging"
    infrastructureservices "grovehub/services/infrastructure/internal/app/services/messaging"
    "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
)

type messagingService struct {
    provider string
    registry *strategies.Registry[strategies.Strategy]
    store messagingclient.EventStore
}

func NewMessagingService(provider string, registry *strategies.Registry[strategies.Strategy], store messagingclient.EventStore) infrastructureservices.Service {
    return &messagingService{provider: provider,registry: registry,store:store}
}

func (service *messagingService) Publish(ctx context.Context, topic string, payload []byte, metadata map[string]string) error {
    provider,err:=service.registry.Resolve(service.provider)
    if err!=nil{return coreerrors.Wrap(errorcode.MessagingProviderUnavailable,err)}
    if err:=provider.Publish(ctx,topic,payload,metadata);err!=nil{
        return coreerrors.Wrap(errorcode.MessagingPublishFailed,err)
    }
    return nil
}

// PublishWithOptions routes best-effort messages directly to NATS, while
// persistent messages are written to the transactional outbox for the worker
// to deliver after commit. Persist=true must use the same transaction context
// as the business write when atomic delivery is required.
func (service *messagingService) PublishWithOptions(ctx context.Context, topic string, payload []byte, metadata map[string]string, options coremessaging.PublishOptions) error {
    if !options.Persist {
        return service.Publish(ctx,topic,payload,metadata)
    }
    if service.store == nil {
        return errorcode.ErrInvalidStore
    }

    headers:=make(map[string]string,len(metadata)+1)
    for key,value:=range metadata {
        headers[key]=value
    }
    eventID:=strings.TrimSpace(headers["event-id"])
    if eventID=="" {
        eventID=uuid.NewString()
    }
    headers["event-id"]=eventID

    return service.store.Enqueue(ctx,coreevents.Message{
        EventID:eventID,
        Topic:topic,
        Payload:append([]byte(nil),payload...),
        Metadata:headers,
    })
}

func (service *messagingService) Consume(ctx context.Context, topic string, handler messagingdto.MessageHandler) error {
    provider,err:=service.registry.Resolve(service.provider)
    if err!=nil{return coreerrors.Wrap(errorcode.MessagingProviderUnavailable,err)}
    if err:=provider.Consume(ctx,topic,handler);err!=nil{
        return coreerrors.Wrap(errorcode.MessagingConsumeFailed,err)
    }
    return nil
}

func (service *messagingService) Health(ctx context.Context) error {
    provider,err:=service.registry.Resolve(service.provider)
    if err!=nil{return coreerrors.Wrap(errorcode.MessagingProviderUnavailable,err)}
    return provider.Health(ctx)
}

func (service *messagingService) Close() error {
    provider,err:=service.registry.Resolve(service.provider)
    if err!=nil{return coreerrors.Wrap(errorcode.MessagingProviderUnavailable,err)}
    return provider.Close()
}
