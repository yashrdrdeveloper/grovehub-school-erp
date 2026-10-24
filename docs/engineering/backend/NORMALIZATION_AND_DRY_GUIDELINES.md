# GroveHub Normalization and DRY Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** all backend services, shared libraries, search helpers, validators, and future modules

## 1. Principle

GroveHub follows DRY for stable technical normalization rules that are used by more than one domain.

A service must not duplicate generic normalization logic merely to translate the result into its own domain error.

Shared technical normalization belongs under:

```text
backend/libs/core/helpers/normalize/
```

Resource request validation and domain-specific error mapping belong under the owning service's validator boundary:

```text
validators/<resource>.go
validators/impl/<resource>.go
validators/helpers/<resource>.go
```

Business policy remains inside the owning service.

## 2. Current shared normalizers

```text
backend/libs/core/helpers/normalize/
├── email.go
├── phone.go
└── uuid.go
```

The package currently owns:

```text
normalize.Email(value)
normalize.Phone(value)
normalize.PhoneCharacters(value)
normalize.UUIDv7(value)
```

These functions own only stable technical behavior such as trimming, canonical casing, syntax parsing, allowed phone characters, E.164-compatible phone shape, UUID parsing, UUID version checks, and canonical string output.

## 3. Domain error boundary

`libs/core` must never import a service's error definitions.

Correct:

```text
Identity validator/helper
    ↓
normalize.Email(...)
    ↓
shared generic normalization error
    ↓
Identity maps failure to IDENTITY_ACCOUNT_INVALID
```

```text
Foundation validator/helper
    ↓
normalize.Email(...)
    ↓
shared generic normalization error
    ↓
Foundation maps failure to FOUNDATION_PERSON_INVALID_CONTACT
```

Incorrect:

```text
libs/core/helpers/normalize/email.go
    ↓
imports foundation/constants/error
```

or:

```text
libs/core/helpers/normalize/email.go
    ↓
imports identity/constants/error
```

Shared helpers know technical formats, not domain semantics.

## 4. What belongs in shared normalization

Promote logic to `libs/core/helpers/normalize` when all of the following are true:

- the input represents the same technical concept across services;
- the canonical output must be identical across services;
- the rule does not depend on domain state, repositories, authorization, or workflow;
- at least two current consumers need the behavior, or the behavior is an established repository-wide technical primitive.

Examples:

```text
email canonicalization
phone canonicalization
UUIDv7 validation/canonicalization
```

## 5. What belongs in validators

Validators own request/service-input policy around shared normalizers.

Examples:

```text
email is required for an Identity account
phone may be empty for an optional account phone
Foundation allows at most N contacts per type
only one primary contact is allowed per type
search filter values are valid for their operators
cursor fields and values match the requested sort
```

A validator may call `libs/core/helpers/normalize`, then translate a generic normalization failure into the owning domain's invalid-request error.

Do not repeat these checks directly inside service methods.

## 6. What must remain service-owned

Do not move these into shared normalizers or request validators:

```text
whether an email already belongs to another account
whether a Foundation contact belongs to another Person
whether an account may transition lifecycle state
whether the caller is authorized
authorization or lifecycle rules
repository-backed uniqueness/existence checks
```

Those decisions require business state or I/O and remain in the service.

## 7. Search behavior

Search validators must reuse the same canonical normalization when an exact search value represents the same stored value.

Examples:

```text
exact email search -> normalize.Email
exact phone search -> normalize.Phone
GID/cursor UUID -> normalize.UUIDv7
```

Partial-search behavior may use a lower-level shared lexical helper such as `normalize.PhoneCharacters` when full-value validation would be incorrect for `contains` or `startsWith` queries.

Do not copy phone character-stripping logic into validator helpers.

## 8. Helper placement by consumer

Place helpers with the layer that actually consumes them.

```text
validators/helpers/<resource>.go
    -> payload validation and normalization helpers

services/helpers/<resource>.go
    -> helpers used by business/service orchestration
```

Do not keep a function under `services/helpers` merely because it historically lived there.

If only validators use it, move it to `validators/helpers`.

If both layers need a stable technical primitive across domains, consider `libs/core` only when the behavior is genuinely domain-neutral.

## 9. Interface / implementation consistency

The same boundary principle used for shared normalization applies to service architecture:

```text
services/<resource>.go
    -> public service interface
services/impl/<resource>.go
    -> concrete implementation

validators/<resource>.go
    -> public validator interface
validators/impl/<resource>.go
    -> concrete implementation
validators/helpers/<resource>.go
    -> validator-only resource helpers

mappers/<resource>.go
    -> public mapper interface when a real mapping boundary exists
mappers/impl/<resource>.go
    -> concrete implementation

repository/<resource>.go
    -> repository interface
repository/impl/<resource>.go
    -> persistence implementation
```

Do not create empty mapper abstractions when there is no current representation change. Introduce the mapper boundary when DTO/entity or other real representation mapping exists.

## 10. Hard rules

```text
NO duplicate email normalization in Foundation and Identity
NO duplicate phone normalization in Foundation and Identity
NO duplicate UUIDv7 validation across services
NO service error definitions imported by libs/core
NO repository calls from shared normalizers
NO repository/network calls from validators
NO authorization/business policy in shared normalizers
NO repository-backed business rule inside request validators
NO payload normalization duplicated inside service implementations
NO validator-only helpers left under services/helpers
NO giant generic utils/common package
NO domain-specific behavior promoted to core just to reduce line count
```

> **Guiding principle:** share stable technical mechanics once; validate/canonicalize payloads in validators; keep business meaning, policy, and I/O with the owning service.
