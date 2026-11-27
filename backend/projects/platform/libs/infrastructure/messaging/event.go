package messaging

// Event is a transport-neutral, immutable pending outbox message.
// Domain payloads, including complete old/new entity snapshots, stay opaque.
type Event struct {
    EventID string
    Topic string
    Payload []byte
    Metadata map[string]string
}
