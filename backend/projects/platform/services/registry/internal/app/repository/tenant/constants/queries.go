package constants

const FindByID = "project_id = ? AND tenant_id = ?"

const ChangeStatus = `WITH parent AS (
	SELECT project_id, status FROM registry.projects WHERE project_id = ? FOR UPDATE
)
UPDATE registry.tenants AS t
SET status = ?, version = t.version + 1,
	session_generation = t.session_generation + CASE WHEN ? = 'DISABLED' THEN 1 ELSE 0 END,
	updated_at = now()
FROM parent
WHERE t.project_id = parent.project_id AND t.tenant_id = ?
	AND t.version = ? AND t.status <> ?
	AND (? = 'DISABLED' OR parent.status = 'ACTIVE')
RETURNING t.*`
