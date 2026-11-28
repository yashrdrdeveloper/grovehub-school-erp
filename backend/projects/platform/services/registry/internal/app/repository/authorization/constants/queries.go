package constants

// Resolve uses one statement to observe a consistent snapshot of the live key,
// workload, owning project, exact grant and target state. NULL is never a wildcard.
const Resolve = `SELECT w.workload_id, COALESCE(w.project_id::text, '') AS caller_project_id,
	w.issuer, w.subject, k.public_key, k.valid_from, k.valid_until, w.updated_at AS workload_updated_at
FROM registry.workload_keys k
JOIN registry.workloads w ON w.workload_id = k.workload_id
LEFT JOIN registry.projects owner ON owner.project_id = w.project_id
JOIN registry.workload_grants g ON g.workload_id = w.workload_id
LEFT JOIN registry.projects target ON target.project_id = g.target_project_id
LEFT JOIN registry.tenants tenant ON tenant.project_id = g.target_project_id AND tenant.tenant_id = g.target_tenant_id
WHERE k.key_id = ?
	AND k.revoked_at IS NULL AND k.valid_from <= now() AND k.valid_until > now()
	AND w.status = 'ACTIVE' AND (w.project_id IS NULL OR owner.status = 'ACTIVE')
	AND g.scope = ? AND g.capability = ?
	AND g.target_project_id IS NOT DISTINCT FROM NULLIF(?, '')::uuid
	AND g.target_tenant_id IS NOT DISTINCT FROM NULLIF(?, '')::uuid
	AND (? OR ((g.target_project_id IS NULL OR target.status = 'ACTIVE')
		AND (g.target_tenant_id IS NULL OR tenant.status = 'ACTIVE')))`

const TenantState = `SELECT (p.status = 'ACTIVE' AND t.status = 'ACTIVE') AS active,
	t.session_generation
FROM registry.tenants t
JOIN registry.projects p ON p.project_id = t.project_id
WHERE t.project_id = ? AND t.tenant_id = ?`
