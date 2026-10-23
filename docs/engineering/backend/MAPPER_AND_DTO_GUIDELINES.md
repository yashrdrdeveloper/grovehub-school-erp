# GroveHub Mapper and DTO Guidelines

> **Status:** Mandatory GroveHub backend standard
> **Applies to:** all backend services/modules
> **Precedence:** when older examples conflict with this document, this document wins together with `GROVEHUB_BACKEND_STRUCTURE.md`.

---

## 1. Purpose

Use `mappers/` only for real representation changes.

Examples:

```text
request DTO -> entity/value
entity -> response DTO
error context -> error metadata DTO
persistence model -> entity, only when a separate persistence model actually exists
```

Do not add mapping abstractions when no representation change exists.

Creating a new entity from values already decided by the service is not automatically a mapper responsibility.

Changing business state on an existing entity is business logic and belongs in `services/`.

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

`mappers/` owns pure representation transformation code.

A mapper may:

```text
convert request DTO values into entity/value types
convert entity values into response DTOs
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

## 5. Builder Rule

Builders are **not** the default GroveHub mapping pattern.

Do not create a builder for a simple entity merely because it has several fields.

For simple business-owned entity construction, construct the entity directly in the service where the values are decided:

```go
person := entity.Person{
    GID:       gid,
    Status:    entity.PersonStatusActive,
    Version:   1,
    CreatedAt: now,
    UpdatedAt: now,
}
```

For a business state transition, update the entity directly in the service:

```go
person.Status = entity.PersonStatusDeactivated
person.Version++
person.UpdatedAt = now
```

Introduce a builder only when there is a real complexity that justifies it, such as:

```text
many optional fields
repeated complex construction across multiple callers
ordered/conditional assembly that remains non-business transformation logic
construction that is otherwise difficult to read or maintain
```

A builder must never be added only to satisfy a pattern.

If a builder is introduced later, it still must not own business defaults, ID generation, clocks, validation policy, or lifecycle decisions.

---

## 6. Service Usage

Current GroveHub application services may accept request DTOs and return response DTOs at their application boundary.

Typical read flow:

```text
request DTO
   ↓ mapper when conversion is needed
entity/value
   ↓ business logic
repository
   ↓ entity
mapper
   ↓
response DTO
```

Typical create flow when there is no client-controlled entity payload:

```text
service decides ID/defaults/time
   ↓
service constructs entity directly
   ↓
repository
   ↓
mapper.To<Resource>Response
```

Business decisions must remain visible in the service.

Correct:

```text
service chooses ACTIVE
service chooses version 1
service chooses now
service constructs/updates Person
mapper converts Person -> response DTO
```

Incorrect:

```text
mapper decides ACTIVE
mapper decides version 1
mapper calls time.Now()
mapper generates UUID
builder added only to wrap field assignment
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

For the current Person implementation, `mappers/person.go` only needs real conversions such as:

```go
ToPersonGID(...)
ToPersonResponse(...)
```

It does not need a `PersonBuilder`.

Foundation tests are added together after the Foundation flow is complete for the current implementation phase.

---

## 9. Hard Rules

```text
NO direct entity -> API response serialization
NO repository -> response DTO
NO field-by-field representation conversion in controllers
NO business decisions in mappers
NO ID/time generation in mappers
NO builder for simple field assignment
NO pattern abstraction without real complexity
NO ad-hoc error metadata maps outside mappers
NO generic mapper dumping files; group mapping by resource/responsibility
```
