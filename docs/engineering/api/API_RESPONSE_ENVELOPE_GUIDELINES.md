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

The outer envelope is shared infrastructure. The `data` payload remains owned by the service contract.

## 2. Shared backend placement

Shared response DTOs live under:

```text
backend/libs/core/dto/response/
```

Shared HTTP response writing lives under:

```text
backend/libs/core/http/response/
```

Shared error-code typing lives under:

```text
backend/libs/core/errors/
```

Do not create a generic `common/`, `misc/`, or giant `utils/` package for response handling.

A future shared request DTO package may live under:

```text
backend/libs/core/dto/request/
```

but it must be created only when a real organization-wide request primitive exists. Services must not be forced into a generic request envelope without a concrete requirement.

## 3. Generics

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

## 4. Service-owned error codes

Error-code values belong to the service that owns the failure semantics.

Example:

```text
backend/services/identity/constants/errorcode/
```

```go
const ServiceDown coreerrors.Code = "IDENTITY_SERVICE_DOWN"
```

Do not place service-specific error codes in `libs/core`.

Use uppercase stable identifiers with the owning service/domain prefix.

Examples:

```text
IDENTITY_SERVICE_DOWN
IDENTITY_INVALID_CREDENTIALS
COMMUNICATION_PROVIDER_UNAVAILABLE
ADMISSIONS_APPLICATION_NOT_FOUND
```

Do not rename an externally published error code casually; treat it as part of the API contract.

## 5. Error fields

`success`
: Always `false` for the error envelope.

`errorCode`
: Stable machine-readable service-owned code.

`errorMessage`
: Safe human-readable message. Never expose SQL, stack traces, secrets, tokens, provider internals, or sensitive configuration.

`httpErrorCode`
: Must equal the actual HTTP response status code.

`errorMetaData`
: Typed per call through generics. It may contain structured diagnostic or validation context that is safe for the caller. It must not contain secrets.

## 6. Controller rule

Controllers should use the shared HTTP writer rather than manually rebuilding the envelope.

```go
httpresponse.WriteSuccess(writer, http.StatusOK, data)

httpresponse.WriteError(
    writer,
    http.StatusServiceUnavailable,
    errorcode.ServiceDown,
    "Identity service is unavailable",
    metadata,
)
```

Controllers still own transport mapping; business services should not depend on `net/http`.

## 7. Contract rule

OpenAPI remains the source of truth for externally visible HTTP payloads.

Whenever an endpoint changes to or adds one of these envelopes, update its OpenAPI schema and examples in the same change.

## 8. Exceptions

Do not wrap responses where the transport semantics require a different body shape, for example:

- file/download streaming
- SSE/WebSocket streams
- protocol-specific callbacks where an external provider mandates a payload
- `204 No Content`

Any exception must be explicit in the endpoint contract.
