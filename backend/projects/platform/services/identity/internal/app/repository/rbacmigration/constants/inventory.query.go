package constants

const Accounts = `SELECT account_id::text AS account_id, version, primary_role_id AS primary_code,
secondary_role_ids::text AS secondary_json FROM identity.accounts
WHERE tenant_id = ? ORDER BY account_id`

const Roles = `SELECT role_id::text AS role_id, normalized_name FROM identity.roles
WHERE tenant_id = ? ORDER BY role_id`

const AssignmentCount = `SELECT count(*) FROM identity.account_roles WHERE tenant_id = ?`
