# Integration Service — canonical architecture

> This is the **current** Integration implementation contract. Follow
> `GROVEHUB_BACKEND_STRUCTURE.md`, `EVENT_GUIDELINES.md`,
> `NORMALIZATION_AND_DRY_GUIDELINES.md`, and
> `docs/engineering/testing/BACKEND_UNIT_TEST_STRUCTURE.md` as well.

## Ownership

**Integration** owns provider catalogs, provider selection, OAuth, encrypted tokens, refresh, provider connection lifecycle, connectivity tests, provider webhook subscriptions, authenticated provider notifications, and stateless provider actions.

**Communication** (not implemented) owns mailbox/user messaging records, messages, threads, attachments, message deduplication, delivery state and synchronization checkpoints. Integration does not persist messages or consumer cursors. **Temporal is not an Integration dependency.**

## Canonical package layout

```text
integration/
  integration.go                              composition + Module.Client()
  client/
    client.go                                public in-process Client interface
    impl/client.go                           service-delegating implementation
  config/config.go                           public Loader interface and config contracts
  config/impl/config.go                      environment implementation
  constants/
    actions/gmail.go                         type-safe mail actions and capabilities
    providers/providers.go                   provider identifiers
    connections/status.go                    connection status and resource identifiers
    error/integration.go                     errors
  controller/
    connections.go                           all connection HTTP method contracts
    provider.go                              provider discovery HTTP contract
    controller.go                            composed controller
    impl/connections.go                      all connection HTTP handling
    impl/provider.go                         provider discovery HTTP handling
  routes/
    http.go                                  compose resource route groups
    connections.go                           EVERY /integrations/connections/** endpoint
    provider.go                              EVERY /integrations/providers/** endpoint
    webhooks.go                              /integrations/webhooks/** endpoint
    helpers/http.go                          shared route registration and tracing
  services/
    connections.go                           connection and provider-action service interface
    providers.go                             separate ProviderService interface
    impl/connections.go                      connection lifecycle implementation
    impl/providers.go                        provider discovery service implementation
    impl/oauth.go                            connection OAuth implementation
    impl/action.go                           provider action execution implementation
    impl/test.go                             connectivity tests
    impl/watch.go                            provider watch/notifications
    helpers/connections/                     pure UUIDv7 checks and response conversion
  providers/
    providers.go                             generic provider/strategy registry contracts
    oauth.go                                 generic optional OAuth strategy contract
    impl/                                    generic registry implementations
    helpers/                                 shared strict action input parser
    mailprovider/
      mailprovider.go                        mail strategy, factory, runtime/resolver contracts
      impl/factory.go                        provider factory and capability-based resolver
      providers/
        gmail/
          gmail.go                           Gmail public facade
          strategy.go                        Gmail mailprovider.Strategy
          impl/                              Gmail OAuth/Gmail API/webhook implementation
        outlook/
          outlook.go                         Outlook public facade
          strategy.go                        Outlook mailprovider.Strategy
          impl/                              Outlook OAuth/Graph API implementation
  events/
    connections/mailbox/
      changed.go                             alias to coreevents.Envelope[Payload]
      helpers/changed.go                     pure envelope/metadata builders
      impl/changed.go                        NATS publish orchestration
  entity/                                    connection/OAuth/test/webhook GORM entities only
  repository/                                public data access contracts
    impl/                                    GORM-first concrete repositories
  migrations/000001_create_integration_schema.sql
  tests/unit-test/                           separate Ginkgo suites for config/controller/providers/services
```

**Public vs private:** Consumers of Integration import `services/integration/client`, or stable provider DTO contracts where needed. They do not import Integration's GORM entities, `services/impl`, `providers/impl`, or vendor implementation packages. Composition roots may import concrete implementation packages.

**Factory:** A `mailprovider.Strategy` declares its descriptor and optionally registers actions/OAuth. The factory builds the provider registry, OAuth registry and action registry. `Runtime.Resolver.Resolve(providerID, capability)` selects a declared provider without a hard-coded provider switch in the consumer; it does not imply that OAuth credentials are configured. The provider details endpoint shows runtime `availableActions`.

## HTTP routes by resource

```text
GET    /integrations/providers
GET    /integrations/providers/{provider}

POST   /integrations/connections
GET    /integrations/connections
GET    /integrations/connections/{id}
DELETE /integrations/connections/{id}
POST   /integrations/connections/{id}/oauth/start
GET    /integrations/connections/oauth/{provider}/callback
POST   /integrations/connections/{id}/test
GET    /integrations/connections/{id}/tests
POST   /integrations/connections/{id}/watch
POST   /integrations/connections/{id}/actions/{action}

POST   /integrations/webhooks/gmail
```

