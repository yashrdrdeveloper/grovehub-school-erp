# Identity Account Contract — Phase 2.1

## Scope

Phase 2.1 establishes the internal Identity Account foundation only. It does not expose an Account HTTP API.

Account lifecycle operations are internal Identity service capabilities. Registration/provisioning will create accounts later, while authentication flows will read accounts by email or account identity.

## Ownership

Foundation owns the canonical Person and GID.

Identity owns the login account used to authenticate that Person.

Identity stores the Foundation Person GID as a cross-domain reference but does not create a database foreign key to `foundation.persons`. Cross-service existence and orchestration must go through Foundation's public client.

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
- Supported lifecycle states in this phase are `ACTIVE`, `SUSPENDED`, and `DEACTIVATED`.
- PersonGid and AccountId are immutable after account creation.

## Internal lifecycle contract

`services.AccountService` owns the Identity account use cases and is implemented under `services/impl`.

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

`Update` currently changes email and/or phone only. Email can never be cleared. Phone may be cleared.

`Deactivate` is a lifecycle state change and does not delete the row.

`PermanentlyDelete` irreversibly removes the Identity account row and is intentionally distinct from deactivation. Any future registration/admin flow exposing permanent deletion must apply its own authorization and safety policy before calling this operation.

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

Business lifecycle decisions such as deactivation remain in the service layer and are persisted through `Update`. The repository uses the generic persistence verb `Delete`; the service exposes the more explicit `PermanentlyDelete` use case.

## Registration boundary

There is deliberately no `/identity/accounts` endpoint.

The later registration/provisioning use case will perform this flow:

```text
Registration / Provisioning
        ↓
Resolve or create Foundation Person
        ↓
Foundation GID
        ↓
Identity AccountService.Create(...)
        ↓
identity.accounts
```

The Account service is therefore an internal Identity use-case dependency, not an HTTP resource.

## Authentication readiness

Phase 2.1 does not implement authentication.

Later Google and Microsoft login can use the provider-verified email to locate the existing GrooveHub Identity Account through `AccountService.GetByEmail`. No GrooveHub password is required for those provider login flows.

Password credentials, Google/Microsoft protocol handling, sessions, access tokens, refresh tokens, logout, OTP, MFA, and launcher authentication middleware are explicitly outside Phase 2.1.
