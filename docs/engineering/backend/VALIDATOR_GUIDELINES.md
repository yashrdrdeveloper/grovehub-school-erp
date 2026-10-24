# GroveHub Validator Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** request/service-input validation in all backend services and modules

## 1. Purpose

Validators keep request-payload validation and normalization out of business services.

The standard flow is:

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

A validator answers: **is this payload structurally and semantically valid before business orchestration, and what are its canonical values?**

A service answers: **what business action should happen with valid values?**

## 2. Canonical structure

```text
validators/
├── <resource>.go
├── impl/
│   └── <resource>.go
└── helpers/
    └── <resource>.go
```

Use the same contract/implementation convention as services, mappers, repositories, and controllers.

```text
validators/person.go
    -> PersonValidator interface
    -> reuses existing DTO/service-input/core contracts

validators/impl/person.go
    -> concrete validation orchestration

validators/helpers/person.go
    -> pure resource-specific validation helpers
```

Concrete validator implementations are constructed only by the owning service composition root.

## 3. Reuse existing contracts

Validation does not create a new data model by itself.

If validation only trims, normalizes, canonicalizes, or validates values while the shape remains the same, return the existing type.

Examples:

```text
CreatePerson DTO
    -> validate/normalize
    -> CreatePerson DTO

UpdatePerson DTO
    -> validate/normalize
    -> UpdatePerson DTO

CreateAccountInput
    -> validate/normalize
    -> CreateAccountInput

UpdateAccountInput
    -> validate/normalize
    -> UpdateAccountInput
```

For shared generic models, return the generic model directly:

```text
search.Request
    -> validate/normalize
    -> coresearch.Query[entity.GID]
    -> coresearch.CursorPayload
```

Do not create wrappers such as:

```text
ValidatedCreatePerson
ValidatedUpdatePerson
ValidatedCreateAccount
ValidatedUpdateAccount
ValidatedPersonSearch
```

when those types only duplicate an existing contract and add no new semantics.

A validator-specific result type is allowed only when validation produces genuinely new information that cannot be represented safely by an existing DTO, service input, entity value object, or generic core contract.

## 4. Validator responsibilities

Validators may own:

```text
required/optional field validation
string length validation
request-local relationship checks
normalization of request values
UUID/email/phone format validation through libs/core helpers
contact-list shape validation
search request validation
filter/operator/value-count validation
cursor payload validation
canonicalization of existing request/service-input contracts
```

Validators should validate request payloads, not construct a parallel validated model just to transfer values to the service.

## 5. Validator forbidden responsibilities

Validators must not:

```text
query repositories
call external services
perform authorization
choose business workflows
make lifecycle transitions
generate IDs
read clocks
write entities to persistence
publish events
select HTTP status codes
```

If validation requires a repository query, that check belongs in the service because it is a business/integration rule, not payload validation.

Examples that stay in the service:

```text
whether an email already belongs to another account
whether a Foundation contact belongs to another Person
whether the caller is allowed to deactivate an account
whether an entity exists
whether a lifecycle transition is allowed
```

## 6. Validator vs mapper vs service helper

Do not mix these responsibilities.

```text
Mapper
    -> representation conversion

Validator
    -> payload validation + canonicalization

Service helper
    -> service-owned entity/value construction from already valid data

Service
    -> business orchestration
```

Example:

```text
CreatePerson.Contacts
    ↓
Validator
    -> checks counts, duplicates, primary rules
    -> normalizes email/phone values
    -> returns normalized request.Contacts
    ↓
Service helper
    -> constructs entity.Contact values
    -> preserves existing verification state for update
```

Repository-backed checks remain in the service.

A mapper must not become a validation layer, and a validator must not produce transport responses.

## 7. Shared normalization

Validators must reuse stable technical normalizers from:

```text
backend/libs/core/helpers/normalize/
```

Examples:

```text
normalize.Email
normalize.Phone
normalize.PhoneCharacters
normalize.UUIDv7
```

Do not duplicate email, phone, or UUID parsing inside service validators.

The shared normalizer returns a generic technical failure. The owning validator maps that failure to the service/domain error contract.

## 8. Validator helpers

A helper belongs under `validators/helpers` when its consumers are validator implementations and its responsibility is validation/normalization.

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

Likewise, do not move genuine service helpers into validators merely to reduce the size of the service package.

Examples that may remain under `services/helpers`:

```text
entity/value construction from validated DTOs
preserving existing entity state while rebuilding values
entity comparison used by business update orchestration
next-page cursor value extraction after repository results return
service-result construction helpers
```

Placement follows the actual consumer and responsibility.

## 9. Service methods after validation

Service methods should read primarily as business steps.

Preferred shape:

```text
validate payload
    ↓
construct domain values if required
    ↓
load/check related state
    ↓
apply business rule
    ↓
construct/update entity
    ↓
persist
    ↓
map response
```

Avoid embedding repeated field parsing, trimming, regex validation, or search-filter validation directly inside service methods.

## 10. Foundation pattern

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

`PersonValidator` returns normalized `foundationrequest.CreatePerson` / `foundationrequest.UpdatePerson` values directly.

For search it returns:

```text
coresearch.Query[entity.GID]
coresearch.CursorPayload
```

directly, without a `ValidatedPersonSearch` wrapper.

`PersonService` retains repository-backed rules such as contact ownership/conflict checks and uses service helpers for entity contact construction from already validated contact DTOs.

## 11. Identity pattern

```text
identity/
├── validators/
│   ├── account.go
│   ├── impl/account.go
│   └── helpers/account.go
└── services/
    ├── account.go
    └── impl/account.go
```

`AccountValidator` returns normalized `CreateAccountInput` / `UpdateAccountInput` directly instead of duplicate validated wrapper types.

`AccountService` owns lifecycle and persistence orchestration.

## 12. Hard rules

```text
NO request-payload validation duplicated inside service implementations
NO validator implementation mixed into validator contract files
NO concrete validator dependency outside the owning composition root
NO repository/network calls from validators
NO authorization/business workflow in validators
NO transport response mapping in validators
NO duplicated technical normalization when libs/core already owns it
NO validator-only helpers left under services/helpers
NO generic validators/common or validators/utils dumping ground
NO ValidatedCreateX/ValidatedUpdateX type when the existing request/service-input type has the same shape
NO ValidatedSearchX wrapper when existing generic search types already represent the result
```

> **Guiding principle:** validate and canonicalize before business orchestration, reuse existing contracts whenever possible, and create new validator result types only when they carry genuinely new semantics.