All `/integrations/connections/**` routes are registered from `routes/connections.go`, even when their handlers represent OAuth, action execution or tests. Public controller declarations likewise share `controller/connections.go`; the actual methods live in `controller/impl/connections.go`. Provider discovery remains separate: `controller/provider.go`, `controller/impl/provider.go`, `routes/provider.go` and `services/impl/providers.go`.

Gmail and Outlook OAuth redirect URIs must use the corresponding `/integrations/connections/oauth/{provider}/callback` endpoint; update registered provider-console callback URLs after the route change.

## In-process client for other backend services

Follow Foundation's established `client/` + `client/impl/` convention:

```go
integrationClient := integrationModule.Client()

providers := integrationClient.Providers(ctx)
result, err := integrationClient.ExecuteAction(
    ctx, tenantID, connectionID, integrationproviders.ActionMailGet,
    json.RawMessage(`{"messageId":"provider-id"}`),
)
```

The client calls the application services in the same modular monolith; it does **not** make an unnecessary HTTP loopback call. The caller is responsible for established cross-service authorization and tenant/owner scopes. Do not expose this client directly to an untrusted public caller.

## Events

`mailbox.Changed` aliases `coreevents.Envelope[mailbox.Payload]`. Gmail Pub/Sub authentication and mailbox connection resolution happen in Integration. The publisher emits `integration.mailbox.changed.v1` to the stable `integration_mailbox_changed` subject using the generic core envelope. The event contains provider/connection/owner/tenant and a history **hint** only, not message contents or a committed checkpoint. A future Communication consumer will own idempotent ingestion and cursor advancement.

## Schema, testing and startup

The single `000001` migration is a **fresh-development baseline only**; do not apply the squashed migration list on existing databases with Goose versions 2–5 or important archived data. Use forward-only migrations for deployed databases.

```bash
make backend-format FORMAT_SCOPE=integration
cd backend && go test ./services/integration/...
make dev-start
```

Tests belong under `tests/unit-test/<owner>/` with `suite_test.go`, Ginkgo v2 and Gomega. Never put `*_test.go` in runtime implementation folders.

`dev-start` initializes Foundation, Identity and Integration schemas (Infrastructure has no migration directory), then starts NATS, backend Air and Swagger. Gmail/Outlook API, OAuth, webhook and database tests must still be verified against real configuration.

## Temporary local-only Integration identity

While IAM's verified principal middleware is not integrated, the local GroveHub
launcher passes a fixed development principal to Integration. There are **no
authentication bypass environment variables**. The development identity is
configured in `backend/launcher/launcher.go`:

- Tenant ID: `018f5d04-315c-75a9-b890-df0d02e787c1`
- Actor GID: `01a1b210-402c-7622-a33c-f82ceba26c7d`

Use this actor as `ownerGid` when creating connections. All protected Integration
routes can use the fixed principal **only for loopback clients**, with the launcher
binding `127.0.0.1:8080`. HTTP headers never supply or override the identity.
OAuth state and Pub/Sub webhook signature checks remain active.

**Development-only security exception:** this code must not be deployed or
exposed through Cloudflare Tunnel, reverse proxies, or other externally
accessible interfaces. Replace the launcher constant with verified IAM principal
middleware before deployment.

Gmail `POST /integrations/connections/{id}/watch` requires a connected,
active Gmail account **and** `GROVEHUB_GMAIL_PUBSUB_TOPIC` and corresponding
Google Pub/Sub setup. Missing configuration correctly returns HTTP 503 and
is unrelated to the temporary identity.

## Known limitations

Verified IAM middleware has not been connected to the Integration HTTP routes; protected endpoints deny rather than trusting arbitrary tenant headers. Outlook webhook watch is not implemented. Integration notifications rely on messaging delivery but do not yet implement a transactional outbox. No Communication implementation is included.

### Gmail mail.send input validation

`mail.send` accepts `{"input":{"raw":"..."}}` where `raw` is a base64url (unpadded) encoded RFC 5322 message with From/To headers and a header/body separator. Validation accepts CRLF or LF line endings; generating RFC 5322 via Python `email.policy.SMTP` is recommended for Gmail interoperability. An explanatory placeholder is not sendable. Invalid encodings or malformed messages return HTTP 400 `INTEGRATION_INVALID`, not a generic HTTP 500. The mail provider still enforces scopes and may reject syntactically valid messages. Never store message contents in Integration.
