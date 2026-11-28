# Infrastructure service

Messaging, event store, replay source and outbox dispatch.

`main.go` exposes the typed bootstrap and transitional compatibility functions. Public client contracts live under `client/`; implementation is under `internal/app`.

## Documentation

- [Engineering](engineering/README.md)
- [Product](product/README.md)

Test suites, fixtures, helpers and mocks belong under `internal/resources/tests/`. Unit and integration tests remain separated by suite type.
