package constants

const LockAccount = `SELECT account_id::text AS account_id FROM identity.accounts
WHERE tenant_id = ? AND account_id = ? FOR UPDATE`

const LockCandidates = `SELECT role_id::text AS role_id FROM identity.roles
WHERE tenant_id = ? AND role_id::text IN ? ORDER BY role_id FOR SHARE`

const DeleteAssignments = `DELETE FROM identity.account_roles WHERE tenant_id = ? AND account_id = ?`

const InsertAssignment = `INSERT INTO identity.account_roles
(tenant_id, account_id, role_id, assignment_type, granted_by) VALUES (?, ?, ?, ?, ?)`

const AdvanceAccountVersion = `UPDATE identity.accounts SET version = version + 1, updated_at = now()
WHERE tenant_id = ? AND account_id = ? AND version = ? AND status IN ('ACTIVE', 'PENDING_SETUP')`
