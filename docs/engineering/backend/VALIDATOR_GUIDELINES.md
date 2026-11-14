# GroveHub Validator Guidelines

> **Status:** Mandatory GroveHub backend standard  
> **Applies to:** request/service-input validation in all backend services and modules

Also read:

```text
docs/engineering/backend/CORE_LIBRARY_GUIDELINES.md
docs/engineering/backend/NORMALIZATION_AND_DRY_GUIDELINES.md
```

---

## Common Request and REST pipeline

The shared `libs/core/request.Request` is embedded in every endpoint
contract. Concrete DTOs declare only the required **typed** sections:
`Params`, `Query`, and/or `RequestBody`; omitted fields do not need
placeholder or `Empty` types.

```go
type UpdatePersonRequest struct {
    corerequest.Request
    Params PersonByGID `validate:"valid"`
    RequestBody UpdatePerson `validate:"valid"`
}
```

`libs/core/pipeline/rest.RunStruct` receives the DTO via
`StructConfig[Request, Response]`. The shared inspector discovers sections;
the pipeline binds present Params and Query fields and decodes RequestBody
using JSON, then runs `annotationvalidation.Process` and custom
resource validation before calling `Service(ctx, input)` exactly once.
The callback forwards `input` unchanged to the business service.
Services independently reuse their validators on public entry points.

```go
func (controller *PersonController) Update(w http.ResponseWriter, r *http.Request) {
    rest.RunStruct(w, r, rest.StructConfig[foundationrequest.UpdatePersonRequest, foundationresponse.Person]{
        Rules: controller.validator.UpdateRules(),
        InvalidParams: personerror.InvalidPersonGID,
        InvalidRequestBody: personerror.InvalidPersonRequest,
        Fallback: personerror.PersonUpdateFailed,
        Service: func(ctx context.Context, input foundationrequest.UpdatePersonRequest) (rest.Result[foundationresponse.Person], error) {
            person, err := controller.service.Update(ctx, input)
            return rest.OK(person), err
        },
    })
}
```

A `@RequestLine` generated client takes the same
`foundationrequest.UpdatePersonRequest` type and uses the common inspector
to map its `Params` to `{gid}`, `Query` to URL query fields, and
`RequestBody` to JSON. Missing sections are left out automatically.
`Params` and `Query` structs use the `path:"..."` and
`query:"..."` tags. Scalar numbers, strings, booleans, optional pointers
and repeated query slices are supported. Invalid or duplicated scalar
values fail binding. `StructConfig.StrictQuery` rejects unknown query
values even if the request DTO has no Query section.

Normalization precedes validation, with tags such as
`normalize:"trim,uuidv7"`, `validate:"uuidv7"`,
`validate:"notblank,maxrunes=100"`, `validate:"min=1,max=100"`,
and nested `validate:"valid"`. PATCH fields use
`annotationdto.Patch[T]` to preserve absent/null/value state.
`RunStruct` returns no value; `rest.OK`, `rest.Created`, and
`rest.CreatedOrOK` describe the HTTP result. The writer preserves
semantic status, safe error metadata, headers and HTTP 204.

Configure `InvalidParams`, `InvalidQuery`, `InvalidRequestBody`
and `Fallback`. Application errors retain their own semantic status.
Use `Decode: httprequest.DecodeStrictJSON` when rejecting unknown JSON
body keys. Validation rules do not implement authorization; controllers,
services and IAM still enforce their own responsibilities.


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
└── person/
    ├── person.go                  # validator contract
    ├── impl/person.go             # composed normalization + validation rules
    └── helpers/
        ├── requests.go            # Create/Update/Search rule callbacks
        ├── contacts.go            # contact constraints and canonicalization
        ├── search.go              # query normalization + filter/sort/cursor flow
        ├── search_fields.go       # searchable field/operator policies
        ├── cursor.go              # sort cursor value validation
        ├── gid.go                 # GID validation
        └── errors.go              # domain error mapping
```

`validators/person/person.go` defines the interface.

`validators/person/impl/person.go` composes typed validation rules.

`validators/person/helpers/` groups reusable, resource-specific validation and normalization policies by responsibility.

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
├── validators/person/
│   ├── person.go
│   ├── impl/person.go
│   └── helpers/{requests,contacts,search,search_fields,cursor,gid,errors}.go
├── services/person/
│   ├── person.go
│   ├── impl/person.go
│   └── operations/{create,update,search,getbygid,deactivate}/
└── mappers/person/
    ├── person.go
    └── impl/
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

## Person validator composition and ownership

The Person validator composes `annotationvalidation.Rules[T]` using method
references from `validators/person/helpers/requests.go`; no inline closures
are needed for Create, Update or Search. The core library still determines the
order: DTO normalization -> custom normalization -> DTO tag validation -> custom
validation. Create and Update reuse the same contact canonicalization and
cross-field checks, while Update guards absent and explicit-null PATCH contacts.

`helpers/errors.go` owns `MapPersonDTOError`, which maps tagged Contact
failures to `InvalidPersonContact` and other tagged failures to
`InvalidPersonRequest`. Search field/operator configuration remains strictly
Person-specific and is split between `search_fields.go` and `cursor.go`.
Search filtering and pagination retain their existing external error codes.
