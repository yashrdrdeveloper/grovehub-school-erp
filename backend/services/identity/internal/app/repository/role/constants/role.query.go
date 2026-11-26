package constants

const InsertDraft = `INSERT INTO identity.roles
(role_id, tenant_id, name, normalized_name, description, search_details, created_by, updated_by)
VALUES (?, ?, ?, ?, ?, ?, ?, ?)`

const InsertPersona = `INSERT INTO identity.role_personas (tenant_id, role_id, persona_code)
VALUES (?, ?, ?)`

const FindDetail = `SELECT r.*, COALESCE(
    (SELECT json_agg(p.persona_code ORDER BY p.persona_code)::text
     FROM identity.role_personas p WHERE p.tenant_id = r.tenant_id AND p.role_id = r.role_id),
    '[]') AS persona_json
FROM identity.roles r WHERE r.tenant_id = ? AND r.role_id = ?`
