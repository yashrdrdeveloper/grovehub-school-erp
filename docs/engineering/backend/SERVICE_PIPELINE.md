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
- backend/services/foundation/services/person/flows/create.go: typed per-request state
- backend/services/foundation/services/person/strategies/create/: independently named business strategies
- backend/services/foundation/services/person/impl/person.go: thin Create entry point

## Annotation semantics
- @Service marks the service type; registration still occurs explicitly in foundation.New.
- @ChainOfResponsibility marks the Create method; Run executes the chain explicitly.
- @Strategy marks a typed strategy implementation; the constructor supplies repositories.
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
