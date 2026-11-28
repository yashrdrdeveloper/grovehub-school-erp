package constants

// Offline provisioning locks all authority and policy writers; the CLI also
// requires stopped application writers. This is not a live mutation protocol.
const Lock = `LOCK TABLE identity.accounts, identity.roles, identity.role_personas,
identity.role_permissions, identity.account_roles, identity.features, identity.permissions,
identity.tenant_features, identity.approval_policies, identity.approval_policy_steps
IN SHARE ROW EXCLUSIVE MODE`

const Author = `SELECT count(*) FROM identity.accounts WHERE tenant_id = ? AND account_id = ? AND status = 'ACTIVE'`

const Version = `SELECT COALESCE(max(version),0) FROM identity.approval_policies WHERE tenant_id = ? AND request_type = ?`

const Eligible = `SELECT count(DISTINCT a.account_id)
FROM identity.accounts a
JOIN identity.account_roles ar ON ar.tenant_id = a.tenant_id AND ar.account_id = a.account_id
JOIN identity.roles r ON r.tenant_id = ar.tenant_id AND r.role_id = ar.role_id
JOIN identity.role_permissions rp ON rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id
JOIN identity.permissions p ON p.permission_id = rp.permission_id
LEFT JOIN identity.features f ON f.feature_id = p.feature_id
LEFT JOIN identity.tenant_features tf ON tf.tenant_id = a.tenant_id AND tf.feature_id = f.feature_id
WHERE a.tenant_id = ? AND a.status = 'ACTIVE' AND r.status = 'ACTIVE'
AND p.active AND p.module = 'IDENTITY' AND p.code = ?
AND ((? = '' AND p.feature_id IS NULL) OR (f.code = ? AND f.active AND tf.enabled))
AND (CAST(? AS text) IS NULL OR r.role_id::text = CAST(? AS text))`

const Deactivate = `UPDATE identity.approval_policies SET active = false WHERE tenant_id = ? AND request_type = ? AND active`
const Insert = `INSERT INTO identity.approval_policies(policy_id, tenant_id, request_type, version, active, created_by) VALUES (?, ?, ?, ?, true, ?)`
const InsertStep = `INSERT INTO identity.approval_policy_steps(tenant_id, policy_id, step_number, required_permission_code, required_role_id) VALUES (?, ?, ?, ?, ?)`
const Receipt = `INSERT INTO identity.approval_policy_receipts(receipt_id, tenant_id, operator_id, author_id, input_sha256, result) VALUES (?, ?, ?, ?, ?, CAST(? AS jsonb))`
