package strategies

import (
    "context"

    servicepipeline "grovehub/libs/core/pipeline/service"
    persondto "grovehub/services/foundation/dto/person"
    updateflow "grovehub/services/foundation/services/person/operations/update/flows"
)

// @Strategy
type UpdateStrategy interface {
    Load(ctx context.Context, state *updateflow.Update) servicepipeline.Outcome[persondto.Person]
    Apply(ctx context.Context, state *updateflow.Update) servicepipeline.Outcome[persondto.Person]
    Persist(ctx context.Context, state *updateflow.Update) servicepipeline.Outcome[persondto.Person]
}
