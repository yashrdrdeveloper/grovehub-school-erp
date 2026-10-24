# Identity Registration Contract

## Scope

Registration is the first public Identity orchestration flow. It creates or resolves the canonical Foundation Person, creates the internal Identity Account, and creates the optional password credential used by password login.

Public endpoint:

```text
POST /auth/register
```

There is still no public `/identity/accounts` endpoint and no standalone password endpoint.

## Request contract

Identity reuses Foundation's `request.CreatePerson` contract directly and adds only the registration-specific password field.

```go
type Register struct {
    foundationrequest.CreatePerson
    Password string `json:"password"`
}
```

Because `CreatePerson` is embedded, the HTTP JSON remains flat. Example shape:

```json
{
  "firstName": "Asha",
  "lastName": "Sharma",
  "contacts": {
    "emails": [
      {"value": "asha@example.com", "isPrimary": true}
    ],
    "phones": [
      {"value": "+919999999999", "isPrimary": true}
    ]
  },
  "password": "example-password"
}
```

Registration does not duplicate Foundation's name/contact DTOs.

## Registration validator

`RegistrationValidator` owns only registration-specific structural rules:

```text
new-password policy
exactly one primary email
at most one primary phone
```

It does not duplicate Foundation name/email/phone format validation. Foundation validates and normalizes its own `CreatePerson` payload when Registration calls the Foundation client.

The primary email is mandatory because every Identity Account requires an email. The primary phone is optional.

## Orchestration

```text
POST /auth/register
        ↓
RegistrationController
        ↓
RegistrationService
        ↓
RegistrationValidator
        ↓
Foundation Client.CreatePerson
        ↓
normalized Foundation Person + GID
        ↓
AccountService.Create
        ↓
PasswordCredentialService.Set
        ↓
Registration response
```

The Account email/phone values come from the normalized Foundation Person response, not directly from the raw HTTP request. Registration uses the primary Foundation email and, when present, the primary Foundation phone.

## Cross-service boundary

Registration depends on Foundation only through `foundation/client.Client`.

Identity must not import Foundation repositories, service implementations, or database tables. The launcher constructs Foundation first and injects `foundationModule.Client()` into the Identity module.

```text
Launcher
  ├── Foundation Module
  │      └── public Foundation Client
  │
  └── Identity Module
         └── RegistrationService
                └── Foundation Client
```

## Failure and compensation

Foundation Person and Identity Account are separate service responsibilities, so Registration does not attempt a cross-service database transaction.

If Foundation Person creation succeeds but Account creation fails, Registration keeps the Person. Foundation Person creation is contact-aware/idempotent, so a later registration retry can resolve the same Person rather than creating another one.

If Account creation succeeds but password credential creation fails, Registration compensates by calling `AccountService.PermanentlyDelete` for the newly created Account. The password credential foreign key cascades on Account deletion if a credential row exists.

Registration never deactivates or hard-deletes the Foundation Person as compensation.

If compensation itself fails, Registration returns the generic `IDENTITY_REGISTRATION_FAILED` error with both causes retained internally.

## Error boundary

Registration translates dependency errors into registration-level public errors:

```text
IDENTITY_REGISTRATION_INVALID        -> 400
IDENTITY_REGISTRATION_ALREADY_EXISTS -> 409
IDENTITY_REGISTRATION_FAILED         -> 500
```

Downstream Foundation/Account/Password persistence details are not exposed as the public registration contract.

## Response

The current response contains Identity identifiers/state only:

```text
accountId
personGid
email
phone (optional)
status
```

Registration does not create a session or tokens yet. Login/session/token work remains a later authentication phase.
