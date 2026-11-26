package constants

const Now = `SELECT clock_timestamp()`

const ActivePolicy = `SELECT policy_id::text AS policy_id FROM identity.approval_policies
WHERE tenant_id = ? AND request_type = ? AND active FOR SHARE`

const PolicySteps = `SELECT step_number AS number, required_permission_code AS permission,
required_role_id::text AS role_id FROM identity.approval_policy_steps
WHERE tenant_id = ? AND policy_id = ? ORDER BY step_number FOR SHARE`

const FindByKey = `SELECT * FROM identity.approval_requests WHERE tenant_id = ? AND requester_id = ? AND idempotency_key = ?`
const LockRequest = `SELECT * FROM identity.approval_requests WHERE tenant_id = ? AND request_id = ? FOR UPDATE`

const RequestStep = `SELECT step_number AS number, required_permission_code AS permission,
required_role_id::text AS role_id FROM identity.approval_request_steps
WHERE tenant_id = ? AND request_id = ? AND step_number = ?`

const InsertRequest = `INSERT INTO identity.approval_requests
(request_id, tenant_id, policy_id, requester_id, request_type, payload, payload_sha256,
 reason, idempotency_key, status, current_step, step_count, version, expires_at, created_at, updated_at)
VALUES (?, ?, ?, ?, ?, CAST(? AS jsonb), ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

const InsertStep = `INSERT INTO identity.approval_request_steps
(tenant_id, request_id, step_number, required_permission_code, required_role_id) VALUES (?, ?, ?, ?, ?)`

const InsertDecision = `INSERT INTO identity.approval_decisions
(tenant_id, request_id, requester_id, step_number, reviewer_id, decision, comment, created_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

const Advance = `UPDATE identity.approval_requests SET status = ?, current_step = ?, version = ?, updated_at = ?
WHERE tenant_id = ? AND request_id = ? AND version = ? AND status = ? AND current_step = ?`
