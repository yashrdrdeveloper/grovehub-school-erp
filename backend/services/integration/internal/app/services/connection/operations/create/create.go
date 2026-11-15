package create

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/services/connection/operations/create/strategies"
)

type CreateOperation struct{strategy strategies.Strategy}
func New(strategy strategies.Strategy) *CreateOperation {return &CreateOperation{strategy:strategy}}
func (operation *CreateOperation) Create(ctx context.Context, tenantID string, input connectiondto.CreateConnection) (connectiondto.Connection,error) {return operation.strategy.Execute(ctx,tenantID,input)}
