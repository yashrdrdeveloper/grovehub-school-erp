package strategies

import (
    "context"

    servicepipeline "grovehub/libs/core/pipeline/service"
    persondto "grovehub/services/foundation/dto/person"
    deactivateflow "grovehub/services/foundation/services/person/operations/deactivate/flows"
)

// @Strategy
type DeactivateStrategy interface {
    Load(ctx context.Context, state *deactivateflow.Deactivate) servicepipeline.Outcome[persondto.Person]
    Prepare(ctx context.Context, state *deactivateflow.Deactivate) servicepipeline.Outcome[persondto.Person]
    Persist(ctx context.Context, state *deactivateflow.Deactivate) servicepipeline.Outcome[persondto.Person]
}
