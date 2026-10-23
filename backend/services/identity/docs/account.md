# Identity Account Contract — Phase 2.1

## Scope

Phase 2.1 establishes the internal Identity Account foundation only. It does not expose an Account HTTP API.

Account creation will be invoked later by the Identity registration/provisioning flow after the Foundation Person has been resolved or created.

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

Later Google and Microsoft login can use the provider-verified email to locate the existing GrooveHub Identity Account. No GrooveHub password is required for those provider login flows.

Password credentials, Google/Microsoft protocol handling, sessions, access tokens, refresh tokens, logout, OTP, MFA, and launcher authentication middleware are explicitly outside Phase 2.1.
