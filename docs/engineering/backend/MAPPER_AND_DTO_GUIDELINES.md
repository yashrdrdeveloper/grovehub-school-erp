# GroveHub Mapper and DTO Guidelines

> **Status:** Mandatory GroveHub backend standard

## 1. Purpose

DTOs define request/response shapes. Mappers perform real representation changes only. Validators validate and normalize request/service-input payloads.

```text
dto/request/<resource>.go
dto/response/<resource>.go
mappers/<resource>.go
mappers/impl/<resource>.go
validators/<resource>.go
validators/impl/<resource>.go
```

`mappers/<resource>.go` defines the public mapper interface. The concrete mapper implementation belongs under `mappers/impl/`.

`validators/<resource>.go` defines the public validator interface. The concrete validator implementation belongs under `validators/impl/`.

Do not create a mapper for a resource that has no real representation change yet. Do not create extra files solely because a resource has an error metadata shape.

## 2. Resource Cohesion

If an error metadata type belongs only to Person, keep it with the Person response DTO:

```go
type PersonErrorMetadata struct {
    GID string `json:"gid,omitempty"`
}
```

Its mapper contract belongs in `mappers/person.go` and the implementation belongs in `mappers/impl/person.go`.

Do not create:

```text
dto/response/error.go
mappers/error.go
```

unless the file contains a genuinely shared multi-resource error representation/transformation.

## 3. Mapper Contract / Implementation Rule

The public mapper package contains interfaces only:

```go
type PersonMapper interface {
    ToPersonGID(request request.PersonByGID) entity.GID
    ToPersonResponse(person entity.Person) response.Person
    ToPersonErrorMetadata(gid entity.GID) response.PersonErrorMetadata
}
```

The concrete implementation belongs under `mappers/impl` and is constructed only by the owning service composition root.

Consumers depend on `mappers.PersonMapper`, never directly on `mappers/impl.PersonMapper`.

## 4. Mapper vs Validator

Use the mapper for representation conversion and the validator for payload correctness/canonicalization.

```text
Mapper
    request/response representation conversion

Validator
    request/service-input validation
    normalization/canonicalization

Service
    business orchestration
```

Validation does not require a second DTO. When a request remains the same shape after normalization, the validator returns the existing request DTO directly.

Example:

```text
CreatePerson DTO
    ↓
validator
    ↓
normalized CreatePerson DTO
```

Do not create `ValidatedCreatePerson` merely to duplicate the fields of `CreatePerson`.

For generic infrastructure contracts, reuse the existing generic type directly. Example: Person search validation returns `coresearch.Query[entity.GID]` and `coresearch.CursorPayload`, not a `ValidatedPersonSearch` wrapper.

For representation changes:

```text
PersonByGID DTO -> entity.GID
    -> mapper

validate UUIDv7 format
    -> validator

load Person and decide not-found behavior
    -> service
```

Do not hide payload validation inside mapper code simply because a mapper already touches the field.

## 5. Mapper Rules

Mappers may:

```text
request DTO -> entity/value
entity getters -> response DTO
resource context -> resource error metadata
```

Mappers must not:

```text
query repositories
perform request validation
perform business validation
generate IDs
read clocks
select HTTP status
own lifecycle decisions
```

## 6. Validator interaction

A service may use both a mapper and validator. The order depends on the responsibility:

```text
request DTO
   ↓
validator for payload validity/canonicalization
   ↓
mapper/service helper for a real representation conversion when needed
   ↓
service business orchestration
   ↓
mapper for response conversion
```

For a primitive representation that must be converted before validation, conversion may happen first, as with `PersonByGID -> entity.GID`.

Do not duplicate normalization in mapper and validator. Shared technical normalization belongs in `libs/core/helpers/normalize` and is consumed by the validator.

See:

```text
docs/engineering/backend/VALIDATOR_GUIDELINES.md
docs/engineering/backend/NORMALIZATION_AND_DRY_GUIDELINES.md
```

## 7. Entity Encapsulation

Read state through getters. Copy-style entity changes use value-receiver setters selected by the service.

Never mutate private entity state through mapper code.

## 8. SRP / OCP

When mapping behavior is selected by a stable resource key (for example contact type), prefer a small registry of mapping functions over a growing switch tree when it improves extensibility.

Do not introduce builders/factories for simple assignments.

## 9. Hard Rules

```text
NO direct entity serialization
NO repository returning response DTOs
NO business logic in DTOs/mappers
NO request validation hidden inside mappers
NO concrete mapper dependency outside the owning composition root
NO separate resource error files without an independent/shared responsibility
NO mapper without a real representation change
NO duplicate validated DTO when the existing request/service-input contract already represents the data
NO wrapper around an existing generic core contract without new semantics
NO builder for simple assignment
NO generic mapper dumping-ground file
```
