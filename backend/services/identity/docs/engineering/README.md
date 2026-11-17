# Identity — Engineering Documentation

This module follows `docs/engineering/backend/backend.md`, `searchDetails.md` and `annotation.md`. Do not replicate generic SRP, ACID, DRY, pipeline, repository or messaging standards in this folder.

## Architecture references

- [Registration, login, tokens, password setup and system mail integration](authentication.md).
- [Product API workflows](../product/README.md).

## Current ownership

- Foundation owns person/contact records; Identity owns accounts, password credentials, one-time setup tokens and sessions.
- `controller/public/` holds login, token refresh and password setup.
- `controller/protected/` holds administrator registration.
- `services/impl/auth.go` orchestrates Identity use cases. The existing registration controller implementation is temporarily wrapped by the protected controller package.
- `repository/impl/` encapsulates GORM persistence, atomic consumption of setup tokens and DB constraints.
- `services/impl/password_setup_invitation.go` persists AES-GCM encrypted delivery jobs and processes them through a leased worker.
- `services/impl/system_mail.go` integrates with the existing `mail.send` action for a dedicated system-owned provider connection.
- `services/impl/session_issuer.go` persists and rotates opaque session bearer credentials; `session_auth.go` validates access credentials for protected routes.

## Technical debt / safety gates

The current administrative bootstrap policy uses `IDENTITY_ADMIN_ACCOUNT_IDS` rather than RBAC. No configured allowlist or mail sender means registration cannot succeed. Remote Foundation calls are not proven atomic with Identity writes; transaction boundaries and retry/compensation need testing. OAuth/SSO, recovery, logout, password reset, abuse controls, migrations and production config validation need dedicated work. Documentation is source-based; compilation and end-to-end tests have not been executed.
