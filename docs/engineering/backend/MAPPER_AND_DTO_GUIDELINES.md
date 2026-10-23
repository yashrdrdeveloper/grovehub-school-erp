# GroveHub Mapper and DTO Guidelines

> **Status:** Mandatory GroveHub backend standard

## 1. Purpose

DTOs define request/response shapes. Mappers perform real representation changes only.

```text
dto/request/<resource>.go
dto/response/<resource>.go
mappers/<resource>.go
```

Do not create extra files solely because a resource has an error metadata shape.

## 2. Resource Cohesion

If an error metadata type belongs only to Person, keep it with the Person response DTO:

```go
type PersonErrorMetadata struct {
    GID string `json:"gid,omitempty"`
}
```

Likewise, its mapper belongs in `mappers/person.go`:

```go
func ToPersonErrorMetadata(gid entity.GID) response.PersonErrorMetadata
```

Do not create:

```text
dto/response/error.go
mappers/error.go
```

unless the file contains a genuinely shared multi-resource error representation/transformation.

## 3. Mapper Rules

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

## 4. Entity Encapsulation

Read state through getters. Copy-style entity changes use value-receiver setters selected by the service.

Never mutate private entity state through mapper code.

## 5. SRP / OCP

When mapping behavior is selected by a stable resource key (for example contact type), prefer a small registry of mapping functions over a growing switch tree when it improves extensibility.

Do not introduce builders/factories for simple assignments.

## 6. Hard Rules

```text
NO direct entity serialization
NO repository returning response DTOs
NO business logic in DTOs/mappers
NO separate resource error files without an independent/shared responsibility
NO builder for simple assignment
NO generic mapper dumping-ground file
```
