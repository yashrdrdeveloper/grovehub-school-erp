// Package strategies declares Person Create strategy contracts.
// Implementations belong to strategies/impl and domain decisions stay there.
package strategies

import (
    "context"

    corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
    persondto "grovehub/services/foundation/internal/app/dto/person"
    personflow "grovehub/services/foundation/internal/app/services/person/operations/create/flows"
)

// @Strategy
// CreateStrategy runs the four ordered steps of the Person creation use case.
// Each step returns Next, Return (with created flag), or Fail.
type CreateStrategy interface {
    Prepare(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person]
    Idempotency(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person]
    Contacts(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person]
    Persist(ctx context.Context, state *personflow.Create) corepipeline.Outcome[persondto.Person]
}
