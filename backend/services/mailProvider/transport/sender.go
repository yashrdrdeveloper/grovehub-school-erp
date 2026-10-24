package transport

import (
	"context"

	"grovehub/services/mailProvider/entity"
)

type Sender interface {
	Send(ctx context.Context, message entity.Message, sender entity.Sender) error
}
