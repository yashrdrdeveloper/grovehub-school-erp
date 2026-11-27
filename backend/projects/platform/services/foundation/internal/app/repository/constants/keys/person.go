package keys

// PersonIdempotencyKeyUniqueConstraint is the PostgreSQL uniqueness constraint
// used to prevent duplicate Person creation requests.
const PersonIdempotencyKeyUniqueConstraint = "uq_persons_create_idempotency_key"
