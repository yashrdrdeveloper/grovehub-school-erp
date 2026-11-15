package helpers

import (
    "strings"
    "github.com/google/uuid"
    inframessaging "grovehub/libs/infrastructure/messaging"
)

// NewMessage defensively copies the immutable event and preserves an explicit
// event-id; absent IDs are generated before persistence.
func NewMessage(topic string, payload []byte, metadata map[string]string) inframessaging.Event {
    headers := make(map[string]string, len(metadata)+1)
    for key, value := range metadata { headers[key] = value }
    eventID := strings.TrimSpace(headers["event-id"])
    if eventID == "" { eventID = uuid.NewString() }
    headers["event-id"] = eventID
    return inframessaging.Event{
        EventID:eventID,
        Topic:topic,
        Payload:append([]byte(nil), payload...),
        Metadata:headers,
    }
}
