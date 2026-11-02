package config

import "time"

const ProviderNATS = "nats"

type Messaging struct {
	Provider         string
	NATSURL          string
	QueueGroupPrefix string
	DurablePrefix    string
	CloseTimeout     time.Duration
	AckWaitTimeout   time.Duration
	SubscribeTimeout time.Duration
}

func DefaultMessaging() Messaging {
	return Messaging{
		Provider:         ProviderNATS,
		NATSURL:          "nats://127.0.0.1:4222",
		QueueGroupPrefix: "grovehub",
		DurablePrefix:    "grovehub",
		CloseTimeout:     30 * time.Second,
		AckWaitTimeout:   30 * time.Second,
		SubscribeTimeout: 30 * time.Second,
	}
}
