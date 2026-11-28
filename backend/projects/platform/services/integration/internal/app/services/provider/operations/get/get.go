package get

import (
    "context"
    dto "grovehub/services/integration/internal/app/dto/provider"
    "grovehub/services/integration/internal/app/services/provider/operations/get/strategies"
)

type GetOperation struct {strategy strategies.Strategy}

func New(strategy strategies.Strategy) *GetOperation {return &GetOperation{strategy:strategy}}

func (service *GetOperation) Provider(ctx context.Context, id string) (dto.ProviderDetails,error) {return service.strategy.Execute(ctx,id)}
