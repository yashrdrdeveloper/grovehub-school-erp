package constants

// Explicit projection avoids duplicate legacy and normalized column aliases.
// Inactive assignments remain visible as assignments, not effective grants.
const NormalizedProjection = `a.account_id, a.person_gid, a.tenant_id, a.username, a.email,
a.phone, a.status, a.version, a.search_details, a.created_at, a.updated_at,
(SELECT ar.role_id::text FROM identity.account_roles ar
 WHERE ar.tenant_id = a.tenant_id AND ar.account_id = a.account_id AND ar.assignment_type = 'PRIMARY') AS primary_role_id,
COALESCE((SELECT json_agg(ar.role_id::text ORDER BY ar.role_id)::jsonb
 FROM identity.account_roles ar
 WHERE ar.tenant_id = a.tenant_id AND ar.account_id = a.account_id AND ar.assignment_type = 'SECONDARY'), '[]'::jsonb) AS secondary_role_ids`
