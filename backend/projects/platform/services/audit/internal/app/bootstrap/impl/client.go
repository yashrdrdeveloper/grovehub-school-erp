package impl

import (
	client "grovehub/services/audit/client/producer"
	clientimpl "grovehub/services/audit/internal/app/client/producer/impl"
	serviceimpl "grovehub/services/audit/internal/app/services/producer/impl"
	"grovehub/services/infrastructure/client/messaging"
)

func NewProducer(store messaging.EventStore) client.Client {
	return clientimpl.NewProducerClient(serviceimpl.New(store))
}
