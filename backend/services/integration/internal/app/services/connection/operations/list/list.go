package list

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/services/connection/operations/list/strategies"
)

type ListOperation struct{strategy strategies.Strategy}

func New(strategy strategies.Strategy) *ListOperation {return &ListOperation{strategy:strategy}}

func (operation *ListOperation) List(ctx context.Context, tenantID, ownerGID string) ([]connectiondto.Connection,error) {return operation.strategy.Execute(ctx,tenantID,ownerGID)}
