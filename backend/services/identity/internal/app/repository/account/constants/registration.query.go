package constants

const FindRegistration = `SELECT a.account_id::text AS account_id, a.person_gid::text AS person_gid, a.email, r.request_sha256 AS digest
FROM identity.registration_receipts r JOIN identity.accounts a ON a.tenant_id = r.tenant_id AND a.account_id = r.account_id
WHERE r.tenant_id = ? AND r.actor_id = ? AND r.idempotency_key = ?`

const InsertPending = `INSERT INTO identity.accounts
(account_id, person_gid, tenant_id, email, phone, status, version, search_details, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, 'PENDING_SETUP', 1, ?, ?, ?)`

const InsertRegistrationReceipt = `INSERT INTO identity.registration_receipts
(tenant_id, actor_id, idempotency_key, request_sha256, account_id, event_id, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?)`
