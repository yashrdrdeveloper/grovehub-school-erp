package events

// Message is an immutable event envelope payload for transactional delivery.
// It contains no database logic and is shared across all service event rules.
type Message struct {
    EventID string
    Topic string
    Payload []byte
    Metadata map[string]string
}
