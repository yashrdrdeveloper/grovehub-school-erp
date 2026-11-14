# Foundation module: public API, private implementation, embedded migrations

The public Go packages are `services/foundation` and
`services/foundation/client/person`. The client package exposes
Go type aliases for DTO definitions stored privately under `internal/app/dto/person/`.
The OpenAPI specification stays under `services/foundation/api/openapi/`.
Documentation is a repository asset rather than an importable Go contract.

```text
backend/services/foundation/
├── foundation.go
├── api/openapi/v1/openapi.yaml
├── client/person/
│   ├── person.go
│   ├── errors.go
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
│       ├── migrate.go
│       └── migrations/{000001_create_foundation_schema.sql,000002_create_event_outbox.sql}
└── docs/
```

## Migration contract

`resources.Migrations` uses Go's `//go:embed migrations/*.sql`.
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
