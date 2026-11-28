# CMD architecture

`main.go` delegates to `internal/src/bootstrap/impl` and exits with the returned code. The execution contract aliases `command.Command`; the private executor forwards arguments and streams. CMD does not own database connections, migrations or service strategies.

Launcher owns the registry, database adapters and serve/migrate strategies; service bootstrap contributes commands and migrations. Unknown commands return 2 without opening a database; success returns 0, operational failure 1.

Run `make backend-build` and `make unit-coverage package=cmd` from the repository root. Operator coverage is checked separately by `scripts/testing/platform-operator-coverage.py`.
