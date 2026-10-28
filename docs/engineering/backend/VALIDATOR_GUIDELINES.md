# GroveHub Validator Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** request/service-input validation in all backend services and modules

Also read:

```text
docs/engineering/backend/CORE_LIBRARY_GUIDELINES.md
docs/engineering/backend/NORMALIZATION_AND_DRY_GUIDELINES.md
```

---

## 1. Purpose

Validators keep payload validation and canonicalization out of business services.

```text
Request DTO / service input
        ↓
Validator interface
        ↓
Validator implementation
        ↓
Same normalized contract or existing generic contract
        ↓
Service implementation
        ↓
Business rules + repository orchestration
```

A validator answers whether input is valid and what its canonical values are. A service decides what business action occurs with those valid values.

---

## 2. Canonical Structure

```text
validators/
├── <resource>.go
├── impl/
│   └── <resource>.go
└── helpers/
    └── <resource>.go
```

`validators/<resource>.go` defines the interface.

`validators/impl/<resource>.go` owns validation orchestration.

`validators/helpers/<resource>.go` owns pure resource-specific validation/normalization policy.

Concrete validator implementations are constructed only by the owning module composition root.

---

## 3. Reuse Existing Contracts

Validation does not create a parallel data model by itself.

If validation only trims, normalizes, canonicalizes, or checks values while the shape remains the same, return the existing type.

Examples:

```text
CreatePerson -> normalized CreatePerson
UpdatePerson -> normalized UpdatePerson
CreateAccountInput -> normalized CreateAccountInput
UpdateAccountInput -> normalized UpdateAccountInput
```

For shared generic search:

```text
coresearch.Request
      ↓ validator
coresearch.Query[entity.GID]
coresearch.CursorPayload
```

Do not create:

```text
ValidatedCreatePerson
ValidatedUpdatePerson
ValidatedCreateAccount
ValidatedUpdateAccount
ValidatedPersonSearch
```

unless validation genuinely produces new information that cannot be represented by an existing DTO, service input, entity value, or generic core contract.

---

## 4. Validator Responsibilities

Validators may own:

```text
required/optional field checks
string length checks
request-local relationship checks
request normalization
UUID/email/phone canonicalization through core/normalize
contact-list shape checks
search field/operator/value validation
cursor payload validation
canonicalization of existing request/service-input contracts
```

Validators must not:

```text
query repositories
call external services
perform authorization
choose business workflows
make lifecycle transitions
generate IDs
read clocks
persist entities
publish events
select HTTP status codes
```

Repository-backed existence/uniqueness checks are service/business orchestration, not request validation.

---

## 5. Core Capability Usage

Validators consume stable core capability façades only.

Correct:

```go
import corenormalize "grovehub/libs/core/normalize"
import coresearch "grovehub/libs/core/search"
```

Incorrect:

```text
grovehub/libs/core/helpers/normalize
grovehub/libs/core/normalize/helpers
grovehub/libs/core/search/contracts
grovehub/libs/core/search/helpers
```

Core subpackages such as `contracts`, `helpers`, `mappers`, and `validators` are implementation details of the owning capability.

Current shared normalization operations include:

```text
corenormalize.Email
corenormalize.Phone
corenormalize.PhoneCharacters
corenormalize.UUIDv7
```

The shared normalizer returns a generic technical failure. The owning resource validator maps that failure to its domain/application error contract.

---

## 6. Validator vs Mapper vs Service Helper

```text
Mapper
    -> representation conversion

Validator
    -> payload validation + canonicalization

Service helper
    -> pure service-owned construction/comparison/result mechanics

Service
    -> business orchestration
```

Example:

```text
CreatePerson.Contacts
    ↓
PersonValidator
    -> count/duplicate/primary checks
    -> canonical email/phone values
    -> normalized request.Contacts
    ↓
Service helper
    -> entity.Contact values
    ↓
Service
    -> ownership/conflict checks + persistence workflow
```

A mapper must not become a validation layer, and a validator must not produce transport responses.

---

## 7. Validator Helpers

Place a function under `validators/helpers` when its consumers are validator implementations and its responsibility is validation/canonicalization.

Examples:

```text
NormalizePersonNames
NormalizeContacts
NormalizePersonSearchRequest
NormalizeAccountId
NormalizePersonGid
NormalizeEmail
NormalizePhone
```

Do not leave validator-only functions under `services/helpers`.

Service helpers may own behavior such as:

```text
entity/value construction from already validated data
preserving existing state while rebuilding values
entity comparison used by update orchestration
next-page cursor value extraction after repository results return
service-result construction
```

Placement follows responsibility and consumer, not the generic word “helper”.

---

## 8. Search Validation

Resource validators consume the public `core/search` façade.

The resource validator owns:

```text
allowed public fields
allowed operators per field
value normalization
sort policy
cursor value validation
resource-specific cursor ID validation
```

It returns:

```text
coresearch.Query[TId]
coresearch.CursorPayload
```

The generic search capability owns query mechanics. The repository owns persistence expression mapping.

Do not import `core/search/contracts` directly from resource validators.

---

## 9. Foundation Pattern

```text
foundation/
├── validators/
│   ├── person.go
│   ├── impl/person.go
│   └── helpers/person.go
├── services/
│   ├── person.go
│   ├── impl/person.go
│   └── helpers/person.go
└── mappers/
    ├── person.go
    └── impl/person.go
```

`PersonValidator` returns normalized request DTOs directly.

For search it returns:

```text
coresearch.Query[entity.GID]
coresearch.CursorPayload
```

without a `ValidatedPersonSearch` wrapper.

Repository-backed contact ownership/conflict checks remain in `PersonService`.

---

## 10. Identity Pattern

`AccountValidator` returns normalized `CreateAccountInput` / `UpdateAccountInput` directly.

It uses `core/normalize` for UUID/email/phone technical canonicalization and maps technical failures to Identity errors.

`AccountService` owns lifecycle, existence, uniqueness, and persistence orchestration.

---

## 11. Hard Rules

```text
NO payload validation duplicated inside service implementations
NO validator implementation mixed into validator interface files
NO concrete validator dependency outside the owning composition root
NO repository/network calls from validators
NO authorization/business workflow in validators
NO transport response mapping in validators
NO duplicate technical normalization when core/normalize owns it
NO direct imports of core capability internals
NO validator-only helpers left under services/helpers
NO validators/common or validators/utils dumping ground
NO ValidatedCreateX/ValidatedUpdateX when the existing type has the same shape
NO ValidatedSearchX when coresearch.Query already represents the result
YES resource validators map generic technical failures to domain errors
YES services receive normalized existing contracts
```

> **Guiding principle:** validate and canonicalize before business orchestration, reuse existing contracts, consume shared capabilities through their root façades, and create new validator result types only for genuinely new semantics.
