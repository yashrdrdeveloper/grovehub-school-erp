package constants

const Insert = `INSERT INTO identity.account_metadata_audit
(event_id, tenant_id, actor_id, account_id, account_version, operation, event_payload, occurred_at)
VALUES (?, ?, ?, ?, ?, 'UPDATED', CAST(? AS jsonb), ?)`
