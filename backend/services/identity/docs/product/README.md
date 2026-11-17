# Identity — Product and API Workflows

Identity currently documents an **admin-invited registration flow** and **public credential authentication**.

- [Protected registration and activation](registration.md)
- [Invitation delivery, token redemption and password setup](password_setup.md)
- [Public login and session tokens](login.md)

## API surface

| Method | Path | Actor |
|---|---|---|
| POST | `/identity/protected/auth/register` | Verified and explicitly authorized administrator |
| POST | `/identity/public/auth/set-password` | Holder of valid single-use setup token |
| POST | `/identity/public/auth/login` | Anonymous user submitting email and password |
| POST | `/identity/public/auth/refresh` | Holder of valid, unexpired refresh credential |
| POST | `/identity/protected/auth/logout` | Verified ACTIVE session |

All external responses use the existing REST envelope; examples show data payloads rather than a captured wire response.

## Not yet available

Logout-all, admin resend invitation, password reset, MFA, OAuth2/OIDC, SAML and complete RBAC are not yet exposed as supported endpoints. Production rollout requires migration, mail-sandbox and concurrency verification.
