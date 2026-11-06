# Integration Service — canonical contract

This is the current source of truth for `backend/services/integration`. For **all** backend modules, also follow `GROVEHUB_BACKEND_STRUCTURE.md`.

## Ownership

Integration owns provider discovery, external account connections, OAuth/refresh credentials, encrypted secret records, connectivity tests, provider webhook subscriptions, token validation and stateless provider operations.

Communication (not yet implemented) owns business mailbox setup, message history, threads, attachments, outbound intent/delivery status, sync checkpoints, duplicate suppression and per-provider message ingestion. Integration must not persist those resources.

An external Gmail notification is authenticated by Integration and published as a generic event; Communication will later decide when and how to synchronize. Temporal belongs to a future workflow orchestrator, never a mandatory Integration dependency.

## Layer architecture

```text
providers/providers.go                   public registry + capability interfaces
providers/impl/{registry,action,oauth}.go concrete registries
providers/{gmail,outlook}/              public adapter facade
providers/{gmail,outlook}/impl/         provider HTTP/OAuth/actions implementations
config/config.go                         public configuration Loader interface and types
config/impl/config.go                    environment-backed implementation
constants/action.go                      action, capability, provider, resource identifiers
constants/error/integration.go           error definitions and sentinels
controller/{connection,oauth,...}.go     separate public controller interfaces
controller/impl/{connection,oauth,...}   matching HTTP implementations
controller/impl/common.go                shared trusted principal/error helpers
routes/{connection,oauth,...}.go         independent route groups
routes/http.go                           route-group composition
entity/                                  Integration-owned GORM persistence models only
repository/{resource}.go                 public repository interfaces
repository/impl/                          ORM-first repository implementations
events/mailbox/changed.go               alias of coreevents.Envelope[Payload]
events/mailbox/helpers/changed.go        event envelope helpers
events/mailbox/impl/changed.go           event publication
```

## Provider capabilities

Every provider registers only its implemented operations via the generic action strategy registry. Consumers use:

```http
GET /integrations/providers
GET /integrations/providers/{provider}
POST /integrations/connections/{id}/actions/{action}
Content-Type: application/json

{"input": {"messageId":"..."}}
```

Gmail supports profile, list, get, attachment, history, send and Gmail watch. Outlook supports profile, list, get, attachment, delta and send. Outlook webhook watch and durable event ingestion have not been implemented. The generic controller does not choose providers by hard-coded brand.

OAuth initiation uses the connection's provider strategy and a one-time hashed state. Provider callbacks are at `/integrations/oauth/{provider}/callback`. Each active connection has its own encrypted credentials. Both Gmail and Outlook require provider secrets, redirect configuration and a stable 32-byte encryption key. Provider testing supports independent INBOUND read probes and OUTBOUND self-addressed sends, saving sanitized outcomes only.

## Events

The type `mailbox.Changed` is `coreevents.Envelope[mailbox.Payload]`. A changed notification is published on `integration_mailbox_changed` as event type `integration.mailbox.changed.v1`; payload includes tenant, owner, provider, connection ID and provider history *hint*. No authoritative Communication cursor lives in Integration.

## Database

This development branch has **one** `000001_create_integration_schema.sql` baseline, representing connections, encrypted credentials, OAuth states, webhook subscriptions and connection tests. No email ledger/synchronization tables exist in the canonical schema. For already-migrated databases, do **not** replace applied migrations with a squashed baseline: preserve data and create an explicit forward-only upgrade plan. `db-schema-reset` is destructive.

## Developer commands

```bash
make dev-start
make backend-format FORMAT_SCOPE=integration
make db-schema-init DB_SERVICE=integration
make db-schema-status DB_SERVICE=integration
```

`make dev-start` loads Integration OAuth/Pub/Sub environment variables, starts NATS, runs all schema migrations, starts the backend in Air, and starts Swagger. `scripts/config/services.sh` registers Integration OpenAPI.

## Known blockers

Authenticated owner/tenant principal middleware has not yet been wired; protected Integration endpoints deny requests rather than trusting headers. External provider calls, DB migration and Ginkgo tests have not been validated in a live deployment.
