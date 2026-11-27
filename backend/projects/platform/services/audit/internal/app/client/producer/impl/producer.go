package impl

import (
	client "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/services/producer"
)

// Event resolves the public contract's unqualified DTO in generated proxy code.
type Event = client.Event
type Activity = client.Activity

type producerClient struct { service producer.Service }

func NewProducerClient(service producer.Service) client.Client {
	return &producerClient{service:service}
}
