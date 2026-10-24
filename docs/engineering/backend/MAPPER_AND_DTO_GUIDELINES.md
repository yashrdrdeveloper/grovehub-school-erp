# GroveHub Mapper and DTO Guidelines

> **Status:** Mandatory GroveHub backend standard

## 1. Purpose

DTOs define request/response shapes. Mappers perform real representation changes only.

```text
dto/request/<resource>.go
dto/response/<resource>.go
mappers/<resource>.go
mappers/impl/<resource>.go
```

`mappers/<resource>.go` defines the public mapper interface. The concrete mapper implementation belongs under `mappers/impl/`.

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

## 4. Mapper Rules

Mappers may:

```text
request DTO -> entity/value
entity getters -> response DTO
resource context -> resource error metadata
```

Mappers must not:

```text
query repositories
perform business validation
generate IDs
read clocks
select HTTP status
own lifecycle decisions
```

## 5. Entity Encapsulation

Read state through getters. Copy-style entity changes use value-receiver setters selected by the service.

Never mutate private entity state through mapper code.

## 6. SRP / OCP

When mapping behavior is selected by a stable resource key (for example contact type), prefer a small registry of mapping functions over a growing switch tree when it improves extensibility.

Do not introduce builders/factories for simple assignments.

## 7. Hard Rules

```text
NO direct entity serialization
NO repository returning response DTOs
NO business logic in DTOs/mappers
NO concrete mapper dependency outside the owning composition root
NO separate resource error files without an independent/shared responsibility
NO mapper without a real representation change
NO builder for simple assignment
NO generic mapper dumping-ground file
```
