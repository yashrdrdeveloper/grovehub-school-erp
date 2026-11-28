package constants

const FindByID = "project_id = ?"

const ChangeStatus = `WITH changed AS (
	UPDATE registry.projects
	SET status = ?, version = version + 1, updated_at = now()
	WHERE project_id = ? AND version = ? AND status <> ?
	RETURNING *
), invalidated AS (
	UPDATE registry.tenants AS t
	SET session_generation = t.session_generation + 1, version = t.version + 1, updated_at = now()
	FROM changed
	WHERE t.project_id = changed.project_id AND changed.status = 'DISABLED'
	RETURNING t.tenant_id
)
SELECT * FROM changed`
