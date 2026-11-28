package impl

import (
	"context"
	coreerrors "grovehub/libs/core/errors"
	errorcode "grovehub/services/infrastructure/internal/app/constants/error"
	dto "grovehub/services/infrastructure/internal/app/dto/messaging"
	consumestrategies "grovehub/services/infrastructure/internal/app/services/messaging/operations/consume/strategies"
	providerstrategies "grovehub/services/infrastructure/internal/app/services/messaging/strategies"
)

type consumeStrategy struct{ provider providerstrategies.Strategy }

var _ consumestrategies.Strategy = (*consumeStrategy)(nil)

func New(provider providerstrategies.Strategy) consumestrategies.Strategy {
	return &consumeStrategy{provider: provider}
}
func (strategy *consumeStrategy) Execute(ctx context.Context, topic string, handler dto.MessageHandler) error {
	if err := strategy.provider.Consume(ctx, topic, handler); err != nil {
		return coreerrors.Wrap(errorcode.MessagingConsumeFailed, err)
	}
	return nil
}
