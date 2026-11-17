# Providers, Capabilities and External Service Boundaries

## 1. Separation of responsibilities

A **provider** is a registered implementation of a capability. A **capability** expresses what an integration can do; an **action** is an executable command; an **adapter** translates a vendor's API into an internal contract; a **registry** resolves implementations.

Do not put capability-specific constants or OAuth/token models into the generic provider contract. The generic provider engine must be reusable for storage, payment, communication or any other capability without importing its specific DTOs or vendor SDKs.

## 2. Generic provider blueprint

```text
internal/app/providers/
├── providers.go                        # ONLY top-level Go file: generic Provider interface
├── dto/
│   └── provider.go                     # generic Descriptor, Capability, Action, Session
├── contracts/
│   ├── registry.go                     # Find, List
│   └── executor.go                     # Execute, action handler registry
├── helpers/
│   ├── errors.go                       # generic error classification
│   └── input.go                        # strict input decoding
└── operations/
    ├── register/                       # uniqueness + registration
    ├── resolve/                        # stable identity lookup
    ├── catalog/                        # enumerate providers/capabilities
    └── execute/                        # typed action dispatch
```

Operation-specific `strategies/strategy.go` and `strategies/impl/<operation>.go` exist only when a replaceable algorithm or substantial implementation exists. The generic `providers.go` MUST NOT re-export every action, token, error or capability category.

## 3. Capability specialization

```text
providers/<capability>provider/
├── <capability>provider.go             # typed strategy, factory, resolver
├── constants/                          # <capability>.* actions and capabilities
├── dto/                                # capability-specific request/result models
├── contracts/                          # capability-only OAuth, read/write or watch
├── helpers/                            # capability-owned encoding/error logic
└── operations/
    ├── build/
    ├── authorize/                      # only if meaningful
    └── <capability-operation>/

externalService/<capability>/<vendor>/
├── <vendor>.go                          # @ExternalService contract when meaningful
└── adapters/ or impl/                   # real HTTP/SDK behavior
```

Technical OAuth mechanics may be reused if protocol behavior is stable, but scopes, credential formats, refresh policy, callbacks and token encryption remain specialized. Generic provider code must never import specialized capability DTOs.

## 4. Adding a provider

1. Define narrow capability interfaces, DTOs, errors and feature negotiation.
2. Implement a vendor adapter behind an explicit external-service boundary.
3. Mark one concrete adapter type `@Provider id=NAME category=kind`.
4. Register descriptor, supported actions and optional capabilities in a factory.
5. Validate env secrets and construct clients in bootstrap.
6. Inject generic registry and specialized strategy into operations.
7. Expose actions only through trusted authorization boundaries.
8. Add unit, contract, sandbox/provider and product walkthrough tests.

The `@Provider` static checker detects malformed/duplicate IDs; it **does not connect a vendor**, obtain OAuth credentials, or register a runtime client automatically.

## 5. Pattern decision matrix

| Requirement | Pattern |
|---|---|
| Stable provider ID to implementation | Registry |
| Construct consistent family with validated config | Factory |
| Translate SDK types to DTOs | Adapter |
| Alternate behavior by provider | Strategy |
| Validate/execute named action | Command |
| Public small typed contract | Facade |

Do not force all providers to implement OAuth, subscriptions and actions. Optional capabilities are separate interfaces (Interface Segregation). All `<capability>.*` action strings live under their owning capability constants, not in the generic registry.

## 6. Security and resilience

- OAuth callbacks verify one-time state, origin/audience, replay and scoped identity.
- Tokens are encrypted at rest and excluded from logs, URLs, public responses and product examples.
- Authorize owner/tenant before provider action execution, not merely by valid token.
- Use bounded timeouts, retry budgets, cancellation, input/body limits and explicit errors.
- Review SSRF, webhook signature validation, provider scopes, MIME handling and 429 backoff.
- A verified webhook may emit a domain event; it must honor the upstream retry/ACK contract.

## 7. Provider test matrix

Test duplicate registration, unknown provider, deterministic catalog, disabled capability, invalid input, wrong principal/tenant, token expiry/refresh, remote timeout/429, provider response changes, webhook replay, secret redaction and cancellation. Registry tests must not depend on live vendor APIs.
