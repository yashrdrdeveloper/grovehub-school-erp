# GroveHub Mapper and DTO Guidelines

> **Status:** Mandatory GroveHub backend standard
> **Applies to:** all backend services/modules
> **Precedence:** when older examples conflict with this document, this document wins together with `GROVEHUB_BACKEND_STRUCTURE.md`.

---

## 1. Purpose

All object-to-object transformation belongs in the owning service `mappers/` package.

Examples:

```text
request DTO -> entity/value
entity -> response DTO
existing entity -> updated entity construction
error context -> error metadata DTO
persistence model -> entity, only when a separate persistence model actually exists
```

Do not spread field-by-field mapping across controllers, services, repositories, or routes.

---

## 2. Canonical Service Shape

```text
backend/services/<service>/
├── dto/
│   ├── request/
│   └── response/
├── entity/
├── mappers/
├── repository/
│   └── postgres/
├── services/
├── controller/
└── routes/
```

Use only folders/files required by the current service.

---

## 3. DTO Rules

`dto/request/` contains request/application-input shapes.

`dto/response/` contains response/application-output shapes and error metadata shapes.

DTOs:

```text
contain data only
contain JSON/transport tags when required
contain no business logic
contain no repository/database access
contain no UUID generation
contain no lifecycle/status decisions
```

Entities are never serialized directly as public API responses.

Repositories never return response DTOs.

---

## 4. Mapper Rules

`mappers/` owns pure transformation/construction code.

A mapper may:

```text
convert request DTO values into entity/value types
convert entity values into response DTOs
copy one entity into a builder for an updated entity
construct error metadata from known service context
map a real persistence model into an entity
```

A mapper must not:

```text
query repositories/databases
call network services
read environment variables
generate IDs
read the current clock
perform authorization
validate business policy
choose lifecycle/status transitions
decide version increments
choose HTTP status codes
```

Those decisions remain in the owning service or transport layer.

---

## 5. Builder Pattern

When an entity/DTO has multiple fields or is constructed in more than one use case, prefer a focused builder in the mapper package.

Example:

```go
person := mappers.NewPersonBuilder().
    WithGID(gid).
    WithStatus(entity.PersonStatusActive).
    WithVersion(1).
    WithCreatedAt(now).
    WithUpdatedAt(now).
    Build()
```

The builder only assigns values supplied by the caller.

The builder must not silently choose defaults such as:

```text
ACTIVE
version = 1
current time
new UUID
```

The service owns those business/default decisions and passes the values into the builder.

For updates, prefer copying an existing entity and changing only values chosen by the service:

```go
updated := mappers.NewPersonBuilderFrom(person).
    WithStatus(entity.PersonStatusDeactivated).
    WithVersion(person.Version + 1).
    WithUpdatedAt(now).
    Build()
```

---

## 6. Service Usage

Current GroveHub application services may accept request DTOs and return response DTOs at their application boundary.

Inside the service:

```text
request DTO
   ↓ mapper
entity/value
   ↓ business logic
repository
   ↓ entity
mapper
   ↓
response DTO
```

Business decisions must remain visible in the service.

Correct:

```text
service chooses ACTIVE
service chooses version 1
service chooses now
mapper builds Person with those supplied values
```

Incorrect:

```text
mapper decides ACTIVE
mapper decides version 1
mapper calls time.Now()
mapper generates UUID
```

---

## 7. Error Metadata Mapping

Error definitions remain under:

```text
constants/error/<resource>.go
```

Error metadata DTOs live under:

```text
dto/response/
```

Conversion of service context into error metadata belongs in:

```text
mappers/error.go
```

Example flow:

```text
GID/entity/service context
      ↓
mappers.ToPersonErrorMetadata(...)
      ↓
coreerrors.Error metadata
      ↓
controller
      ↓
shared HTTP error response envelope
```

Do not assemble `errorMetaData` maps ad hoc in controllers or services.

---

## 8. Foundation Example

```text
backend/services/foundation/
├── dto/
│   ├── request/
│   │   └── person.go
│   └── response/
│       ├── person.go
│       └── error.go
├── entity/
│   ├── gid.go
│   └── person.go
├── mappers/
│   ├── person.go
│   └── error.go
├── repository/
│   ├── person.go
│   └── postgres/
│       └── person.go
└── services/
    └── person.go
```

Foundation tests are added together after the Foundation flow is complete for the current implementation phase.

---

## 9. Hard Rules

```text
NO direct entity -> API response serialization
NO repository -> response DTO
NO field-by-field object conversion in controllers
NO field-by-field object conversion in services when a mapper exists
NO business decisions in mappers
NO ID/time generation in builders
NO ad-hoc error metadata maps outside mappers
NO generic mapper dumping files; group mapping by resource/responsibility
```
