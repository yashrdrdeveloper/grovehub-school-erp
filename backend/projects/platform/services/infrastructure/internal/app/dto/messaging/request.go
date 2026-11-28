package messaging

// Event is a transport-neutral, immutable pending event DTO.
type Event struct {
	EventID string
	Topic string
	Payload []byte
	Metadata map[string]string
}

type PublishOptions struct {
	Persist bool
}
