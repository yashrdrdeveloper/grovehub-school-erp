# Identity registration behavior

`POST /auth/register` accepts `firstName`, `middleName`, `lastName`, `email`, optional `phoneNumber`, and (currently) `password`, plus an `Idempotency-Key` header. The registration JSON does not accept `gid` or any unknown property.

The service normalizes the email, searches Foundation for an existing Person with that email, reuses and updates the same Person if found, or creates one if not found. It does not create a duplicate intentionally for a unique matching email. It returns both `gid` and `personGid` as the same Foundation identifier, alongside the Account ID.

**Atomic registration:** Person, Account and password credential writes share one PostgreSQL transaction; events publish only after commit. **Not implemented yet:** password-setup email and durable email/event outbox. Existing name updates by email alone should be reviewed against email-verification/ownership requirements before public exposure.

See `openapi.yaml` for the API contract.
