package strategies

import (
    "context"

    persondto "grovehub/services/foundation/client/person/dto"
)

// @Strategy
type GetByGIDStrategy interface {
    Execute(ctx context.Context, request persondto.GetPersonByGIDRequest) (persondto.Person, error)
}
