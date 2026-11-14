// Package handlers decodes and validates Person Changed events for inbound subscriptions.
package handlers

import (
    "context"
    "encoding/json"
    "errors"

    coreevents "grovehub/libs/core/events"
    personevents "grovehub/services/foundation/events/person"
)

var ErrInvalidChangedEvent = errors.New("invalid foundation Person Changed event")

type ChangedHandler func(context.Context, coreevents.Envelope[personevents.Changed]) error

func DecodeChanged(ctx context.Context, payload []byte, handler ChangedHandler) error {
    var envelope coreevents.Envelope[personevents.Changed]
    if err := json.Unmarshal(payload, &envelope); err != nil { return err }
    if envelope.EventID == "" || envelope.AggregateID == "" || envelope.EventType != personevents.ChangedType ||
        envelope.Version != personevents.ChangedVersion { return ErrInvalidChangedEvent }
    return handler(ctx, envelope)
}
