package test

import (
    "context"
    connectiondto "grovehub/services/integration/internal/app/dto/connection"
    "grovehub/services/integration/internal/app/services/connection/operations/test/strategies"
)

type TestOperation struct{ strategy strategies.Strategy }

func New(strategy strategies.Strategy) *TestOperation {return &TestOperation{strategy:strategy}}

func (operation *TestOperation) RunTest(ctx context.Context, tenantID, id, direction string) (connectiondto.ConnectionTest,error) {return operation.strategy.RunTest(ctx,tenantID,id,direction)}

func (operation *TestOperation) ListTests(ctx context.Context, tenantID, id string) ([]connectiondto.ConnectionTest,error) {return operation.strategy.ListTests(ctx,tenantID,id)}
