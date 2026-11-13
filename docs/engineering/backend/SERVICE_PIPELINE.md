# GrooveHub service pipeline: Strategy + Chain of Responsibility

Foundation Person creation uses a typed, ordered chain:
Prepare -> Idempotency -> Contacts -> Persist.

The controller validates/normalizes the DTO once. The service does **not**
call CreateRules again. Every handler returns Next, Return(value, created), or
Fail(error). Run stops immediately on Return or Fail and rejects chains that
finish without a response.

## Layout
- backend/libs/core/pipeline/service/{chain,strategy,transaction}.go: runner, predicates, outcomes, transactional execution and post-rollback recovery
- backend/libs/core/annotation/{chain,strategy}: directive metadata
- backend/services/foundation/services/person/operations/create/create.go: operation-owned Create method and chain orchestration
- backend/services/foundation/services/person/operations/create/flows/create.go: typed Create request state
- backend/services/foundation/services/person/operations/create/strategies/strategy.go: CreateStrategy interface; service depends on this contract
- backend/services/foundation/services/person/operations/create/strategies/impl/: Create-only business strategy implementations
- backend/services/foundation/services/person/operations/create/helpers/: Create-only Person construction and error mapping helpers
- backend/services/foundation/services/person/helpers/: Person utilities shared between Create and other operations
- backend/libs/core/helpers/string.go: generic nullable/optional string conversions
- backend/services/foundation/services/person/impl/person.go: composing service embeds the Create operation and retains other Person methods

## Annotation semantics
- @Service marks the service type; registration still occurs explicitly in foundation.New.
- @ChainOfResponsibility marks `operations/create/create.go` method `(*Operation).Create`; Run executes the chain explicitly.
- @Strategy marks both the CreateStrategy interface in operations/create/strategies/strategy.go and its concrete Creator implementation in operations/create/strategies/impl/create.go; NewCreateStrategy injects repositories and returns the interface. The annotation checker recognizes both targets, but actual strategy wiring remains explicit.
- @Transactional marks the persistWrites method; the **actual** transaction is
  enforced by Transactional(unitOfWork, persistWrites), not by the comment.

Go has no runtime annotation interception. These are currently pilot, checked
metadata; do not remove the explicit runner/wrapper without introducing and
verifying a code generator.

## Correctness notes
- Person and Contact are written atomically through the existing UnitOfWork.
- Duplicate idempotency key and contact conflicts are re-read only after the
  failed transaction exits; this is necessary on PostgreSQL.
- A replay with a different request fingerprint returns a conflict.
- Contact matches to multiple Person records return a contact conflict.
- Neither lookup nor policy selection assumes that concurrent writes cannot occur:
  database uniqueness constraints remain authoritative.
- Other Foundation Person operations retain their existing behavior; this
  migration targets the Create workflow first.

## Reusable helper ownership
`OptionalString` and `NullableString` live in `libs/core/helpers`; the Person
service, mapper and idempotency implementation all use these conversions.
Person-specific mapping, contact ownership, and search projection remain in
Foundation. Redundant pass-through cursor encoding has been removed in favor
of direct `coresearch.EncodeCursor` calls.

## Use-case ownership
Person Create owns its `flows`, `strategies`, and `helpers` beneath `services/person/operations/create`.
Other Person operations continue using `services/person/helpers` for shared contact
normalization, association, comparison, and search projection. `NewPerson` and
Create-only error outcome helpers live in `operations/create/helpers`, not the shared
Person helper package.

## Operation entry points
The `PersonService` in `services/person/impl` anonymously embeds
`*create.Operation`. Go promotes its `Create` method onto `PersonService`, so
`foundationservices.PersonService` keeps the same method signature without
an extra forwarding method. The strategy interface and implementations stay
inside `operations/create/strategies`. A method cannot declare a receiver type
from another Go package; `Create` therefore belongs to `*create.Operation`.
