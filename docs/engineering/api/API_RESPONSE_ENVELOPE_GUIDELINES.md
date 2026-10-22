# API Response Envelope Guidelines

> **Status:** Mandatory engineering standard  
> **Applies to:** GroveHub REST/JSON APIs across all backend services.  
> **Relationship:** Additive to `API_CONTRACT_GUIDELINES.md` and `BACKEND_CODE_GUIDELINES.md`.

## 1. Organization-wide response contract

Successful JSON responses use:

```json
{
  "success": true,
  "data": {}
}
```

Error JSON responses use:

```json
{
  "success": false,
  "errorCode": "IDENTITY_SERVICE_DOWN",
  "errorMessage": "Identity service is unavailable",
  "httpErrorCode": 503,
  "errorMetaData": {}
}
```

The outer envelope is shared infrastructure. The `data` payload and semantic error definitions remain owned by the service contract.

## 2. Shared backend placement

Shared response DTOs live under:

```text
backend/libs/core/dto/response/
```

Shared HTTP response writing lives under:

```text
backend/libs/core/http/response/
```

Shared error primitives live under:

```text
backend/libs/core/errors/
```

Do not create a generic `common/`, `misc/`, or giant `utils/` package for response handling.

A future shared request DTO package may live under:

```text
backend/libs/core/dto/request/
```

only when a real organization-wide request primitive exists.

## 3. Generic response DTOs

The shared response DTOs use Go generics:

```go
type Success[T any] struct {
    Success bool `json:"success"`
    Data    T    `json:"data"`
}

type Error[M any] struct {
    Success       bool            `json:"success"`
    ErrorCode     coreerrors.Code `json:"errorCode"`
    ErrorMessage  string          `json:"errorMessage"`
    HTTPErrorCode int             `json:"httpErrorCode"`
    ErrorMetaData M               `json:"errorMetaData"`
}
```

This keeps payload and error metadata strongly typed without duplicating envelopes per service.

## 4. Semantic error definitions

A semantic error definition contains only stable service-owned semantics:

```go
type Definition struct {
    // private fields
}
```

It exposes:

```text
Code()
Message()
```

The definition does not own an HTTP status code and does not own request-specific metadata.

Reason:

```text
semantic failure
    !=
transport mapping
```

The same semantic error may eventually be represented through HTTP, gRPC, events, or an in-process client.

## 5. Service-owned error definitions

Error definitions belong to the service that owns the failure semantics.

Group them by responsibility instead of one giant file.

Example:

```text
backend/services/identity/constants/errorcode/
├── health_code.go
├── authentication_code.go
└── session_code.go
```

For health errors:

```go
package errorcode

import coreerrors "grovehub/libs/core/errors"

var (
    ServiceDown = coreerrors.NewDefinition(
        "IDENTITY_SERVICE_DOWN",
        "Identity service is unavailable",
    )
)
```

Go does not allow a struct value to be declared with `const`, so semantic definitions are package variables created from immutable-style core definitions. Their fields are private and consumers can only read `Code()` and `Message()`.

Do not place service-specific error definitions in `libs/core`.

Use stable uppercase identifiers with the owning service/domain prefix:

```text
IDENTITY_SERVICE_DOWN
IDENTITY_INVALID_CREDENTIALS
COMMUNICATION_PROVIDER_UNAVAILABLE
ADMISSIONS_APPLICATION_NOT_FOUND
```

Published error codes are API contract values and must not be renamed casually.

## 6. HTTP status ownership

HTTP status belongs to the controller/transport boundary, not the service error definition.

Correct:

```go
httpresponse.WriteError(
    writer,
    http.StatusServiceUnavailable,
    errorcode.ServiceDown,
)
```

This produces:

```json
{
  "success": false,
  "errorCode": "IDENTITY_SERVICE_DOWN",
  "errorMessage": "Identity service is unavailable",
  "httpErrorCode": 503,
  "errorMetaData": {}
}
```

`httpErrorCode` is populated from the actual HTTP status passed to the shared writer. Do not duplicate the status inside service constants.

## 7. Error metadata

Most errors do not need metadata.

Use:

```go
httpresponse.WriteError(
    writer,
    http.StatusServiceUnavailable,
    errorcode.ServiceDown,
)
```

When safe structured context is useful, use:

```go
httpresponse.WriteErrorWithMetadata(
    writer,
    http.StatusBadRequest,
    errorcode.SomeValidationError,
    validationMetadata,
)
```

Metadata must be typed, caller-safe, and secret-free. Never return SQL details, stack traces, tokens, passwords, provider secrets, or sensitive configuration.

## 8. Success payload ownership

The controller must not invent public response DTO structs inline when the payload is part of the API contract.

Service response DTOs live under the owning service boundary, for example:

```text
backend/services/identity/dto/v1/response/health.go
```

Controller usage:

```go
httpresponse.WriteSuccess(
    writer,
    http.StatusOK,
    identityresponse.HealthData{
        Service: "identity",
        Status:  "ok",
    },
)
```

Controllers coordinate transport. They should not become DTO definition packages.

## 9. Controller rule

Controllers:

- select the HTTP status code
- map service results into versioned response DTOs
- select the appropriate service-owned semantic error definition
- append safe metadata only when required
- use the shared response writer

Controllers must not duplicate `json.NewEncoder`, envelope fields, error messages, or semantic error codes per endpoint.

Business services must not depend on `net/http`.

## 10. Contract rule

OpenAPI remains the source of truth for externally visible HTTP payloads.

Whenever an endpoint changes its response envelope, DTO, fields, or error contract, update the OpenAPI schema and examples in the same change.

## 11. Exceptions

Do not wrap responses where transport semantics require another shape, including:

- file/download streaming
- SSE/WebSocket streams
- provider callbacks with mandated payloads
- `204 No Content`

Any exception must be explicit in the endpoint contract.
