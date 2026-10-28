# GroveHub Mapper and DTO Guidelines

> **Status:** Mandatory GroveHub backend standard

Also read:

```text
docs/engineering/backend/CORE_LIBRARY_GUIDELINES.md
docs/engineering/backend/VALIDATOR_GUIDELINES.md
docs/engineering/backend/NORMALIZATION_AND_DRY_GUIDELINES.md
```

---

## 1. Purpose

DTOs define request/response shapes. Mappers perform real representation changes only. Validators validate and canonicalize request/service-input payloads.

```text
dto/request/<resource>.go
dto/response/<resource>.go
mappers/<resource>.go
mappers/impl/<resource>.go
validators/<resource>.go
validators/impl/<resource>.go
```

Do not create a mapper merely because a resource exists. Create it when there is an actual representation boundary.

---

## 2. Resource Cohesion

Resource-specific response/error metadata stays with that resource.

Example:

```go
type PersonErrorMetadata struct {
    GID string `json:"gid,omitempty"`
}
```

Its mapping belongs to the Person mapper boundary.

Do not create generic service-level files such as:

```text
dto/response/error.go
mappers/error.go
```

when the representation belongs only to one resource.

Shared **HTTP transport envelopes** are different: those belong to the HTTP core capability under `backend/libs/core/http/contracts`, not a repository-wide `core/dto` bucket.

---

## 3. Mapper Contract / Implementation

The public mapper package contains the interface:

```go
type PersonMapper interface {
    ToPersonGID(request request.PersonByGID) entity.GID
    ToPersonResponse(person entity.Person) response.Person
    ToPersonErrorMetadata(gid entity.GID) response.PersonErrorMetadata
}
```

The concrete implementation belongs under `mappers/impl` and is constructed by the owning module composition root.

Consumers depend on `mappers.PersonMapper`, not `mappers/impl.PersonMapper`.

---

## 4. Mapper vs Validator vs Service

```text
Mapper
    -> representation conversion

Validator
    -> request/service-input validation + canonicalization

Service helper
    -> pure service-owned construction/comparison/result mechanics

Service
    -> business orchestration
```

Validation does not require a second DTO. If normalization leaves the shape unchanged, return the existing request/service-input type.

Example:

```text
CreatePerson
    ↓ validator
normalized CreatePerson
```

Do not create `ValidatedCreatePerson` merely to duplicate `CreatePerson`.

For generic shared contracts, reuse the public capability types directly. Person search validation returns:

```text
coresearch.Query[entity.GID]
coresearch.CursorPayload
```

not a `ValidatedPersonSearch` wrapper.

---

## 5. Mapper Responsibilities

Mappers may:

```text
request DTO -> entity/value identifier
entity state -> response DTO
resource context -> resource error metadata
```

Mappers must not:

```text
query repositories
perform request validation
perform business validation
generate IDs
read clocks
select HTTP status codes
own lifecycle/workflow decisions
perform persistence queries
```

---

## 6. Entity and Persistence Mapping

For current GORM-backed GroveHub modules, an entity may also be the database/GORM model.

When the persisted record and entity are the same representation, do **not** create a duplicate persistence model and mapper merely to copy identical fields.

Foundation Person is the canonical example:

```text
entity.Person == GORM/database record
```

Therefore do not create:

```text
person_model.go
person_mapper.go   # solely for DB <-> entity copying
```

A service mapper may still exist for a real boundary such as:

```text
entity.Person -> dto/response.Person
```

Persistence/query-derived behavior belongs in the repository implementation, not in a mapper or entity.

---

## 7. Core Capability Usage

Mappers and validators consume stable core capability roots where relevant.

Correct:

```go
corenormalize "grovehub/libs/core/normalize"
coresearch "grovehub/libs/core/search"
```

Incorrect:

```text
grovehub/libs/core/helpers/normalize
grovehub/libs/core/search/contracts
```

Service code must not depend on a core capability's internal `contracts/helpers/mappers/validators` packages.

---

## 8. SRP / OCP

When mapping behavior is selected by a stable key such as contact type, a small registry of mapping functions may be preferable to a growing switch tree when it genuinely improves extensibility.

Do not introduce builders/factories for simple assignments.

Do not create a generic mapper dumping ground.

---

## 9. Hard Rules

```text
NO repository returning transport response DTOs
NO business logic in DTOs/mappers
NO request validation hidden inside mappers
NO concrete mapper dependency outside the owning composition root
NO resource error files without an independent/shared responsibility
NO mapper without a real representation change
NO duplicate validated DTO when the existing contract represents the data
NO wrapper around an existing core capability contract without new semantics
NO duplicate persistence model/mapper when entity already is the DB model
NO direct imports of core capability internals
NO builder for simple assignment
NO generic mapper dumping-ground package
YES mapper for real representation boundaries
YES HTTP envelopes under core/http/contracts
```

> **Guiding principle:** map only when representation changes; validate in validators; orchestrate in services; keep persistence/query behavior in repositories; and avoid duplicate models when the entity already represents the database record.
