package constants

const Insert = `INSERT INTO identity.resource_audit
(event_id, tenant_id, request_id, actor_id, resource_type, resource_id, resource_version, operation, event_payload, occurred_at)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, CAST(? AS jsonb), ?)`
