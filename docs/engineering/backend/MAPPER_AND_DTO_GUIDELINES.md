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

Entity construction and entity state transitions are not automatically mapper responsibilities.

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

Shared cross-service contract shapes may live under:

```text
backend/libs/core/contracts/<contract>/
```

For example, GroveHub search requests/results use `backend/libs/core/contracts/search` instead of duplicating the generic search envelope in each service.

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

Nested public value shapes are still DTOs. A nested value does not require an independent resource or independent identifier.

Example:

```text
Person response
└── Contacts
    ├── Emails[]
    └── Phones[]
```

---

## 4. Mapper Rules

`mappers/` owns pure representation transformation code.

A mapper may:

```text
convert request DTO values into entity/value types when no business decision is required
convert entity getter values into response DTOs
reshape entity collections into nested response structures
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
normalize values when normalization is a business/application rule
choose lifecycle/status transitions
decide version increments
choose HTTP status codes
mutate entity state
```

Those decisions remain in the owning service or transport layer.

---

## 5. Entity Access and Copy Semantics

When an entity encapsulates its fields, mappers must read it through entity getters.

Example:

```go
func ToPersonResponse(person entity.Person) response.Person {
    return response.Person{
        GID:        string(person.GetGID()),
        FirstName:  person.GetFirstName(),
        MiddleName: person.GetMiddleName(),
        LastName:   person.GetLastName(),
        Status:     string(person.GetStatus()),
        Version:    person.GetVersion(),
        CreatedAt:  person.GetCreatedAt(),
        UpdatedAt:  person.GetUpdatedAt(),
    }
}
```

The actual Foundation mapper also reads `person.GetContacts()` and groups those entity values into response `emails[]` and `phones[]` arrays. That is a representation change and therefore belongs in the mapper.

Incorrect outside the entity implementation:

```go
person.Status = ...
person.Version++
person.Contacts = append(...)
```

Entity setters use value receivers and return a new entity copy. They do not mutate through pointer receivers.

Collection getters/setters must preserve encapsulation. A collection getter returns a copy instead of exposing the entity's internal slice.

Example:

```go
updatedPerson := person.
    SetContacts(contacts).
    SetVersion(person.GetVersion() + 1).
    SetUpdatedAt(now)
```

The service chooses the new values. The entity setter only copies the entity, assigns the supplied value, and returns the copy.

---

## 6. Builder Rule

Builders are **not** the default GroveHub mapping pattern.

Do not create a builder for simple field assignment.

For simple entities, prefer a pure entity constructor that only assigns caller-supplied values:

```go
person := entity.NewPerson(
    gid,
    firstName,
    middleName,
    lastName,
    contacts,
    entity.PersonStatusActive,
    1,
    now,
    now,
)
```

The constructor must not decide defaults, normalize contacts, generate IDs, read the clock, validate business policy, or perform lifecycle transitions.

Introduce a builder only when there is real construction complexity, such as many optional fields or repeated complex assembly across multiple callers.

---

## 7. Service Usage

Current GroveHub application services may accept request DTOs/shared contracts and return response DTOs/shared result contracts at their application boundary.

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

Typical create flow:

```text
service validates/normalizes input and decides ID/defaults/time
   ↓
entity pure constructor
   ↓
repository
   ↓
mapper.To<Resource>Response
```

Typical update flow:

```text
repository returns entity
   ↓
service validates/normalizes and decides new business values
   ↓
entity copy-setter methods return updated copy
   ↓
repository.Update(updated copy)
   ↓
mapper.To<Resource>Response
```

Business decisions must remain visible in the service.

---

## 8. Error Metadata Mapping

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

## 9. Foundation Example

For the current Person implementation:

```text
entity/person.go
    private Person fields
    NewPerson(...)
    getters
    copy-setters for mutable fields

entity/contact.go
    private Contact value fields
    NewContact(...)
    getters

mappers/person.go
    ToPersonGID(...)
    ToPersonResponse(...)
        -> groups Contacts into emails[] / phones[]

mappers/error.go
    ToPersonErrorMetadata(...)
```

There is no `PersonBuilder` and no `ContactBuilder`.

Contact normalization, duplicate detection, primary validation, verification preservation, and search validation remain service responsibilities rather than mapper responsibilities.

---

## 10. Hard Rules

```text
NO direct entity -> API response serialization
NO repository -> response DTO
NO direct entity field reads outside entity package when getters exist
NO direct entity field mutation outside entity package
NO leaking mutable entity collection slices
NO pointer-receiver setters for copy-style entities
NO business decisions in entity setters
NO business decisions in mappers
NO ID/time generation in mappers
NO contact/search normalization in mappers
NO builder for simple field assignment
NO pattern abstraction without real complexity
NO ad-hoc error metadata maps outside mappers
NO generic mapper dumping files
```
