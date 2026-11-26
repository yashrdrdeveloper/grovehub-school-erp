package constants

// Every branch is scoped before union. Approval events target approval requests;
// resource effects target the affected account or role. No payload is selected.
const Timeline = `
SELECT event_id, source, actor_id, resource_type, resource_id, request_id, version, operation, occurred_at
FROM (
	SELECT event_id, 'APPROVAL'::text AS source, actor_id, 'APPROVAL'::text AS resource_type,
	       request_id AS resource_id, request_id, request_version AS version, operation, occurred_at
	FROM identity.approval_audit WHERE tenant_id = ?
	UNION ALL
	SELECT event_id, 'RESOURCE', actor_id, resource_type, resource_id, request_id, resource_version, operation, occurred_at
	FROM identity.resource_audit WHERE tenant_id = ?
	UNION ALL
	SELECT event_id, 'ROLE_METADATA', actor_id, 'ROLE', role_id, NULL::uuid, role_version, operation, occurred_at
	FROM identity.role_metadata_audit WHERE tenant_id = ?
	UNION ALL
	SELECT event_id, 'ACCOUNT_METADATA', actor_id, 'ACCOUNT', account_id, NULL::uuid, account_version, operation, occurred_at
	FROM identity.account_metadata_audit WHERE tenant_id = ?
	UNION ALL
	SELECT event_id, 'REGISTRATION', actor_id, 'ACCOUNT', account_id, NULL::uuid, 1::bigint, 'CREATED', created_at
	FROM identity.registration_receipts WHERE tenant_id = ?
) AS events WHERE TRUE
`
