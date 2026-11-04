package events

// Operation describes the Foundation aggregate lifecycle change in an event.
// Keep this domain-owned rather than adding it to generic libs/core.
type Operation string

const (
    Created     Operation = "CREATED"
    Updated     Operation = "UPDATED"
    Deactivated Operation = "DEACTIVATED"
)
