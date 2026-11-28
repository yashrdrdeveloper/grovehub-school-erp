package constants

const SummaryProjection = "r.request_id, r.requester_id, r.request_type, r.status, r.current_step, r.step_count, r.version, r.created_at, r.updated_at, r.expires_at"

const Inbox = `r.status = 'PENDING' AND r.expires_at > ? AND r.requester_id <> ?
AND r.request_type IN ('ACCOUNT_LIFECYCLE_CHANGE','ROLE_LIFECYCLE_CHANGE','ROLE_PERMISSION_CHANGE','ROLE_ASSIGNMENT_CHANGE')
AND (r.request_type <> 'ROLE_ASSIGNMENT_CHANGE' OR lower(r.payload->>'accountId') <> ?)
AND EXISTS (
 SELECT 1 FROM identity.approval_request_steps s
 JOIN jsonb_to_recordset(CAST(? AS jsonb)) AS g(role_id text, permission_code text)
 ON g.permission_code = s.required_permission_code
 AND (s.required_role_id IS NULL OR s.required_role_id::text = g.role_id)
 WHERE s.tenant_id = r.tenant_id AND s.request_id = r.request_id AND s.step_number = r.current_step
)`
