package strategies

import (
    "context"

    persondto "grovehub/services/foundation/dto/person"
)

// @Strategy
type GetByGIDStrategy interface {
    Execute(ctx context.Context, request persondto.GetPersonByGIDRequest) (persondto.Person, error)
}
