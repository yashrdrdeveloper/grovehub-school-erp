# Engineering Documentation Index

**Purpose:** route an engineer or code assistant to the authoritative document for the task without loading unrelated material.

## Canonical locations

```text
docs/
├── AGENT.md                            # repository agent routing
└── engineering/
    ├── README.md                       # this index
    ├── backend/
    │   ├── backend.md                  # GENERIC backend master
    │   ├── annotation.md
    │   ├── searchDetails.md
    │   ├── events.md
    │   ├── persistence.md
    │   ├── api.md
    │   ├── providers.md
    │   ├── patterns.md
    │   ├── testing.md
    │   └── libraries.md
    ├── api/                            # broad API governance
    ├── frontend/                       # frontend code standards
    ├── database/                       # DB operations
    ├── testing/                        # shared test strategy
    ├── documentation/                  # documentation standards
    ├── providers/                      # provider SDK governance
    ├── security/                       # security policy
    ├── observability/                  # trace/log/metrics
    ├── scripts/                        # Makefile + local tooling
    ├── ci/                             # continuous integration
    ├── deployment/                     # Docker, Kubernetes, releases
    ├── release/
    └── git/
```

Modules own **additional implemented architecture** in `backend/projects/<project>/services/<module>/internal/resources/docs/engineering/`; they own **business API execution** in `backend/projects/<project>/services/<module>/internal/resources/docs/product/`. Shared libraries own `backend/libs/<library>/internal/resources/docs/engineering/`. Generic documentation must not enumerate particular module endpoints, events, tables or provider names.

## Task routing

| Task | Read first | Then |
|---|---|---|
| Go module creation/refactor | [backend/backend.md](backend/backend.md) | Owning module `internal/resources/docs/engineering/` |
| Annotations, generated Go, overlay | [backend/annotation.md](backend/annotation.md) | Library owner annotation engine docs |
| Search/index/cursor | [backend/searchDetails.md](backend/searchDetails.md) | Module-specific projection rules |
| Domain event and broker/outbox | [backend/events.md](backend/events.md) | Publisher/consumer module architecture |
| ACID, migrations, idempotency | [backend/persistence.md](backend/persistence.md) | [database guideline](database/DATABASE_GUIDELINES.md) |
| HTTP, DTO, validator, mapper | [backend/api.md](backend/api.md) | [API contract guideline](api/API_CONTRACT_GUIDELINES.md) |
| Provider strategies and adapters | [backend/providers.md](backend/providers.md) | Module's provider engineering docs |
| Operation/pattern review | [backend/patterns.md](backend/patterns.md) | Module services implementation |
| Go tests and CI | [backend/testing.md](backend/testing.md) | [testing guideline](testing/TESTING_GUIDELINES.md) |
| Core/Infrastructure library | [backend/libraries.md](backend/libraries.md) | Library's `internal/resources/docs/engineering/` |
| Frontend | [frontend guideline](frontend/FRONTEND_CODE_GUIDELINES.md) | Component/package owner docs |
| Security | [security guideline](security/SECURITY_GUIDELINES.md) | Implemented IAM/middleware |
| Tracing and monitoring | [observability guideline](observability/OBSERVABILITY_GUIDELINES.md) | [trace policy](observability/TRACE_CONTEXT_GUIDELINES.md) |
| Local developer scripts | [script guidelines](scripts/SCRIPT_AND_MAKEFILE_GUIDELINES.md) | [local runbook](scripts/GROVEHUB_LOCAL_DEVELOPMENT.md) |
| Documentation changes | [documentation guideline](documentation/DOCUMENTATION_GUIDELINES.md) | Owning module/library docs |
| CI, Git, release, deployment | Respective category below | Runtime/operational runbook |

## Module ownership links

- [Foundation engineering](../../backend/projects/platform/services/foundation/internal/resources/docs/engineering/README.md) and [product](../../backend/projects/platform/services/foundation/internal/resources/docs/product/README.md)
- [Integration engineering](../../backend/projects/platform/services/integration/internal/resources/docs/engineering/README.md) and [product](../../backend/projects/platform/services/integration/internal/resources/docs/product/README.md)
- [Infrastructure engineering](../../backend/projects/platform/services/infrastructure/internal/resources/docs/engineering/README.md) and [product](../../backend/projects/platform/services/infrastructure/internal/resources/docs/product/README.md)
- [Core library engineering](../../backend/libs/core/internal/resources/docs/engineering/README.md)
- [Platform library engineering](../../backend/libs/platform/internal/resources/docs/engineering/README.md)
- [Backend documentation index](../../backend/internal/resources/docs/README.md)
- [Platform CMD architecture](../../backend/projects/platform/cmd/internal/resources/docs/engineering/architecture.md)
- [Audit engineering](../../backend/projects/platform/services/audit/internal/resources/docs/engineering/README.md)

Identity remains on its current documentation/code layout and is not being changed in this documentation refactor.

## Cross-cutting standards

[API response envelopes](api/API_RESPONSE_ENVELOPE_GUIDELINES.md) · [Backend unit-test structure](testing/BACKEND_UNIT_TEST_STRUCTURE.md) · [Regression testing](testing/REGRESSION_TESTING.md) · [CI/CD](ci/CI_CD_GUIDELINES.md) · [Git/PR](git/GIT_AND_PR_GUIDELINES.md) · [Docker](deployment/DOCKER_GUIDELINES.md) · [Kubernetes](deployment/KUBERNETES_GUIDELINES.md) · [Deployment](deployment/DEPLOYMENT_GUIDELINES.md) · [Release](release/RELEASE_GUIDELINES.md) · [Observability](observability/GROVEHUB_OBSERVABILITY.md).

**Rule:** a generic rule lives in the central owner; a concrete module contract lives with that module; an API walkthrough lives in `internal/resources/docs/product/`. Delete obsolete alternatives instead of maintaining contradictory duplicates.
