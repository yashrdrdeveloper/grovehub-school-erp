package contracts

import (
    "context"
    "encoding/json"
    dto "grovehub/services/integration/internal/app/providers/dto"
)

type ActionHandler interface {
    Execute(ctx context.Context, session dto.Session, input json.RawMessage) (any,error)
}

type ActionRegistry interface {
    Register(provider string, action dto.Action, handler ActionHandler) error
    Execute(ctx context.Context, session dto.Session, action dto.Action, input json.RawMessage) (any,error)
    Actions(provider string) []dto.Action
}
