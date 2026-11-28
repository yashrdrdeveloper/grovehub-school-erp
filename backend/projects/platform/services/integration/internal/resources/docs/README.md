# Integration service

Provider connections, email reconciliation and provider workers.

`main.go` exposes the typed bootstrap and transitional compatibility functions. Public client contracts live under `client/`; implementation is under `internal/app`.

## Documentation

- [Engineering](engineering/README.md)
- [Product](product/README.md)

Test suites, fixtures, helpers and mocks belong under `internal/resources/tests/`. Unit and integration tests remain separated by suite type.
