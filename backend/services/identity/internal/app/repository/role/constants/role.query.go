package constants

const InsertDraft = `INSERT INTO identity.roles
(role_id, tenant_id, name, normalized_name, description, search_details, created_by, updated_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

const InsertPersona = `INSERT INTO identity.role_personas (tenant_id, role_id, persona_code)
VALUES (?, ?, ?)`

const DetailProjection = `r.*, COALESCE(
    (SELECT json_agg(p.persona_code ORDER BY p.persona_code)::text
     FROM identity.role_personas p WHERE p.tenant_id = r.tenant_id AND p.role_id = r.role_id),
    '[]') AS persona_json`

const FindDetail = `SELECT ` + DetailProjection + ` FROM identity.roles r WHERE r.tenant_id = ? AND r.role_id = ?`

const UpdateMetadata = `UPDATE identity.roles SET name = ?, normalized_name = ?, description = ?,
search_details = ?, updated_by = ?, updated_at = now(), version = version + 1
WHERE tenant_id = ? AND role_id = ? AND version = ?`

const DeletePersonas = `DELETE FROM identity.role_personas WHERE tenant_id = ? AND role_id = ?`

const FindPermissions = `SELECT r.role_id::text AS role_id, r.version, COALESCE(
    (SELECT json_agg(p.code ORDER BY p.code)::text FROM identity.role_permissions rp
     JOIN identity.permissions p ON p.permission_id = rp.permission_id
     WHERE rp.tenant_id = r.tenant_id AND rp.role_id = r.role_id), '[]') AS codes_json
FROM identity.roles r WHERE r.tenant_id = ? AND r.role_id = ?`
