package outbox

import (
    "context"
    "encoding/json"
    "errors"
    "log/slog"
    "time"

    coremessaging "grovehub/libs/core/events/messaging"
)

// Dispatcher uses at-least-once delivery: downstream handlers must deduplicate event IDs.
type Dispatcher struct {
    store Store
    publisher coremessaging.Publisher
    logger *slog.Logger
}

func NewDispatcher(store Store, publisher coremessaging.Publisher, logger *slog.Logger) *Dispatcher {
    if logger == nil { logger = slog.Default() }
    return &Dispatcher{store: store, publisher: publisher, logger: logger}
}

func (dispatcher *Dispatcher) Dispatch(ctx context.Context) error {
    // Four sequential ten-second publish attempts fit inside the ninety-second lease.
    entries, err := dispatcher.store.Claim(ctx, 4, 90*time.Second)
    if err != nil { return err }
    var failures []error
    for _, entry := range entries {
        if err := dispatcher.deliver(ctx, entry); err != nil { failures = append(failures, err) }
    }
    return errors.Join(failures...)
}

func (dispatcher *Dispatcher) deliver(ctx context.Context, entry Entry) error {
    metadata := map[string]string{}
    err := json.Unmarshal(entry.Metadata, &metadata)
    if err == nil {
        publishCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
        err = dispatcher.publisher.Publish(publishCtx, entry.Topic, entry.Payload, metadata)
        cancel()
    }
    if err != nil {
        dispatcher.logger.ErrorContext(ctx, "outbox delivery failed", "eventId", entry.EventID, "attempt", entry.Attempts, "error", err)
        if retryErr := dispatcher.store.MarkRetry(ctx, entry.EventID, entry.Attempts); retryErr != nil { return errors.Join(err, retryErr) }
        return err
    }
    return dispatcher.store.MarkSent(ctx, entry.EventID)
}

// Run polls until cancellation; callers own the dispatcher's lifetime.
func (dispatcher *Dispatcher) Run(ctx context.Context, interval time.Duration) {
    if interval <= 0 { interval = time.Second }
    ticker := time.NewTicker(interval)
    defer ticker.Stop()
    for {
        if err := dispatcher.Dispatch(ctx); err != nil && ctx.Err() == nil {
            dispatcher.logger.ErrorContext(ctx, "outbox batch failed", "error", err)
        }
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
        }
    }
}
