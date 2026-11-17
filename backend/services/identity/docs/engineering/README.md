# Identity Service — Engineering Guide

**Source status:** reviewed against `phase-2.1-identity-account` source for registration, account and password-credential interfaces. This guide documents observed behavior and clearly separates architecture requirements from work not yet verified.

## Shared standards

Follow the generic documents under `docs/engineering/backend/`: `backend.md`, `searchDetails.md`, and `annotation.md`. This file is limited to Identity-specific ownership, actual components, integration boundaries and decisions. Do not duplicate global DRY, SRP, ACID, API or annotation rules here.

## Module boundary

Identity owns **accounts** and **authentication credentials**. Foundation owns person records and contacts. Identity may orchestrate Foundation via its typed client; it must not directly own Foundation persistence. Authentication tokens, OIDC, OAuth2, SAML, sessions, MFA, authorization, RBAC and password-reset capabilities must not be presented as shipped until the respective implementations and contracts are verified.

## Observed directory ownership

| Observed path | Responsibility |
|---|---|
| `controller/auth.go` | Declares `AuthController` and its `Register` HTTP handler contract; annotates `@Routes /auth` and `@POST /register` |
| `controller/impl/auth.go` | Adapts HTTP to `RegisterRequest`, executes core REST pipeline and maps creation to `201` |
| `dto/request/auth.go` | Register request body and request wrapper; structural trim/validation annotations |
| `dto/response/auth.go` | Register output fields |
| `validators/auth.go` | Exposes reusable registration validation rules |
| `services/auth.go` | Application-facing registration interface |
| `services/impl/auth.go` | Orchestrates registration, idempotency and Foundation/account/credential operations |
| `services/account.go` | Account lifecycle service contract |
| `services/impl/account.go` | Implements account lifecycle and persistence error mapping |
| `services/impl/password_credential.go` | Creates, verifies and replaces password hashes |
| `entity/account.go` | Account persistence shape, status and version |
| `repository/account.go` | Account persistence interface and sentinel errors |

This is an inspected file inventory, not an exhaustive package tree. Add remaining files only after verification from source.

## Registration use case (current code path)

```text
HTTP handler
  -> strict JSON decode + core REST validation
  -> AuthService.Register
     -> read Idempotency-Key from request context and validate
     -> attach causation ID
     -> validate/normalize request and fingerprint body
     -> start Identity GORM transaction
        -> find existing account by registration idempotency key
          -> same fingerprint: reconcile password credential and return account
          -> different fingerprint: conflict
        -> Foundation SearchPersons by email
           -> none: Foundation CreatePerson
           -> one: reuse person
           -> multiple: failure
        -> check account by person GID
        -> Foundation UpdatePerson
        -> AccountService.Create
        -> PasswordCredentialService.Set
     -> commit or return error
  -> response mapping
```

### Ownership and invariants

- Idempotency key comes from request context, not the JSON body; the normalized request body is fingerprinted.
- The same key with different normalized payload is rejected.
- Account references a Foundation-owned person GID; it stores email, optional phone, lifecycle status and version.
- Account statuses currently include `ACTIVE`, `SUSPENDED`, and `DEACTIVATED`.
- Password verification and replacement are encapsulated in the password-credential service.
- Registration needs account and credential persistence to be consistent within Identity.

### Critical transaction boundary

`coretransaction.Run(ctx, service.db, ...)` can cover **Identity operations using the same transaction-aware database context**. It cannot guarantee rollback for Foundation operations performed through an independent client / remote service. The registration path currently invokes Foundation create/update inside the local transaction callback.

**Required follow-up:** confirm whether Foundation is in-process with a shared transactional database context or is a separate transactional boundary. If separate, design idempotent cross-boundary coordination, recovery and/or saga compensation. Do not claim end-to-end ACID or automatic Foundation rollback without tests proving it.

### Additional risks to verify before release

1. Confirm transaction propagation from `coretransaction.Run` through account and credential repositories.
2. Confirm retry of a successful registration does not unexpectedly change password credentials, including when reusing the same key with a different password.
3. Test concurrent registration by email, person GID and idempotency key; rely on DB uniqueness constraints.
4. Verify two-step Foundation search/create under concurrency and how duplicate-email races are resolved.
5. Confirm password hash algorithm and parameter management, safe comparisons and secret-redacted logging.
6. Confirm status checks and account restrictions before authentication succeeds.
7. Ensure conflicting identities across email/phone cannot allow account takeover.
8. Confirm actual public route prefix and registration auth/rate limiting in runtime router.
9. Verify `@Routes`, `@Entity`, and validation annotations have active generation/processing support.
10. Document all actual migrations, indexes, tenant scoping and audit events after reading their implementation.

## Dependency rules

HTTP handlers depend on application service interfaces and validator contracts; they do not call repositories or Foundation directly. Orchestration belongs in services; persistence in repository implementations; deterministic mapping in mappers; transport shape in DTOs. Foundation access goes through the typed client.

## Engineering documentation maintenance

Update this guide when a component, invariant, transaction boundary, event contract, persistence assumption or integration changes. Record concrete test commands and configuration only after they are checked against this branch. User journeys and reproducible API calls belong in `../product/`.
