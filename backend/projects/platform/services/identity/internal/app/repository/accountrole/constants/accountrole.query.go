package constants

const FindAssignment = `SELECT a.account_id::text AS account_id, a.status, a.version,
COALESCE((SELECT ar.role_id::text FROM identity.account_roles ar
WHERE ar.tenant_id = a.tenant_id AND ar.account_id = a.account_id AND ar.assignment_type = 'PRIMARY'), '') AS primary_role_id,
COALESCE((SELECT json_agg(ar.role_id::text ORDER BY ar.role_id)::text FROM identity.account_roles ar
WHERE ar.tenant_id = a.tenant_id AND ar.account_id = a.account_id AND ar.assignment_type = 'SECONDARY'), '[]') AS secondary_json
FROM identity.accounts a WHERE a.tenant_id = ? AND a.account_id = ?`

const FindCandidates = `SELECT r.role_id::text AS role_id, r.status, r.version,
COALESCE((SELECT json_agg(p.code ORDER BY p.code)::text FROM identity.role_permissions rp
JOIN identity.permissions p ON p.permission_id = rp.permission_id
WHERE rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id), '[]') AS codes_json
FROM identity.roles r WHERE r.tenant_id = ? AND r.role_id::text IN ?
ORDER BY r.role_id`
