# Foundation module: public API, private implementation, embedded migrations

Public Go packages are `services/foundation` and
`services/foundation/client/person`, including its `dto/` child package.
The OpenAPI specification stays under `services/foundation/api/openapi/`.
Documentation is a repository asset rather than an importable Go contract.

```text
backend/services/foundation/
├── foundation.go
├── api/openapi/v1/openapi.yaml
├── client/person/
│   ├── person.go
│   ├── errors.go
│   └── dto/{operations,params,queryparams,request,response}.go
├── internal/
│   ├── app/
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
overlay. Public Person DTO PATCH generation uses
`.generated/client/person/dto`. Client and route annotations continue
generating the same endpoints and in-process request pipeline.
