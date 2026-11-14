# Foundation module: public API, private implementation, embedded migrations

The public Go packages are `services/foundation` and
`services/foundation/client/person`. The client package exposes
Go type aliases for DTO definitions stored privately under `internal/app/dto/person/`.
The OpenAPI specification is stored under `services/foundation/internal/resources/api/openapi/`.
Documentation is a repository asset rather than an importable Go contract.

```text
backend/services/foundation/
├── foundation.go
├── client/person/
│   ├── person.go
│   └── constants/error.go
├── internal/
│   ├── app/
│   │   ├── dto/person/{operations,params,queryparams,request,response}.go
│   │   ├── bootstrap/{bootstrap.go,impl/}
│   │   ├── client/person/impl/
│   │   ├── constants/
│   │   ├── controller/
│   │   ├── entity/
│   │   ├── events/
│   │   ├── mappers/
│   │   ├── repository/
│   │   ├── routes/
│   │   ├── services/
│   │   └── validators/
│   └── resources/
│       ├── resources.go
│       ├── api/openapi/v1/openapi.yaml
│       └── migrations/
│           ├── migrate.go
│           ├── 000001_create_foundation_schema.sql
│           └── 000002_create_event_outbox.sql
└── docs/
```

## Migration contract

`resources.Migrations` uses Go's `//go:embed migrations/*.sql`.
`resources.OpenAPI` embeds the YAML specification independently;
Swagger development configuration reads the relocated YAML file directly.
Run `groovehub-platform migrate foundation` in production after configuring
`DATABASE_URL`, before the HTTP process starts. This invokes Goose on the
embedded files, with the `foundation_schema_state` version table used by
`scripts/database/migrate.sh` in development. Never edit an applied migration.
The embedded command only applies Foundation migrations; Identity and
Integration continue using their existing migration scripts.

## Annotation generation

Original source is nested under `internal/app`, while generated Go files are
written under `.generated/internal/app` and applied with the existing Go
overlay. Private Person DTO PATCH generation uses
`.generated/internal/app/dto/person`. Client and route annotations continue
generating the same endpoints and in-process request pipeline.

## Public client aliases, private DTO source

`client/person/person.go` imports the private DTO package and re-exports
client-facing types as Go aliases, for example:

```go
import persondto "grovehub/services/foundation/internal/app/dto/person"

type UpdatePerson = persondto.UpdatePerson
type CreatePersonRequest = persondto.CreatePersonRequest
type Person = persondto.Person
```

The interface retains `@Client` and `@RequestLine` annotations and uses the
same original concrete DTO types. Callers such as Identity import **only**
`grovehub/services/foundation/client/person` and use aliases such as
`person.CreatePersonRequest` and `person.Person`. All Foundation controllers,
validators, services, events and mappers import the private DTO definitions
inside the module; public aliases and private structs are the **same types**, not
copies or duplicated transport models.

## Public Person errors

Person error definitions remain owned in `internal/app/constants/error/`. The
client-facing aliases are grouped under `client/person/constants/error.go`,
so consumers import `grovehub/services/foundation/client/person/constants`
for error contracts, rather than mixing errors into the Person client interface.

## Resource ownership

The runtime migration implementation now lives in
`internal/resources/migrations/migrate.go` and accesses the SQL embedded by
`internal/resources/resources.go`. `foundation.Migrate` delegates to that
package. Foundation's OpenAPI specification lives in
`internal/resources/api/openapi/v1/openapi.yaml` and Swagger uses the
corresponding source path; the embedding also includes the YAML in production.

## Mirrored isolated unit tests

See `docs/engineering/backend/FOUNDATION_UNIT_TEST_GUIDELINES.md` for the
annotation syntax, mock policy, mirror conventions, and strict coverage gate.
Unit tests live under `internal/tests/unit-test/`, production source under
`internal/app/`, and domain-only fixtures and mocks under `internal/tests/`.
Shared core testing utilities are under `backend/libs/core/tests/`.

## Infrastructure-owned durable events

Foundation only constructs Person events and calls the public
`services/infrastructure/client/messaging.EventStore` from the same GORM
transaction as the Person mutation. The durable table, SQL repository,
dispatcher, NATS adapter, retries and 24-hour retention worker belong to
`services/infrastructure/internal/app/`. Foundation does not independently
start an outbox worker.

The existing Foundation `000002_create_event_outbox.sql` migration is retained
**only as historical Goose migration compatibility**. Infrastructure's
`000001_create_messaging.sql` atomically copies any old queued records to
`infrastructure.messaging` before dropping the legacy table. Never remove an
applied migration from Goose history; deploy by running all migrations in
service order (Foundation then Infrastructure) with messaging workers stopped.
