package action

import (
    providerdto "grovehub/services/integration/internal/app/providers/dto"
    "context"
    "encoding/json"
    "grovehub/services/integration/internal/app/services/connection/operations/action/strategies"
)

type ActionOperation struct{ strategy strategies.Strategy }

func New(strategy strategies.Strategy) *ActionOperation {return &ActionOperation{strategy:strategy}}

func (operation *ActionOperation) ExecuteAction(ctx context.Context, tenantID, connectionID string, action providerdto.Action, input json.RawMessage) (any,error) {return operation.strategy.ExecuteAction(ctx,tenantID,connectionID,action,input)}
