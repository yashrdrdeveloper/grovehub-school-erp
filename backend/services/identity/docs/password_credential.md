# Identity Password Credential Contract

## Scope

Password credentials are internal Identity state used by higher-level authentication flows. They are not an HTTP resource and do not have standalone password CRUD endpoints.

There must be no routes such as:

```text
POST /password
PUT /password
DELETE /password
```

Higher-level flows call the internal password credential service instead:

```text
Register              -> PasswordCredentialService.Set
Login                 -> PasswordCredentialService.Verify
Forgot/Reset Password -> PasswordCredentialService.Replace
Change Password       -> verify current password in the flow, then Replace
```

Registration is now implemented through `POST /auth/register`; the other authentication flows remain future work.

The higher-level flow owns authorization, OTP/recovery proof, current-password checks, account status checks, throttling, and response behavior.

## Account relationship

An Identity Account may exist without a password credential.

```text
Account
├── PasswordCredential   optional
├── Google identity      future
├── Microsoft identity   future
└── Passkey              future
```

This allows provider-only or passkey-only accounts without requiring a GrooveHub password.

`identity.password_credentials.account_id` is both the primary key and a foreign key to `identity.accounts.account_id`. The relationship is one-to-zero-or-one. Deleting an Identity Account cascades to its password credential.

## Model

```text
PasswordCredential
├── AccountId
├── PasswordHash
├── PasswordChangedAt
├── Version
├── CreatedAt
└── UpdatedAt
```

The credential entity is internal. Service methods do not return it to callers, avoiding password-hash exposure outside credential orchestration.

## Internal service contract

```text
Set
Verify
Replace
```

`Set` creates the first password credential for an existing account and is currently used by Registration.

`Verify` compares a candidate password with the stored encoded hash and returns a boolean. It does not decide whether login is allowed.

`Replace` stores a new password hash and increments credential version. It does not itself prove authorization for reset/change-password.

There is no standalone Remove/Delete password use case. Account deletion cascades automatically.

## Validator boundary

For newly set/replaced passwords, the current policy is:

```text
valid UTF-8
minimum 8 Unicode code points
maximum 128 Unicode code points
```

Passwords are never trimmed, lowercased, or normalized. The exact submitted value is hashed.

Verification uses candidate-password validation rather than the current minimum-length policy so future policy increases do not invalidate older credentials.

## Password hashing helper

Password hashing remains an internal service helper:

```text
backend/services/identity/services/helpers/password_hasher.go
```

The implementation uses Argon2id from `golang.org/x/crypto/argon2` with 64 MiB memory, 3 passes, 1 thread, a 16-byte random salt, and a 32-byte derived key. Encoded hashes carry algorithm/version/work parameters, and verification validates/bounds those parameters before running Argon2id. Hash comparison is constant-time.

## Repository boundary

`repository.PasswordCredentialRepository` owns only:

```text
Create
FindByAccountId
Update
```

No password verification logic belongs in the repository.

## Registration integration

Registration validates the new password before any cross-service write. After Foundation Person and Identity Account creation, it calls `PasswordCredentialService.Set`.

If `Set` fails after Account creation, Registration compensates by permanently deleting that newly created Account. It does not delete or deactivate the Foundation Person.

See `backend/services/identity/docs/registration.md`.

## Still future

```text
login endpoint
forgot-password endpoint
reset-password endpoint
change-password endpoint
rate limiting
credential lockout
session creation
access/refresh tokens
OTP/recovery proof
Google/Microsoft login
MFA/passkeys
```
