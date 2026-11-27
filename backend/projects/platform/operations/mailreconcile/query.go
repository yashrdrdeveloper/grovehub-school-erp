package mailreconcile

// Filter each side by the exact tenant BEFORE joining. An Integration receipt
// without a same-tenant Identity row must not reveal any foreign account data.
// A receipt confirms provider acceptance only, not inbox delivery.
const reconciliationQuery = `
WITH jobs AS (
	SELECT i.id AS delivery_id, i.account_id, i.purpose, i.sent_at,
	       i.suppressed_at, i.suppression_reason, i.lease_until
	FROM identity.invitations i
	JOIN identity.accounts a ON a.account_id = i.account_id
	WHERE a.tenant_id = ? AND i.id > CAST(? AS uuid)
), receipts AS (
	SELECT delivery_id, state, reserved_at, accepted_at
	FROM integration.mail_deliveries
	WHERE tenant_id = ? AND delivery_id > CAST(? AS uuid)
), joined AS (
	SELECT COALESCE(j.delivery_id, r.delivery_id)::text AS delivery_id,
	       j.account_id::text AS account_id,
	       CASE WHEN j.purpose IN ('ACCOUNT_SETUP','PASSWORD_RESET') THEN j.purpose
	            WHEN j.purpose IS NOT NULL THEN 'UNKNOWN' END AS purpose,
	       CASE WHEN j.delivery_id IS NULL THEN 'MISSING'
	            WHEN j.sent_at IS NOT NULL THEN 'SENT'
	            WHEN j.suppressed_at IS NOT NULL THEN 'SUPPRESSED'
	            ELSE 'PENDING' END AS identity_state,
	       COALESCE(r.state, 'MISSING') AS integration_state,
	       j.suppression_reason, r.reserved_at, r.accepted_at,
	       CASE
	         WHEN j.delivery_id IS NULL THEN 'RECEIPT_WITHOUT_JOB'
	         WHEN j.suppression_reason = 'DELIVERY_UNCERTAIN' THEN 'REVIEW_UNCERTAIN'
	         WHEN r.state = 'ACCEPTED' AND j.sent_at IS NULL THEN 'ACCEPTED_WITHOUT_SENT'
	         WHEN r.state = 'DISPATCHING' AND j.sent_at IS NOT NULL THEN 'SENT_WITHOUT_ACCEPTANCE'
	         WHEN r.state = 'DISPATCHING'
	              AND r.reserved_at < CURRENT_TIMESTAMP - INTERVAL '2 minutes'
	              AND (j.lease_until IS NULL OR j.lease_until < CURRENT_TIMESTAMP)
	           THEN 'STALE_DISPATCH'
	         ELSE NULL
	       END AS disposition
	FROM jobs j FULL OUTER JOIN receipts r USING (delivery_id)
)
SELECT delivery_id, account_id, purpose, identity_state, integration_state,
       suppression_reason, reserved_at, accepted_at, disposition
FROM joined
WHERE disposition IS NOT NULL
ORDER BY delivery_id
LIMIT ?
`
