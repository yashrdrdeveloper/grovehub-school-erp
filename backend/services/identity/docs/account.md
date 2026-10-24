# Identity Account Contract — Phase 2.1

## Scope

Phase 2.1 establishes the internal Identity Account foundation. Account lifecycle operations remain internal service capabilities; there is no public `/identity/accounts` HTTP API.

Registration now consumes `AccountService.Create` as part of `POST /auth/register`. Authentication flows will later read accounts by email or account identity.

## Ownership

Foundation owns the canonical Person and GID.

Identity owns the login account used to authenticate that Person.

Identity stores the Foundation Person GID as a cross-domain reference but does not create a database foreign key to `foundation.persons`. Cross-service existence and orchestration go through Foundation's public client.

## Account model

```text
Account
├── AccountId
├── PersonGid
├── Email        (required)
├── Phone        (optional)
├── Status
├── Version
├── CreatedAt
└── UpdatedAt
```

Rules:

- One Identity Account per Foundation Person GID.
- Email is mandatory for every account.
- Email is normalized to lowercase, trimmed, and globally unique.
- Phone is optional.
- When present, phone is normalized to E.164 form and globally unique.
- Account IDs use UUIDv7.
- Initial account status is `ACTIVE`.
- Supported lifecycle states are `ACTIVE`, `SUSPENDED`, and `DEACTIVATED`.
- PersonGid and AccountId are immutable after account creation.

## Validation boundary

```text
services.CreateAccountInput / UpdateAccountInput / lookup value
        ↓
validators.AccountValidator
        ↓
validators/impl.AccountValidator
        ↓
validated / normalized existing contract
        ↓
services.AccountService
```

The validator owns AccountId/PersonGid UUIDv7 validation, required email normalization/validation, optional phone normalization/validation, and create/update input validation. It reuses `libs/core/helpers/normalize` and performs no repository I/O or lifecycle decisions.

Repository-backed conflicts, reads, lifecycle transitions, persistence, and error mapping remain in `AccountService`.

## Internal lifecycle contract

Supported operations:

```text
Create
GetById
GetByPersonGid
GetByEmail
Update
Deactivate
PermanentlyDelete
```

`Update` changes email and/or phone only. Email cannot be cleared. Phone may be cleared.

`Deactivate` changes lifecycle state without deleting the row.

`PermanentlyDelete` irreversibly removes the Account and is intentionally internal. Registration currently uses it only as compensation when a newly created Account cannot receive its password credential.

## Repository boundary

`repository.AccountRepository` exposes persistence primitives only:

```text
Create
FindById
FindByPersonGid
FindByEmail
Update
Delete
```

Business lifecycle decisions remain in the service layer.

## Registration boundary

Registration performs:

```text
Resolve/create Foundation Person through Foundation Client
        ↓
Foundation GID + normalized primary contacts
        ↓
AccountService.Create
        ↓
PasswordCredentialService.Set
```

The Account service remains an internal Identity dependency rather than an HTTP resource.

See `backend/services/identity/docs/registration.md` for orchestration and compensation rules.

## Authentication readiness

Google and Microsoft login can later use provider-verified email to locate an existing GrooveHub Identity Account through `AccountService.GetByEmail`. A GrooveHub password is not required for provider-only accounts.

Login, sessions, access tokens, refresh tokens, logout, OTP, MFA, provider protocol handling, and launcher authentication middleware remain future work.
