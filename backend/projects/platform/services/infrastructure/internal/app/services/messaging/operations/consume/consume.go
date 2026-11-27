package consume

import (
    "context"
    dto "grovehub/services/infrastructure/internal/app/dto/messaging"
    "grovehub/services/infrastructure/internal/app/services/messaging/operations/consume/strategies"
)

type ConsumeOperation struct { strategy strategies.Strategy }
func New(strategy strategies.Strategy) *ConsumeOperation { return &ConsumeOperation{strategy:strategy} }
func (operation *ConsumeOperation) Consume(ctx context.Context, topic string, handler dto.MessageHandler) error {
    return operation.strategy.Execute(ctx, topic, handler)
}
