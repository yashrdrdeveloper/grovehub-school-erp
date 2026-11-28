# Core Library — Engineering

`libs/core` owns **reusable domain-neutral technical capabilities**, not business vocabulary or provider-specific messaging. The consumer-facing surface is a named capability facade. Service implementations must not import private nested helper/validator packages.

Generic policy: [Backend Standard](../../../../../docs/engineering/backend/backend.md), [Libraries](../../../../../docs/engineering/backend/libraries.md), [Annotations](../../../../../docs/engineering/backend/annotation.md).

## Owner-specific references

- [capabilities.md](capabilities.md) — inventory, public import direction and SRP for Core capabilities.
- [annotation_engine.md](annotation_engine.md) — parser/checker/generation/overlay pipeline, supported effect boundaries.
- [../../annotation/README.md](../../annotation/README.md) — source-level annotation catalog entry point.

## New capability checklist

Document public facade, contract and error behavior, parameters, concurrency and cancellation rules, allocation/PII risks, tests, owners and reasons the helper cannot live in the consuming service. No capability may import a concrete business module or external provider implementation.

No `docs/product/` is created: library behavior is technical, and actual end-user workflows belong in the consuming service's `docs/product/`.

## Module and validation

Module: `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core`. This module can run its tests with GOWORK=off and never imports the school backend. In this private monorepo, the consuming school module uses a local relative replacement. No public/tagged distribution is claimed.

```bash
cd libs/backend/core
export GOWORK=off
go run ./annotation/cmd testing "$PWD" "$PWD"
go run ./annotation/cmd overlay "$PWD" "$PWD"
go test -vet=off -overlay=.generated/overlay.json ./...
```

Use `make unit-test package=core` from the repository root for report-only coverage; `make unit-coverage package=core` requires 100%.
