package impl

import (
	client "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/services/producer"
)

// producerClient forwards publication through the Audit-owned service and
// preserves the caller's transaction context.
type producerClient struct {
	service producer.Service
}

func NewProducerClient(service producer.Service) client.Client {
	return &producerClient{service: service}
}
