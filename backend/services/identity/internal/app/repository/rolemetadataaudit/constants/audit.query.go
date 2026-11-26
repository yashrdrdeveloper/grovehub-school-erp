package constants

const Insert = `INSERT INTO identity.role_metadata_audit
(event_id, tenant_id, actor_id, role_id, role_version, operation, event_payload, occurred_at)
VALUES (?, ?, ?, ?, ?, ?, CAST(? AS jsonb), ?)`
