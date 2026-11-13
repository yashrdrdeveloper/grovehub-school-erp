# GroveHub Validator Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** request/service-input validation in all backend services and modules

Also read:

```text
docs/engineering/backend/CORE_LIBRARY_GUIDELINES.md
docs/engineering/backend/NORMALIZATION_AND_DRY_GUIDELINES.md
```

---

## REST request executor

The `libs/core/pipeline/rest.Run` API is a transport executor. A controller
provides typed `Config[Path, Query, Body, Response]`, then the pipeline:

1. Binds path and query parameters into separate typed DTOs (field tags).
2. Decodes a JSON body unless `Body` is `rest.Empty`.
3. Applies tagged normalization, resource canonicalization, DTO tag validation,
   and custom source validation to each typed input in that order.
4. Applies optional `Config.Validate` cross-source checks.
5. Calls the configured `Service(context.Context, Input[P,Q,B])` once.
6. Serializes one success/error response using the shared writer.

The service callback receives `Input{Path, Query, Body}`. Only the callback
calls the real business service; no business or authorization logic belongs
inside the pipeline.

```go
func (controller *PersonController) Create(w http.ResponseWriter, r *http.Request) {
    rest.Run(w, r, rest.Config[rest.Empty, rest.Empty, foundationrequest.CreatePerson, foundationresponse.Person]{
        Body: controller.validator.CreateRules(),
        InvalidBody: personerror.InvalidPersonRequest,
        Fallback: personerror.PersonCreateFailed,
        Service: func(ctx context.Context, input rest.Input[rest.Empty, rest.Empty, foundationrequest.CreatePerson]) (rest.Result[foundationresponse.Person], error) {
            person, created, err := controller.service.Create(ctx, input.Body)
            return rest.CreatedOrOK(person, created), err
        },
    })
}
```

Parameter binding supports exported scalar `string`, `bool`, signed/unsigned
integer and floating-point fields; optional pointers and repeated query
slices are supported. Use `path:"gid"` or `query:"page"` field tags and
`validate:"uuidv7"` / `validate:"required,min=1,max=100"` rules.
A `normalize:"trim,uuidv7"` tag canonicalizes an existing UUIDv7.
Missing path values, malformed scalar values and repeated scalar query keys
fail binding; unknown query keys are ignored by default, while 
`StrictQuery: true` rejects them. Malformed query encoding is rejected.
Use pointers when absent and
zero-value query fields must be distinguished.

`rest.Run` intentionally returns nothing (HTTP handlers return nothing).
Successful service callbacks return `rest.Result[R]`, normally via
`rest.OK`, `rest.Created` or `rest.CreatedOrOK`.
The default writer uses GroveHub's existing response envelope, optional
result headers and HTTP 204 without a body. `Config.Writer` can override
success/error presentation. `Config.InvalidPath`, `InvalidQuery`,
`InvalidBody`, and `Fallback` determine safe error responses; application
errors use their definition's HTTP status and metadata.
Always configure a safe `Fallback` and source-specific 4xx definitions.

Use `Decode: httprequest.DecodeStrictJSON` when rejecting unknown JSON
body properties. `rest.DecodeError` and `rest.BindError` do not expose
input values in responses. `rest.ReadBody` remains available for use cases
requiring only typed body decoding/validation. Services independently
validate public in-process calls with their reusable
`annotationvalidation.Rules[T]`, preserving current security boundaries.

## Request DTOs and header context

- Foundation uses a plain `dto/request.Person` with string fields embedded in `CreatePerson`; `UpdatePerson` uses shared presence-aware `annotation/dto.Patch[T]` fields, with `@PatchRequestDto` generating JSON serialization so absent, explicit-null and supplied values remain distinct.
- Validation and normalization are applied by existing validator code using the `normalize` and `validate` struct tags; repeating the same rules as `@NotBlank`, `@Size`, `@Valid`, `@Data`, or `@Validator` comments in these DTOs adds no value and is not required.
- `libs/core/headercontext.Middleware` captures and validates `Idempotency-Key`, `X-Correlation-ID` and `X-Request-ID` once at the launcher HTTP boundary. Services consume `headercontext.IdempotencyKey(ctx)` directly; no request-body header fields or controller extraction. Never treat client-supplied `X-Trace-ID` as authoritative telemetry.

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
coresearch.Query[string]
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
    -> canonical email/phone values (normalization)
    -> DTO tag validation of normalized values
    -> count/duplicate/primary checks (custom validation)
    -> normalized, validated request.Contacts
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
NormalizeContacts
ValidateContacts
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

`PersonValidator` exposes only three rule methods: `CreateRules()`, `UpdateRules()`, and `SearchRules()`. HTTP controllers pass these directly to `rest.Config.Body`; in-process Create and Update service methods use `annotationvalidation.Process(request, rules)`. Search services use the same resource-owned search helper to transform validated requests into `coresearch.Query[string]` and cursor payloads, while `SearchRules()` performs the check for REST callers.

For search it returns:

```text
coresearch.Query[string]
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
