# GroveHub Normalization and DRY Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** all backend services, shared libraries, search helpers, validation helpers, and future modules

## 1. Principle

GroveHub follows DRY for stable technical normalization rules that are used by more than one domain.

A service must not duplicate generic normalization logic merely to translate the result into its own domain error.

Shared technical normalization belongs under:

```text
backend/libs/core/helpers/normalize/
```

Domain-specific policy and domain-specific error mapping remain inside the owning service.

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
Identity helper
    ↓
normalize.Email(...)
    ↓
shared generic normalization error
    ↓
Identity maps failure to IDENTITY_ACCOUNT_INVALID
```

```text
Foundation helper
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

## 5. What must remain service-owned

Do not move these into the shared normalizer:

```text
whether email is required for an Identity account
whether phone is optional
maximum number of Foundation contacts
whether a Person may have multiple primary contacts
whether an account may change its email
which domain error code is returned
authorization or lifecycle rules
repository-backed uniqueness checks
```

For example, `normalize.Phone` validates/canonicalizes a phone number. Identity decides that an empty phone means "no optional phone" before calling the shared normalizer. Foundation decides that a submitted contact value must be a valid contact.

## 6. Search behavior

Search helpers must reuse the same canonical normalization when an exact search value represents the same stored value.

Examples:

```text
exact email search -> normalize.Email
exact phone search -> normalize.Phone
GID/cursor UUID -> normalize.UUIDv7
```

Partial-search behavior may use a lower-level shared lexical helper such as `normalize.PhoneCharacters` when full-value validation would be incorrect for `contains` or `startsWith` queries.

Do not copy phone character-stripping logic into search helpers.

## 7. Interface / implementation consistency

The same boundary principle used for shared normalization applies to service architecture:

```text
services/<resource>.go
    -> public service interface
services/impl/<resource>.go
    -> concrete implementation

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

## 8. Hard rules

```text
NO duplicate email normalization in Foundation and Identity
NO duplicate phone normalization in Foundation and Identity
NO duplicate UUIDv7 validation across services
NO service error definitions imported by libs/core
NO repository calls from shared normalizers
NO authorization/business policy in shared normalizers
NO giant generic utils/common package
NO domain-specific behavior promoted to core just to reduce line count
```

> **Guiding principle:** share stable technical mechanics once; keep domain meaning, policy, and error semantics with the owning service.
