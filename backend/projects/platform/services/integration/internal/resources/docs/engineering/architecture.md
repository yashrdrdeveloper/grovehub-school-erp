# Integration — Runtime Architecture

## 1. Implemented layout

```text
services/integration/
├── main.go                                  # typed Bootstrap + compatibility adapters
├── client/connection/connection.go          # @Client and stable DTO aliases
└── internal/
    ├── app/
    │   ├── bootstrap/impl/
    │   │   ├── bootstrap.go                 # runtime Build/Close
    │   │   ├── client.go                    # client proxy
    │   │   ├── repository.go                # GORM dependencies
    │   │   ├── services.go                  # operation services
    │   │   ├── routes.go                    # HTTP controllers
    │   │   ├── events.go                    # mailbox producer
    │   │   ├── external_services.go         # Gmail/Outlook factory/registry
    │   │   └── workers.go                  # watch renewal lifecycle
    │   ├── config/{config.go,impl/config.go}
    │   ├── controller/
    │   │   ├── connection/{public,protected}/
    │   │   │   └── impl/                  # tier-specific adapters
    │   │   ├── provider/public/impl/     # public catalog
    │   │   ├── webhook/restricted/impl/  # verified Gmail push
    │   │   └── helpers/errors.go
    │   ├── dto/{connection,provider}/
    │   ├── entity/                         # connection, webhook and test records
    │   ├── repository/{connection,watch}/
    │   ├── validators/connection/
    │   ├── mappers/connection/
    │   ├── services/
    │   │   ├── connection/operations/{create,get,list,disconnect,oauth,action,test,watch}/
    │   │   └── provider/operations/{list,get}/
    │   ├── providers/                      # capability-neutral registry/DTO/operations
    │   │   └── mailprovider/              # mail-specific factory, OAuth, actions
    │   ├── externalService/mail/{gmail,outlook}/
    │   └── events/producers/connection/operations/mailbox/
    ├── resources/{resources.go,migrations/,api/openapi/v1/}
    └── tests/unit-test/
```

Only real folders are included; a public `config/` is not permitted. Public consumers use `client/connection` or module entrypoint, not another module's `internal/app`.

## 2. SRP by owned resource

**Connection**: owner identity, provider, lifecycle status, encrypted credential reference, OAuth state; repository maps SQL; service operations enforce authorized creation, access, disconnect, action/test/watch transitions.

**Provider discovery**: read-only catalog from the generic provider registry, available actions and OAuth readiness; no credential persistence inside discovery.

**Controller**: provider discovery, protected connection operations, public OAuth callbacks, and restricted Gmail push each own a tier-specific implementation. There is no shared `controller/connection/impl` package. `controller/helpers/errors.go` normalizes errors for HTTP, not business policy.

**Operations**: `services/connection/impl/connection.go` composes eight typed operations, while `services/provider/impl/provider.go` composes list/get. Complex provider calls are in operation strategies. Service code must not switch across Gmail/Outlook implementations.

**DTOs/mappers/validators**: typed input, output, canonicalization, null rules and representation changes; no GORM or provider SDK calls inside validators/mappers.

## 3. API and authorization

Routes are generated from annotated resource controller contracts and now use `/integrations/public/...`, `/integrations/protected/...`, and `/integrations/restricted/...` tiers. This is a **breaking URL change**: update clients, callback allowlists, webhook subscriptions, API references and reverse-proxy rules before deployment. OAuth redirect URI validation expects `/integrations/public/connections/oauth/{provider}/callback`. The generated typed client currently delegates through an in-process service `@Proxy`; this is a trusted internal call and **does not automatically execute the external HTTP middleware**. Connection HTTP routes expect a verified principal/owner; provider discovery can be public where policy permits. The launcher now starts Integration with its developer-authentication bypass disabled. Protected routes receive a verified actor GID and tenant from the gateway; restricted Gmail webhook routes require a validated Pub/Sub workload token.

## 4. Database ownership and migrations

The module has its own embedded SQL migrations and `integration_schema_state` Goose history. It stores connection, OAuth state/credentials, webhook watch subscriptions and connectivity test outcomes, not full email messages or mailbox sync cursors. Tokens are encrypted; they must never be returned by public DTOs or written to logs. Applied migrations are immutable.

## 5. Bootstrap and shutdown

Bootstrap loads validated private config, constructs provider adapters/registries, repositories, validators, mappers, Connection and Provider services, HTTP handlers, client proxy and ChangedProducer. A cancelable six-hour Gmail watch renewal worker starts only if configured; shutdown cancels work, joins the goroutine and closes dependencies.

## 6. Verification

Tests under `internal/resources/tests/unit-test` validate provider factory/registry, controller auth, action/test/watch behavior, OAuth and remote failures, credential encryption, and change-producer error propagation. Updating an external adapter must preserve the public Integration OpenAPI and test each provider independently.

## External-provider Strategy

Gmail and Outlook remain independently owned adapters under `externalService/mail/{gmail,outlook}`. Each registers a concrete `@Provider` and `@Strategy` implementation of the existing `mailprovider.Strategy` contract. The provider registry, OAuth credential handling and action DTOs remain Integration-owned; Identity must only call its typed public Integration client.


## Tenant-scoped system mail sender

An additive migration creates `integration.system_mail_senders` with one persisted connection selector per tenant. `POST /integrations/protected/connections/system-mail/sender` accepts `{"connectionId":"<uuid>"}` only from a gateway-verified tenant administrator. The operation confirms the selected connection is tenant-owned, ACTIVE and supports `mail.send` before saving the selection. `ExecuteSystemMail` revalidates the active connection and routes the delivery through the registered provider action. The setting can change without redeploying Identity or Integration.
