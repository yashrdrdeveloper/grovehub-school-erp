// Package consumers declares inbound Person event subscription contracts.
package consumers

import "context"

type ChangedConsumer interface {
    Register(ctx context.Context) error
}
