package impl

import (
	client "grovehub/services/audit/client/producer"
	"grovehub/services/audit/internal/app/services/producer"
)

type producerClient struct { service producer.Service }

func NewProducerClient(service producer.Service) client.Client {
	return &producerClient{service:service}
}
