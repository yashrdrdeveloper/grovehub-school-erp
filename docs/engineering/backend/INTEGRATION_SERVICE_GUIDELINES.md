# Integration Service — Foundation-style Structure

> Applies to `backend/services/integration`. Public HTTP contracts and OAuth/webhook URLs are unchanged.

Integration follows Foundation's **resource-first internal/app** boundaries. The public package exposes module lifecycle, typed clients, and only the small `integration.DevAuth` bootstrap alias. Gmail, Outlook, push, and loader configuration types are private.

```text
services/integration/
├── integration.go                         # New + Migrate delegates to bootstrap
├── client/
│   ├── client.go                         # backwards-compatible Client alias
│   └── connection/connection.go          # @Client + public DTO aliases
└── internal/
    ├── app/
    │   ├── bootstrap/
    │   │   ├── bootstrap.go              # Bootstrap + Runtime contracts
    │   │   └── impl/
    │   │       ├── bootstrap.go          # builder.Build + module Close
    │   │       ├── external_services.go # Gmail/Outlook adapter composition
    │   │       ├── repository.go         # Connection + Watch repositories
    │   │       ├── services.go           # Connection + Provider services
    │   │       ├── routes.go             # Controller assembly
    │   │       ├── client.go             # generated client adapter
    │   │       ├── events.go             # Connection ChangedProducer assembly
    │   │       └── workers.go            # cancelable Gmail watch renewal lifecycle
    │   ├── client/connection/impl/client.go
    │   ├── config/
    │   │   ├── config.go                # private Gmail, Outlook, Push, DevAuth, Loader
    │   │   └── impl/config.go           # typed environment configuration loader
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
    │   ├── mappers/connection/{connection.go,impl/connection.go}
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

Configuration contracts for Gmail, Outlook, GmailPush and the loader now live
in `internal/app/config/config.go`, and the private environment loader remains
in `internal/app/config/impl/config.go`. The launcher uses only
`integration.DevAuth`, a minimal public alias of the private bootstrap
option; no provider credential configuration is exported through `integration/config`.
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
and `services/connection/impl/connection.go` embeds eight real operations:

```text
services/connection/
├── connection.go                           # ConnectionService interface
├── impl/connection.go                     # @Service composition only
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

## Connection DTO mapping

Connection and connection-test entity-to-response mapping belongs exclusively
to `mappers/connection` and its `impl` package, through
`ConnectionMapper`, matching Foundation's typed `@Mapper` pattern.
The mapper is injected into Create, Get, List, Disconnect and Test operations
from `bootstrap/impl/services.go`. Shared service helpers no longer own
entity-to-DTO conversion logic. Unit tests cover mapping and UUIDv7 validation.

## Provider boundary and controller composition

Integration's provider-neutral registries, action/OAuth contracts and mail
factory stay under `internal/app/providers`. **Vendor-specific network,
protocol and OAuth implementations** are isolated under
`internal/app/externalService/mail/gmail` and
`internal/app/externalService/mail/outlook`. The bootstrap's
`impl/external_services.go` registers those adapters into provider-neutral
registries. Connection services and controllers never import vendor HTTP SDKs
other than the narrow integration points that require Gmail watch and push.

The HTTP layer is split into `controller/connection/impl` and
`controller/provider/impl` with a typed composite in
`controller/impl`. Shared structured HTTP error mapping lives under
`controller/helpers`; authentication and ownership checks remain exclusive
to Connection endpoints. Both `@Routes` interfaces remain at the controller
root because the current route generator requires this placement to preserve
existing URLs *without adding* `/public` or `/protected` path segments.

Foundation's generated `@RequestLine` client goes through an HTTP route
because Foundation has in-process controller request DTO semantics. Integration
keeps its existing `@Proxy` service client for trusted internal callers: its
method contracts require tenant and owner IDs, whereas the HTTP API derives
identity from **verified principal context**. Switching to HTTP dispatch without
an explicit internal authorization boundary could bypass that protection.

The OAuth webhook watch-renewal worker now derives its timeouts from the
module's cancelable context. Closing Integration cancels in-flight renewal work
and joins its goroutine rather than waiting for an unrelated timeout.

### Refactor invariants

- No changes to public Integration HTTP paths, response envelopes, or OAuth callbacks.
- No changes to Integration SQL migration versions or table names.
- No new direct dependencies on vendor implementations from unrelated services.
- Infrastructure owns messaging transport; Integration emits only mailbox-change hints.
- Existing authenticated Connection routes retain ownership checks; provider discovery remains public.

## Domain event producer (Foundation pattern)

Integration now owns a typed event producer for mailbox-change hints:

```text
events/producers/connection/
├── connection.go                             # ChangedProducer interface
├── impl/connection.go                        # event producer implementation
└── operations/mailbox/
    ├── changed.go                             # @Event schema, unchanged topic
    ├── flows/changed.go                       # typed provider notification payload
    ├── helpers/changed.go                     # envelope / metadata builder
    ├── pipeline.go                            # Build then direct Publish
    └── strategies/
        ├── strategy.go                        # typed @Strategy
        └── impl/changed.go                    # serialization strategy
```

`bootstrap/impl/events.go` builds the producer and injects it into the
ConnectionService's Watch operation. The Watch strategy no longer serializes or
publishes domain events itself; it only selects authorized matching
connections and sends typed notification data to the producer.

**Webhook semantics remain unchanged**: mailbox-change notifications publish
directly to the configured message broker. A failure returns an error so Google
Pub/Sub can retry the delivery. A transactional PostgreSQL outbox is not
appropriate for acknowledging this external provider notification before
NATS delivery. The original event type, subject and schema are preserved.

The six-hour watch renewal worker now lives in
`bootstrap/impl/workers.go`, separate from event-producer composition. Its
context is canceled when the module closes.

## Resource naming

Resource implementation files use singular names:
`controller/connection/impl/connection.go`,
`services/connection/impl/connection.go`,
`services/provider/impl/provider.go`.
The Connection `Watch` response is grouped with its other response DTOs in
`dto/connection/response.go` rather than an unrelated extra file.

## Declaration spacing standard

Every Integration Go source file separates top-level `type`, `func`, `var`,
and `const` declarations with one blank line. Keep GoDoc and annotation
comments immediately attached to the declaration they describe, but add a
blank line before that comment block when it follows another declaration.
This includes methods in external mail adapters, repositories, validators,
operations, and tests. The rule changes only whitespace, not method
signatures or behavior. Run `gofmt` to normalize indentation when editing.

Example:

```go
type TokenCipher struct{ key []byte }

func NewTokenCipher(key []byte) (TokenCipher, error) {
    // ...
}

func (cipher TokenCipher) Seal(token providers.OAuthToken) ([]byte, error) {
    // ...
}

func (cipher TokenCipher) Open(payload []byte) (providers.OAuthToken, error) {
    // ...
}
```
