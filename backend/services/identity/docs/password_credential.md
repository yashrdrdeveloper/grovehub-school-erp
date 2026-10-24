# Identity Password Credential Contract

## Scope

Password credentials are internal Identity state used by higher-level authentication flows. They are not an HTTP resource and do not have standalone password CRUD endpoints.

There must be no routes such as:

```text
POST /password
PUT /password
DELETE /password
```

Future flows call the internal password credential service instead:

```text
Register              -> PasswordCredentialService.Set
Login                 -> PasswordCredentialService.Verify
Forgot/Reset Password -> PasswordCredentialService.Replace
Change Password       -> verify current password in the flow, then Replace
```

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

`identity.password_credentials.account_id` is both the primary key and a foreign key to `identity.accounts.account_id`. The relationship is one-to-zero-or-one.

Deleting an Identity Account cascades to its password credential.

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

The credential entity is internal. Service methods do not return it to callers, which avoids exposing the password hash outside credential orchestration.

## Internal service contract

`services.PasswordCredentialService` exposes only the capabilities needed by higher-level flows:

```text
Set
Verify
Replace
```

`Set` creates the first password credential for an existing account.

`Verify` compares a candidate password with the stored encoded hash and returns a boolean. It does not decide whether login is allowed; the future login flow owns account status, throttling, generic invalid-credential responses, and session creation.

`Replace` stores a new password hash and increments the credential version. It does not itself prove that a caller is authorized to reset/change the password; the higher-level reset/change flow must establish that first.

There is no standalone Remove/Delete password use case in this phase. Account deletion cascades automatically. A future provider-only conversion flow may introduce password removal only when that business requirement exists.

## Validator boundary

`validators.PasswordCredentialValidator` owns request/input validity only.

For a newly set/replaced password, the current policy is:

```text
valid UTF-8
minimum 8 Unicode code points
maximum 128 Unicode code points
```

Passwords are never trimmed, lowercased, or otherwise normalized. The exact submitted value is hashed.

Verification uses a candidate-password check rather than the current minimum-length policy. This prevents a future password-policy increase from making an older valid credential impossible to verify. Candidate verification still rejects invalid/empty/oversized input before expensive hashing.

## Password hashing helper

Password hashing is an internal Identity service helper:

```text
backend/services/identity/services/helpers/password_hasher.go
```

It is intentionally not a separate interface/package while there is only one implementation.

The current implementation uses Argon2id from `golang.org/x/crypto/argon2` with:

```text
memory:  64 MiB
passes:  3
threads: 1
salt:    16 random bytes
key:     32 bytes
```

The encoded hash stores the algorithm/version/work parameters with the salt and derived hash. Verification therefore uses the parameters stored with the credential rather than assuming that every existing password was created with the current defaults.

Hash parsing validates and bounds memory, time, thread, and output-size parameters before running Argon2id so malformed stored data cannot request unbounded work.

Hash comparison uses constant-time comparison.

## Repository boundary

`repository.PasswordCredentialRepository` owns only persistence primitives:

```text
Create
FindByAccountId
Update
```

The repository maps:

```text
primary-key conflict -> password credential already exists
foreign-key conflict -> Identity account missing
no row               -> password credential not found
```

No password verification logic belongs in the repository.

## Security ownership

This phase intentionally does not implement:

```text
login endpoint
registration endpoint
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

Those higher-level authentication flows will consume the internal credential capability later.
