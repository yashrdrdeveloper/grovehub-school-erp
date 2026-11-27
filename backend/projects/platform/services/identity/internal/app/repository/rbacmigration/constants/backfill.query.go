package constants

// Maintenance-only table locks block legacy and normalized DML for the short
// backfill transaction. Acquire before reading inventory; never call online.
const LockTables = `LOCK TABLE identity.accounts, identity.roles, identity.role_personas,
identity.role_permissions, identity.account_roles IN SHARE ROW EXCLUSIVE MODE`

const InsertRole = `INSERT INTO identity.roles
(role_id, tenant_id, name, normalized_name, description, search_details, created_by, updated_by)
VALUES (?, ?, ?, ?, '', ?, ?, ?)`

const InsertPersona = `INSERT INTO identity.role_personas (tenant_id, role_id, persona_code) VALUES (?, ?, ?)`

const InsertAssignment = `INSERT INTO identity.account_roles
(tenant_id, account_id, role_id, assignment_type, granted_by) VALUES (?, ?, ?, ?, ?)`

const VerifyRole = `SELECT count(*) FROM identity.roles r
WHERE r.tenant_id = ? AND r.role_id = ? AND r.name = ? AND r.normalized_name = ?
AND r.status = 'INACTIVE' AND r.version = 1
AND (SELECT count(*) FROM identity.role_personas rp WHERE rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id) = 1
AND EXISTS (SELECT 1 FROM identity.role_personas rp WHERE rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id AND rp.persona_code = ?)
AND NOT EXISTS (SELECT 1 FROM identity.role_permissions rp WHERE rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id)`
