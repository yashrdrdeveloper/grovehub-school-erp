package strategies

import (
    "context"

    corepipeline "grovehub/libs/core/pipeline"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    updateflow "grovehub/services/foundation/internal/app/services/person/operations/update/flows"
)

// @Strategy
type UpdateStrategy interface {
    Load(ctx context.Context, state *updateflow.Update) corepipeline.Outcome[persondto.Person]
    Apply(ctx context.Context, state *updateflow.Update) corepipeline.Outcome[persondto.Person]
    Persist(ctx context.Context, state *updateflow.Update) corepipeline.Outcome[persondto.Person]
}
