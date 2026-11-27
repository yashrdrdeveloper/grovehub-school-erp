package constants

const LockRole = `SELECT status, version FROM identity.roles
WHERE tenant_id = ? AND role_id = ? FOR UPDATE`

const CountAssignments = `SELECT count(*) FROM identity.account_roles WHERE tenant_id = ? AND role_id = ?`

const FindActivePermission = `SELECT p.permission_id::text AS id, p.module,
COALESCE(f.code, '') AS feature_code
FROM identity.permissions p LEFT JOIN identity.features f ON f.feature_id = p.feature_id
WHERE p.code = ? AND p.active AND (p.feature_id IS NULL OR f.active)
FOR SHARE OF p`

const DeletePermissions = `DELETE FROM identity.role_permissions WHERE tenant_id = ? AND role_id = ?`

const InsertPermission = `INSERT INTO identity.role_permissions
(tenant_id, role_id, permission_id, granted_by) VALUES (?, ?, ?, ?)`

const AdvancePermissionVersion = `UPDATE identity.roles
SET version = version + 1, updated_at = now(), updated_by = ?
WHERE tenant_id = ? AND role_id = ? AND version = ? AND status = 'INACTIVE'`
