// Package contracts owns the provider-neutral messaging transport types.
package contracts

import "context"

type Event struct {
	EventID string
	Topic string
	Payload []byte
	Metadata map[string]string
}

type Message struct {
	ID string
	Payload []byte
	Metadata map[string]string
}

type Handler func(context.Context, Message) error

type PublishOptions struct {
	Persist bool
}
