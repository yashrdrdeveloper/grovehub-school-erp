# Integration Service — Foundation-style Structure

> Applies to `backend/services/integration`. Public HTTP contracts and OAuth/webhook URLs are unchanged.

Integration follows Foundation's **resource-first internal/app** boundaries. The public package exposes only module lifecycle and typed clients; configuration contracts required by the launcher remain public.

```text
services/integration/
├── integration.go                         # New + Migrate delegates to bootstrap
├── config/config.go                      # public configuration contracts
├── client/
│   ├── client.go                         # backwards-compatible Client alias
│   └── connection/connection.go          # @Client + public DTO aliases
└── internal/
    ├── app/
    │   ├── bootstrap/
    │   │   ├── bootstrap.go              # Bootstrap + Runtime contracts
    │   │   └── impl/
    │   │       ├── bootstrap.go          # builder.Build + module Close
    │   │       ├── providers.go          # Gmail/Outlook adapter composition
    │   │       ├── repository.go         # Connection + Watch repositories
    │   │       ├── services.go           # Connection + Provider services
    │   │       ├── routes.go             # Controller assembly
    │   │       ├── client.go             # generated client adapter
    │   │       └── events.go             # Gmail watch renewal lifecycle
    │   ├── client/connection/impl/client.go
    │   ├── config/impl/config.go         # typed environment configuration loader
    │   ├── constants/{actions,connections,error,providers}/
    │   ├── controller/
    │   │   ├── connections.go            # @Routes contract (unchanged URL)
    │   │   ├── provider.go               # @Routes contract (unchanged URL)
    │   │   ├── controller.go             # interface composition
    │   │   └── connection/impl/          # request/auth/webhook handlers
    │   ├── dto/
    │   │   ├── connection/{request,response,watch}.go
    │   │   └── provider/response.go
    │   ├── entity/{connection,oauth,test,webhook}.go
    │   ├── events/producers/connection/operations/mailbox/
    │   │   ├── changed.go                # versioned @Event
    │   │   ├── helpers/changed.go
    │   │   └── impl/changed.go
    │   ├── providers/
    │   │   ├── impl/                     # provider registries
    │   │   └── mailprovider/providers/
    │   │       ├── gmail/                # Gmail OAuth/Watch/Actions
    │   │       └── outlook/              # Outlook OAuth/Actions
    │   ├── repository/
    │   │   ├── connection/{connection.go,impl/connection.go}
    │   │   └── watch/{watch.go,impl/watch.go}
    │   ├── validators/connection/{connection.go,helpers/gid.go,impl/connection.go}
    │   ├── services/
    │   │   ├── connection/{connection.go,impl/,helpers/,operations/}
    │   │   └── provider/{provider.go,impl/,operations/}
    │   └── routes/http.go
    ├── resources/
    │   ├── resources.go                  # embeds SQL and OpenAPI YAML
    │   ├── api/openapi/v1/openapi.yaml
    │   └── migrations/{migrate.go,000001_create_integration_schema.sql}
    └── tests/unit-test/{config,controller,providers,services}/
```

## Public boundary

The stable `integration.New(logger, db, messaging, optionalDevAuth...)` entry point and
`integration.Module` contract are unchanged. Consumers depend on
`services/integration/client` (alias to `client/connection`), not on
`internal/app`. The generated `@Client` proxy forwards to Integration services;
it does not make a network call.

Configuration contracts for Gmail, Outlook and opt-in loopback-only DevAuth
remain in public `config/config.go`; parsing/validation is private.
The dev auth bypass must never be enabled on a production listener or trusted
based on arbitrary forwarded HTTP headers.

## Existing APIs (unchanged)

Provider discovery: `GET /integrations/providers`,
`GET /integrations/providers/{provider}`.

Connections: `POST /integrations/connections`,
`GET /integrations/connections`, `GET /integrations/connections/{id}`,
`DELETE /integrations/connections/{id}`.

OAuth: `POST /integrations/connections/{id}/oauth/start`,
`GET /integrations/connections/oauth/{provider}/callback`.

Testing and actions: `POST /integrations/connections/{id}/test`,
`GET /integrations/connections/{id}/tests`,
`POST /integrations/connections/{id}/actions/{action}`.

Webhooks: `POST /integrations/connections/{id}/watch`,
`POST /integrations/webhooks/gmail`.

The route annotations remain in the controller contract package to retain the
existing `/integrations` URL namespace. Do not introduce `/protected` or
`/public` prefixes merely to reproduce Foundation's path naming.

## Event and persistence policy

Integration emits `integration.mailbox.changed.v1` to the existing
`integration_mailbox_changed` subject, only as a provider notification hint.
It does not store email messages or commit Communication-owned sync cursors.

Integration's own repositories use GORM. The module reuses the shared database,
has its own `integration_schema_state` Goose history, and embeds its migration
SQL for binary `migrate integration` and `migrate all`. Dev scripts and Swagger
read the moved embedded resource folders.

## Verification

```bash
make backend-build
cd backend
go test -vet=off -overlay=.generated/overlay.json ./services/integration/internal/tests/unit-test/...
```

Run `make dev-start` for an end-to-end development start after setting your
Gmail/Outlook OAuth values and Infrastructure NATS variables.

## Connection operation pattern

ConnectionService is composed exactly like Foundation's PersonService:
the public interface lives at `internal/app/services/connection/connection.go`,
and `services/connection/impl/connections.go` embeds eight real operations:

```text
services/connection/
├── connection.go                           # ConnectionService interface
├── impl/connections.go                    # @Service composition only
├── helpers/{errors,credential,...}.go     # shared pure/error and crypto helpers
└── operations/
    ├── create/{create.go,strategies/{strategy.go,impl/create.go}}
    ├── get/{get.go,strategies/{strategy.go,impl/get.go}}
    ├── list/{list.go,strategies/{strategy.go,impl/list.go}}
    ├── disconnect/{disconnect.go,strategies/{strategy.go,impl/disconnect.go}}
    ├── oauth/{oauth.go,strategies/{strategy.go,impl/oauth.go}}
    ├── action/{action.go,strategies/{strategy.go,impl/action.go}}
    ├── test/{test.go,strategies/{strategy.go,impl/test.go}}
    └── watch/{watch.go,strategies/{strategy.go,impl/watch.go}}
```

Each strategy owns its operation logic; the composed service does not
duplicate methods or switch on Gmail/Outlook. Shared OAuth token refresh is
exposed through the typed `oauth.TokenProvider` interface to the Action,
Test and Watch operations. Credential encryption uses the same AES-GCM
behavior, key validation and nonce generation as before the refactor.
The public controller error behavior remains compatible through canonical
errors in `services/connection/helpers/errors.go`.

## Provider discovery operations

Provider discovery mirrors the same composition pattern without introducing
vendor-specific behavior into the public service. The
`services/provider/impl/ProviderService` embeds
`list.ListOperation` and `get.GetOperation`, backed by separate typed
`strategies/impl` implementations. List resolves registered provider
descriptors. Get resolves a provider descriptor, configured action catalog,
and OAuth readiness. The existing `Providers` and `Provider` signatures and
public Integration client methods remain unchanged.

## Typed Connection validation

Connection operations `Create`, `Get`, `List`, and `Disconnect` use
`internal/app/validators/connection.ConnectionValidator`, injected through
Integration bootstrap, as Foundation injects PersonValidator. The validator
reuses the platform's existing UUIDv7 normalization/validation semantics.
It does not perform SQL lookups; those remain in the repositories and services.

Repositories are annotated `@Repository entity=Connection` and
`@Repository entity=WebhookSubscription` and remain GORM-only.
