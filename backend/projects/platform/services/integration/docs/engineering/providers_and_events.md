# Integration — Generic Provider Engine, Mail Adapters and Events

## 1. Generic provider engine

`internal/app/providers/providers.go` contains only the generic `Provider` descriptor interface.

- `providers/dto/provider.go`: generic `Action`, `Capability`, `Descriptor`, `Session`.
- `providers/contracts/`: generic registry, action/executor contracts.
- `providers/helpers/`: generic error mapping and strict input decoding.
- `providers/operations/{register,resolve,catalog,execute}/`: typed registry and execution algorithms.
- `providers/mailprovider/`: **mail-only** OAuth token DTOs/contracts, `mail.*` action/capability constants, MIME helpers, factory and OAuth registry operations.
- `externalService/mail/{gmail,outlook}/`: actual provider SDK/HTTP/OAuth implementation.

Provider-neutral packages must not import `mailprovider` or vendor implementations. The `@Provider id=... category=mail` annotations identify concrete Gmail/Outlook strategy types to the annotation checker; they do **not** auto-register runtime clients or tokens. `bootstrap/impl/external_services.go` wires them explicitly.

## 2. Connection transitions

Create a provider connection in pending setup; authorize with one-time OAuth state and code exchange; bind the verified mailbox; activate credentials on success. All connection operations must prove the authenticated actor owns the connection and request scope. Disconnect must invalidate/erase credential state and cease provider actions. A refresh uses the owner-bound token provider and preserves safe secret handling.

## 3. Supported provider actions

Gmail: `mail.profile`, `mail.list`, `mail.get`, `mail.history`, `mail.attachment.get`, `mail.send`.

Outlook: `mail.profile`, `mail.list`, `mail.get`, `mail.delta`, `mail.attachment.get`, `mail.send`.

Only enabled and configured actions should appear in the provider catalog. The `mail.send` API uses structured recipients/body/attachments; provider adapters translate into the appropriate API/MIME representation. Integration does not claim a send API response is a delivery receipt.

## 4. Gmail webhook event

A verified Gmail Pub/Sub push is decoded and matched to active, owned subscriptions, then the typed ChangedProducer emits `integration.mailbox.changed.v1` on `integration_mailbox_changed`.

The event is a **notification hint**, not mail content or a committed sync cursor. It is published **directly** through the Infrastructure messaging client; when broker publication fails, the webhook responds with a retryable error to the upstream sender. The owner of the downstream mailbox synchronization state is separate; Integration does not store messages.

```text
signed webhook -> verify JWT and audience -> find active watch
  -> build typed mailbox-change hint -> ChangedProducer.Publish
  -> NATS -> downstream consumer (owned elsewhere)
```

Never acknowledge a failed broker publish as successful. The mailbox producer is under `events/producers/connection/operations/mailbox/` with typed flow, strategy and helper; the watch operation no longer assembles JSON/event metadata itself.

## 5. Watch renewal

A cancelable background worker schedules renewal roughly every six hours with a bounded child timeout. Its lifecycle is owned by `bootstrap/impl/workers.go`, not by the HTTP controller. Expiring subscription state persists in the module's watch repository.

## 6. Security and tests

Test OAuth state reuse, missing/invalid principal, unowned connection, credential expiry/refresh failure, webhook JWT mismatch, wrong provider capabilities, malformed mail send input, 429/timeouts, remote error mapping, direct-publish failure upstream retry and worker cancellation. No real OAuth secrets in documentation or test fixture.
