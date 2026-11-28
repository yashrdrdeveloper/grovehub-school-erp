# Foundation — Implementation Architecture

**Owner:** Person identity aggregate and contact associations.

## Public contract and bootstrap

```text
services/foundation/
├── main.go                                 # typed Bootstrap + compatibility adapters
├── client/person/person.go                 # @Client, typed public DTO aliases
├── client/person/constants/error.go        # client-visible error aliases
└── internal/
    ├── app/
    │   ├── bootstrap/{bootstrap.go,impl/}   # repository/service/routes/events/client composition
    │   ├── dto/person/{operations,params,queryparams,request,response}.go
    │   ├── entity/{person,contact}.go
    │   ├── validators/person/
    │   ├── mappers/person/
    │   ├── repository/{person,contact}/
    │   ├── services/person/{impl,operations,helpers}/
    │   ├── controller/person/protected/
    │   ├── routes/http.go
    │   └── events/producers/person/
    ├── resources/{resources.go,migrations/,api/openapi/v1/}
    └── tests/{unit-test/,mocks/,fixtures/}
```

`foundation.New(logger, db, eventStore)` composes the typed public client, request validation, Person and Contact repositories, the Person service operations and the event producer. `foundation.Migrate(ctx, db)` applies only the embedded Foundation schema migrations. Generated annotated methods and routes live under an ignored service `.generated/` overlay rather than the maintained source tree.

## Person aggregate and SRP

- `entity/person.go`: UUIDv7 `gid`, lifecycle status, version, normalized search projection, idempotency identity and timestamps.
- `entity/contact.go`: separate persisted contacts with type (email/phone), normalized value, primary/verified flags, timestamps and association to Person GID.
- `repository/person/`: base person CRUD, idempotency lookup, safe search. `@Searchable` uses an explicit resource-owned query config.
- `repository/contact/`: contact batch retrieval/persistence. No independent contact business lifecycle.
- `validators/person/`: name/contact normalization, PATCH presence, search field/operator rules. No SQL.
- `mappers/person/`: entity-to-response and request-to-command representation changes.
- `services/person/impl/`: wires operation implementations, not the algorithms themselves.
- `services/person/operations/`: `create`, `getbygid`, `update`, `search`, `deactivate`. Multistep operations own flows/strategies; simple reads avoid artificial chains.
- `controller/person/protected/`: generated route contract and actual request/response adapter.
- `events/producers/person/`: changed event type/pipeline; transport delivery stays in Infrastructure.

## Atomicity and conflicts

The Create, Update and Deactivate flows operate through shared PostgreSQL transaction scope when writing associated rows, idempotency state and the outbox event. Contact normalized values and primary flags must not violate the aggregate's constraints. Requests returning the same idempotent result must not emit another event. Version predicates protect updates; a conflict is reported with a stable error, not an unseen overwrite.

PATCH has three states for optional fields (omitted, explicit null, value). Replacing contact collections is a whole-set decision, not implicit merge. Verified status is Foundation-managed, not supplied by arbitrary clients.

## External boundaries

The public HTTP endpoints currently live under `/foundation/protected/persons`. **The URL tier is a routing label, not enforced IAM.** Do not expose those endpoints on untrusted networks without verified authentication, authorization and tenant policy. Sibling modules import `client/person`, never `internal/app`.

## Migrations and deployment

Embeds `internal/resources/migrations/000001_create_foundation_schema.sql`. An older `000002_create_event_outbox.sql` remains immutable for historical Goose compatibility and data migration; Infrastructure now owns the live messaging table. Deployment must apply the required schema migrations in the correct service order with messaging workers stopped during one-time legacy queue transfer. Never edit applied SQL history.

## Tests

Tests mirror app ownership under `internal/tests/unit-test` for controllers, event producer/pipeline, operations, helpers and generated client/route constraints. PostgreSQL integration tests should cover contact uniqueness, transaction rollback, optimistic version and outbox atomicity; a mock alone does not establish these guarantees.
