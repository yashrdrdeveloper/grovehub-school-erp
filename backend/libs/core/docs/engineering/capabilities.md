# Core Library — Capability Ownership

This inventory describes the current **Core-owned technical packages**. Their root package is the stable import API where provided. Internal contracts/validators/helpers under a capability must not become arbitrary imports from services.

| Capability path under `libs/backend/core/` | Responsibility | Must not contain |
|---|---|---|
| `annotation/` | Registered syntax and typed Go code generation/checks | Runtime business decisions and hidden authorization |
| `database/transaction/` | GORM context propagation, UnitOfWork, after-commit hooks | A domain repository or distributed remote transaction |
| `diff/` | Generic deterministic comparison of JSON-compatible values | Publishing domain events or PII policy |
| `errors/` | Structured technical error contracts and wrapping | Service-specific error identity |
| `events/` | Typed generic event envelope, metadata, pipeline rules | NATS/Watermill or outbox persistence entities |
| `headercontext/` | Request/correlation/idempotency metadata context | Business authorization |
| `helpers/` | Established narrow compatibility helpers; migrate unrelated additions to named capabilities | New generic dumping ground |
| `http/` | Request routing/binding and response envelope capability | Resource-specific API policies |
| `id/` | Typed generation/validation of generic IDs | Service ownership decisions |
| `idempotency/` | Technical request-key and fingerprint mechanics | Aggregate-specific lifecycle behavior |
| `logging/` | Structured logging configuration and helpers | Raw secret/payload disclosure |
| `module/` | Generic module lifecycle contracts | Service-specific bootstrap |
| `normalize/` | Reusable technical email, phone, name etc. normalizers when semantics align | Unverified business policy |
| `pipeline/` | Typed REST and service execution runner | Domain-specific strategies or operations |
| `request/` | Common typed request section contract | An individual route's payload |
| `search/` | Allowlisted query, keyset, cursor mechanics | Domain field names or SQL supplied by users |
| `telemetry/` | Tracing/service identity and metric mechanics | Tenant authorization or data leakage |
| `tests/` | Shared test DB/mocks/fixtures capability | Production fixture state |

For generic normalization, never move a value into core only because a similarly named field appears in several services. Stable semantics, output equivalence, source ownership and tests must justify reuse.

## Allowed dependency directions

```text
module business code -> libs/backend/core/<public-capability>
module business code -> libs/infrastructure/<public-capability> (when needed)
libs/backend/core/<capability> -> sibling generic facade, stdlib, approved technical libs
libs/backend/core/<capability> -X-> services/*, vendor-specific externalService/*
```

A root facade should be small. Avoid re-exporting every internal helper, impl factory, validator, or generated type from it.

## Runtime boundary notes

- UnitOfWork coordinates **one PostgreSQL database**, not a distributed transaction.
- Search requires trusted resource-owned SQL configs and allowlisted operators.
- Events are domain-neutral envelopes; reliable delivery lives with Infrastructure.
- Code comments with annotations are **not runtime magic**.
- Security scopes/principal context must be enforced by the actual consuming service and launcher middleware.

## Quality gates

Capability changes need unit tests, API compatibility assessment for all imports, generated-overlay compilation where needed and consumer integration tests for changed behavior.
