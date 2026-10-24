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
Validated / normalized values
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
    -> validated result contracts

validators/impl/person.go
    -> concrete validation orchestration

validators/helpers/person.go
    -> pure resource-specific validation helpers
```

Concrete validator implementations are constructed only by the owning service composition root.

## 3. Validator responsibilities

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
creation of validated/normalized input value objects
```

Validators may receive already-loaded entity values when validation needs existing state but does not require I/O. Example: preserving existing contact verification flags while validating an updated contact payload.

## 4. Validator forbidden responsibilities

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

## 5. Validator vs mapper

Do not mix these responsibilities.

```text
Mapper
    -> representation conversion

Validator
    -> payload validation + canonicalization

Service
    -> business orchestration
```

Examples:

```text
PersonByGID DTO -> entity.GID primitive conversion
    -> mapper may perform the representation conversion

validate that GID is UUIDv7
    -> validator

load Person by GID and decide not-found behavior
    -> service
```

A mapper must not become a validation layer, and a validator must not produce transport responses.

## 6. Shared normalization

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

## 7. Validator helpers

A helper belongs under `validators/helpers` when its consumers are validator implementations and its responsibility is validation/normalization.

Examples:

```text
NormalizePersonNames
BuildContacts
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
entity comparison used by business update orchestration
next-page cursor value extraction after repository results return
service-result construction helpers
```

Placement follows the actual consumer and responsibility.

## 8. Service methods after validation

Service methods should read primarily as business steps.

Preferred shape:

```text
validate payload
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

## 9. Foundation pattern

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

`PersonValidator` validates and normalizes create/update/search/GID inputs.

`PersonService` retains repository-backed rules such as contact ownership/conflict checks.

## 10. Identity pattern

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

`AccountValidator` owns account ID, Person GID, email, phone, create-input, and update-input validation.

`AccountService` owns lifecycle and persistence orchestration.

## 11. Hard rules

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
```

> **Guiding principle:** validate and canonicalize before business orchestration; keep I/O and business decisions in the service.
