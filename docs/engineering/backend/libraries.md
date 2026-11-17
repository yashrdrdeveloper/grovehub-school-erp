# Shared Library Ownership and Capability Design

## 1. Library rule

Shared libraries are **capability-first**. A library exists only when a reusable contract has a distinct owner and proven stable semantics. Duplication within one service is not a reason to publish a global helper.

```text
backend/libs/
├── core/
│   ├── <capability>/
│   │   ├── <capability>.go       # stable public facade
│   │   ├── contracts/           # internal typed contracts if actually required
│   │   ├── helpers/             # algorithm support
│   │   ├── mappers/             # transformation only when needed
│   │   └── validators/          # generic technical checks
│   ├── annotation/
│   ├── database/transaction/
│   └── docs/engineering/
└── <technical-library>/
    ├── <capability>/
    └── docs/engineering/
```

Create only folders with real responsibilities. Services import the public capability facade, **never** a library's private helpers, contracts, mappers or validators unless explicitly part of its API. Libraries must follow Go `internal` import restrictions.

## 2. Core capabilities

Core owns domain-neutral mechanics: typed ID generation, normalization, generic error representation, request/HTTP pipeline, header context, structured logging/tracing, typed search, event envelopes, service pipeline, transaction UnitOfWork, JSON diff, annotation parser/generator and test infrastructure.

Core must not know vendor names, business entity names, application permissions, transport-specific broker storage or event delivery policy. Domain validation remains in its owning service.

## 3. Technical infrastructure contracts

Infrastructure-owned shared contracts, such as broker-delivery DTOs, event-store records, publish options and transport configuration, belong in a separate technical library rather than `libs/core/events`. The actual vendor SDK/transport adapter lives in Infrastructure's `externalService/<capability>/<provider>`.

## 4. Avoiding false reuse

- No `libs/core/helpers/` bucket for unrelated routines.
- No service-specific event/error/DTO in core.
- No generic repository replacing aggregate-specific SQL and lifecycle rules.
- No external app importing another service's private/internal files.
- No hardcoded production broker URL, passwords or tokens in library defaults.
- No global mutable registry unless concurrency, lifetime and cleanup are rigorously defined.
- No duplicate DTO copied solely to bypass a private package import; publish a deliberate stable alias or boundary contract.

## 5. Documentation for libraries

Each owned library has `docs/engineering/README.md`: public entry points, import rules, SRP, initialization, error/timeout semantics, lifecycle, security, usage and tests. Complex capabilities can add topic-specific engineering guides. **Product execution guides belong to consuming services**, not shared libraries.

Breaking library changes must update consumers, generated code, tests, and owner-local docs together.
