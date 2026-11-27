package constants

const LockRegistry = `LOCK TABLE identity.features, identity.permissions,
identity.tenant_features IN SHARE ROW EXCLUSIVE MODE`

const InsertGrant = `INSERT INTO identity.role_permissions (tenant_id, role_id, permission_id, granted_by)
VALUES (?, ?, ?, ?)`

const ActivateRole = `UPDATE identity.roles SET status = 'ACTIVE', version = version + 1,
updated_at = now(), updated_by = ? WHERE tenant_id = ? AND role_id = ? AND version = 1 AND status = 'INACTIVE'`
