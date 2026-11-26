package constants

const Insert = `INSERT INTO identity.approval_audit
(event_id, tenant_id, request_id, actor_id, request_version, operation, event_payload, occurred_at)
VALUES (?, ?, ?, ?, ?, ?, CAST(? AS jsonb), ?)`
