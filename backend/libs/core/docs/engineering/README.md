# Core Shared Library — Engineering Documentation

See [Backend Engineering Master Standard](../../../../../docs/engineering/backend/backend.md).

## Scope

Shared libraries provide **generic, independently reusable mechanisms**; they do not own a service's entities, API flows, business errors, permissions, or use cases.

## Rules

- A public helper must have a precise purpose, stable contract, documented failure modes and clear owner.
- Keep libraries small and dependency-light; avoid importing service packages.
- Separate transport pipelines, persistence primitives, search contracts, error primitives, logging/tracing and annotations by responsibility when implemented.
- Inject external dependencies explicitly; avoid global mutable state.
- State concurrency guarantees and context/cancellation behavior for exported interfaces.
- Version breaking contracts and provide migration instructions to consumers.
- Shared mechanisms must not silently enforce a product-specific policy.
- Do not add a library abstraction until at least one real consumer demonstrates stable reusable semantics.

## Module inventory to maintain

For each actual exported package, document: API, examples, ownership, supported configurations, invariants, concurrency safety, failure semantics, extension points, test commands, backwards compatibility and known consumers.

This is an engineering-only module. Product workflows belong to the consuming services' `docs/product/`.

**Status:** This file establishes consistent library documentation; implementation-specific API inventories must be verified from code and filled in.
