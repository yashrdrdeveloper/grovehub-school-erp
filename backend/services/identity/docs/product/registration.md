# Registration — Protected Admin Invitation

Registration is now an administrator-invited flow. **The registration request contains no password** and new accounts are created as `PENDING_SETUP`. The administrator must have an ACTIVE Identity session and be explicitly authorized.

For complete API requests, setup mail semantics, activation and replay behavior, see [Admin Invitation and Password Setup](password_setup.md).

For public sign-in, session-token issuance and refresh behavior, see [Login](login.md).

## Workflow

1. Admin signs in and presents a verified bearer access token.
2. Admin calls `POST /identity/protected/auth/register` with an idempotency key and person details.
3. Identity checks authenticated admin authorization, normalizes and fingerprints the request.
4. Foundation person resolution runs, Identity creates the pending account, and the setup token hash plus encrypted invitation job commit in the same Identity transaction.
5. The worker sends an invitation through a system-owned Integration `mail.send` connection.
6. The invited user opens the HTTPS link, submits the token plus their chosen password to `POST /identity/public/auth/set-password`, and becomes ACTIVE.
7. The user signs in using `POST /identity/public/auth/login`.

## Limitations

The system sender, setup URL, encryption key and admin allowlist must be configured. The initial admin must be bootstrapped independently. Full cross-service rollback between Identity and Foundation is not guaranteed. Resending an expired invitation, standardized operational metrics and rate limiting are follow-up work. These examples have not been executed against a live environment in this change.
