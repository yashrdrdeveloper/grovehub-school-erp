# GroveHub Normalization and DRY Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** all backend services, shared libraries, search validation, validators, and future modules

Also read:

```text
docs/engineering/backend/CORE_LIBRARY_GUIDELINES.md
```

---

## 1. Principle

GroveHub shares stable, domain-neutral technical normalization once.

Normalization is its own core capability:

```text
backend/libs/core/normalize/
├── normalize.go
└── helpers/
    └── normalize.go
```

Consumers import only:

```go
import corenormalize "grovehub/libs/core/normalize"
```

Do not import `normalize/helpers` from services and do not create a repository-wide `libs/core/helpers/normalize` package.

Resource request validation and domain-specific error mapping remain under the owning service validator:

```text
validators/<resource>.go
validators/impl/<resource>.go
validators/helpers/<resource>.go
```

Business policy remains inside the owning service.

---

## 2. Current Shared Normalizers

The `core/normalize` façade currently exposes:

```text
normalize.Email(value)
normalize.Phone(value)
normalize.PhoneCharacters(value)
normalize.UUIDv7(value)
```

It also exposes stable generic technical errors such as:

```text
ErrInvalidEmail
ErrInvalidPhone
ErrInvalidUUIDv7
```

These functions own technical behavior only:

```text
trimming
canonical case
email syntax/canonicalization
allowed phone characters
E.164-compatible phone shape
UUID parsing/version validation
canonical output formatting
```

---

## 3. Domain Error Boundary

`libs/core` must never import service-specific error definitions.

Correct:

```text
Identity validator/helper
    ↓
core/normalize.Email(...)
    ↓
generic normalization failure
    ↓
Identity maps to IDENTITY_ACCOUNT_INVALID
```

```text
Foundation validator/helper
    ↓
core/normalize.Email(...)
    ↓
generic normalization failure
    ↓
Foundation maps to FOUNDATION_PERSON_INVALID_CONTACT
```

Incorrect:

```text
core/normalize/helpers
    ↓
imports foundation/constants/error
```

or:

```text
core/normalize/helpers
    ↓
imports identity/constants/error
```

Shared normalization knows technical formats, not domain semantics.

---

## 4. What Belongs in `core/normalize`

Promote logic when all are true:

- the input represents the same technical concept across services;
- canonical output must be identical across services;
- the rule does not depend on repositories, authorization, workflow, or domain state;
- more than one consumer needs it, or it is an established repository-wide primitive.

Examples:

```text
email canonicalization
phone canonicalization
UUIDv7 validation/canonicalization
```

Do not move domain-specific validation to core merely to reduce line count.

---

## 5. What Belongs in Service Validators

Validators own request/service-input policy around shared normalization.

Examples:

```text
Identity account email is required
optional account phone may be empty
Foundation allows at most N contacts per type
only one primary contact per type
search filter values are valid for their operators
cursor values match the requested resource sort
```

A validator may call `core/normalize`, then translate a generic normalization failure into the owning domain's error.

Do not repeat these checks inside service methods.

---

## 6. What Must Remain Service-Owned

Do not move stateful/business decisions into normalization or request validators:

```text
whether an email belongs to another account
whether a contact belongs to another Person
whether an account may transition lifecycle state
whether a caller is authorized
repository-backed uniqueness/existence rules
workflow decisions
```

Those require business state or I/O and remain in services/repositories.

---

## 7. Search Normalization

Exact search values representing stored canonical values must reuse `core/normalize`.

Examples:

```text
exact email search -> normalize.Email
exact phone search -> normalize.Phone
GID/cursor UUID     -> normalize.UUIDv7
```

Partial search may use a lower-level public capability operation such as `normalize.PhoneCharacters` when full-value validation would reject valid `contains` / `startsWith` input.

Do not copy phone-character stripping into resource validators.

---

## 8. ID Generation Is a Separate Capability

Generation and normalization are different responsibilities.

Generate IDs through:

```go
import coreid "grovehub/libs/core/id"

value, err := coreid.NewUUIDv7()
```

Validate/canonicalize existing IDs through:

```go
import corenormalize "grovehub/libs/core/normalize"

value, err := corenormalize.UUIDv7(input)
```

Do not place either under a generic `core/helpers` package.

---

## 9. Helper Placement by Owner

```text
core/normalize/helpers
    -> implementation mechanics of shared normalization

validators/helpers/<resource>.go
    -> resource payload validation/canonicalization policy

services/helpers/<resource>.go
    -> pure helpers used by service orchestration
```

The presence of the word “helper” does not determine ownership. The capability/layer that owns the behavior determines placement.

---

## 10. Hard Rules

```text
NO libs/core/helpers/normalize
NO service importing core/normalize/helpers directly
NO duplicate email normalization across services
NO duplicate phone normalization across services
NO duplicate UUIDv7 validation across services
NO service error definitions imported by libs/core
NO repository/network calls from shared normalizers
NO authorization/business policy in shared normalizers
NO repository-backed rules in request validators
NO payload normalization duplicated in service implementations
NO giant generic utils/common/helpers package
YES core/normalize as the public normalization façade
YES core/id as the public ID-generation façade
```

> **Guiding principle:** share technical canonicalization once behind the capability façade; keep request policy with validators and business meaning/state with services.
