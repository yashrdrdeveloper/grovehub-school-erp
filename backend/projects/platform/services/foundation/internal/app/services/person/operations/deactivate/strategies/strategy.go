package strategies

import (
    "context"

    corepipeline "grovehub/libs/core/pipeline"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    deactivateflow "grovehub/services/foundation/internal/app/services/person/operations/deactivate/flows"
)

// @Strategy
type DeactivateStrategy interface {
    Load(ctx context.Context, state *deactivateflow.Deactivate) corepipeline.Outcome[persondto.Person]
    Prepare(ctx context.Context, state *deactivateflow.Deactivate) corepipeline.Outcome[persondto.Person]
    Persist(ctx context.Context, state *deactivateflow.Deactivate) corepipeline.Outcome[persondto.Person]
}
