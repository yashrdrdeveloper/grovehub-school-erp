// Package strategies declares the Foundation Person strategy contracts.
// Implementations belong to strategies/impl and domain decisions stay there.
package strategies

import (
    "context"

    servicepipeline "grovehub/libs/core/pipeline/service"
    persondto "grovehub/services/foundation/dto/person"
    personflow "grovehub/services/foundation/services/person/flows"
)

// CreateStrategy runs the four ordered steps of the Person creation use case.
// Each step returns Next, Return (with created flag), or Fail.
type CreateStrategy interface {
    Prepare(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person]
    Idempotency(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person]
    Contacts(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person]
    Persist(ctx context.Context, state *personflow.Create) servicepipeline.Outcome[persondto.Person]
}
