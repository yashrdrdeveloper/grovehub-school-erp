package messaging

import (
    "context"
    "encoding/json"
    "errors"
    "log/slog"
    "time"

    client "grovehub/services/infrastructure/client/messaging"
    "grovehub/services/infrastructure/internal/app/entity"
    repository "grovehub/services/infrastructure/internal/app/repository/messaging"
)

// Dispatcher reliably publishes an event at least once. Consumers must
// deduplicate by event ID; a broker ack and PostgreSQL update aren't atomic.
type Dispatcher struct {
    store repository.Repository
    publisher client.MessagingClient
    logger *slog.Logger
}

func NewDispatcher(store repository.Repository,publisher client.MessagingClient,logger *slog.Logger) *Dispatcher {
    if logger==nil {logger=slog.Default()}
    return &Dispatcher{store:store,publisher:publisher,logger:logger}
}

func (dispatcher *Dispatcher) Dispatch(ctx context.Context) error {
    entries,err:=dispatcher.store.Claim(ctx,4,90*time.Second)
    if err!=nil{return err}
    var failures []error
    for _,entry:=range entries {
        if err:=dispatcher.deliver(ctx,entry);err!=nil{failures=append(failures,err)}
    }
    return errors.Join(failures...)
}

func (dispatcher *Dispatcher) deliver(ctx context.Context,entry entity.Messaging) error {
    metadata:=map[string]string{}
    err:=json.Unmarshal(entry.Metadata,&metadata)
    if err==nil{
        publishCtx,cancel:=context.WithTimeout(ctx,10*time.Second)
        err=dispatcher.publisher.Publish(publishCtx,entry.Topic,entry.Payload,metadata)
        cancel()
    }
    if err!=nil{
        dispatcher.logger.ErrorContext(ctx,"messaging outbox delivery failed","eventId",entry.EventID,"attempt",entry.Attempts,"error",err)
        if retryErr:=dispatcher.store.MarkRetry(ctx,entry.EventID,entry.Attempts);retryErr!=nil{return errors.Join(err,retryErr)}
        return err
    }
    return dispatcher.store.MarkSent(ctx,entry.EventID)
}

// Run dispatches once per interval and prunes SENT records after 24 hours.
// Pending/retrying events are never deleted and backoff remains durable.
func (dispatcher *Dispatcher) Run(ctx context.Context,interval time.Duration) {
    if interval<=0 {interval=time.Second}
    ticker:=time.NewTicker(interval)
    defer ticker.Stop()
    nextPrune:=time.Now().UTC().Add(time.Minute)
    for {
        if err:=dispatcher.Dispatch(ctx);err!=nil && ctx.Err()==nil {
            dispatcher.logger.ErrorContext(ctx,"messaging delivery batch failed","error",err)
        }
        if !time.Now().UTC().Before(nextPrune) && ctx.Err()==nil {
            if _,err:=dispatcher.store.PruneSent(ctx,time.Now().UTC().Add(-24*time.Hour),500);err!=nil {
                dispatcher.logger.ErrorContext(ctx,"messaging retention prune failed","error",err)
            }
            nextPrune=time.Now().UTC().Add(time.Minute)
        }
        select{
        case <-ctx.Done():return
        case <-ticker.C:
        }
    }
}
