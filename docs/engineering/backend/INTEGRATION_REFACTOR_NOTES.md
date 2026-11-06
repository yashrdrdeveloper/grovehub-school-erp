# Integration structural refactor — developer branch

## Current shape

Provider connection and test entities remain in Integration. Provider registry and OAuth/action registry interfaces are in `providers/`, implementations under `providers/impl/`. Gmail and Outlook adapters use `providers/{gmail,outlook}/impl/` with small public provider facades.

Controller interfaces are separated into connection, provider, OAuth, test, watch, webhook and action. Their methods are implemented in matching `controller/impl/` files. `routes/http.go` composes resource route groups.

All runtime environment configuration is loaded by `config.Loader` with the environmental implementation in `config/impl/`. Provider identifiers and action/capability constants live under `constants/`. Mailbox notification events embed `coreevents.Envelope[Payload]` and are built/published by event helpers and implementation.

Webhook subscriptions are queried with GORM entity predicates and association `Preload`. No Integration-owned message or mailbox sync state is introduced.

## Database development baseline

This branch now contains one `000001_create_integration_schema.sql` migration for a **new database** with connections, OAuth, tests and webhook subscriptions only. The old `000002` through `000005` migration history was consolidated as requested.

**Do not run `make db-schema-reset DB_SERVICE=integration` on a database containing important or previously migrated data.** The consolidation is incompatible with an existing Goose schema state that has already applied versions 2–5. Existing message/cursor legacy data must be archived and migrated separately before transitioning an existing database. There is no automatic safe upgrade from an applied 5-migration schema to this one-file baseline. For existing databases preserve the original migration history and use forward-only migrations instead.

`make dev-start` already initializes every registered service including Integration. Integration OAuth/webhook values are now exported by `scripts/lib/env.sh` before startup. `make db-schema-init DB_SERVICE=integration` initializes Integration only.

## Runtime limitations

No verified IAM middleware is wired into the launcher; protected routes correctly deny access. Gmail OAuth and Outlook OAuth require credentials to test. Outlook webhook subscriptions are not implemented, and Communication message persistence remains outside this service. No live tests or migration execution have been performed.
