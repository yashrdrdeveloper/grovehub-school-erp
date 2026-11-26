package producer

import (
	"context"
	client "grovehub/services/audit/client/producer"
)

type Service interface {
	Produce(ctx context.Context, event client.Event) error
}
