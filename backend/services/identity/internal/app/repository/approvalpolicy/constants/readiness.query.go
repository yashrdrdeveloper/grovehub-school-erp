package constants

const LatestReceipt = `SELECT result::text FROM identity.approval_policy_receipts
WHERE tenant_id = ? ORDER BY created_at DESC, receipt_id DESC LIMIT 1`

const ActiveRevision = `SELECT policy_id::text AS policy_id, version, request_type AS type
FROM identity.approval_policies WHERE tenant_id = ? AND request_type = ? AND active`

const ActiveSteps = `SELECT step_number AS number, required_permission_code AS permission, required_role_id::text AS role_id
FROM identity.approval_policy_steps WHERE tenant_id = ? AND policy_id = ? ORDER BY step_number`
