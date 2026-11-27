package search

import coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"

// Search metadata only: never search raw payload, digest or tenant identifier.
var Config=coresearch.Config{
	SearchExpression:"(a.asset_type || ' ' || a.asset_id::text || ' ' || a.audit_type || ' ' || a.operation || ' ' || a.source || ' ' || a.audit_id::text)",
	CursorExpression:"a.audit_id",
	ScalarFields:map[string]coresearch.ScalarField{
		"auditId":{Expression:"a.audit_id::text"},"assetId":{Expression:"a.asset_id::text"},"actorId":{Expression:"a.actor_id::text"},
		"personId":{Expression:"a.person_id::text"},"subjectAccountId":{Expression:"a.subject_account_id::text"},"occurrenceId":{Expression:"a.occurrence_id::text"},
		"assetType":{Expression:"a.asset_type"},"auditType":{Expression:"a.audit_type"},"operation":{Expression:"a.operation"},"source":{Expression:"a.source"},
		"occurredAt":{Expression:"a.occurred_at",Cast:"::timestamptz"},"ingestedAt":{Expression:"a.ingested_at",Cast:"::timestamptz"},
	},
	SortFields:map[string]coresearch.SortField{
		"auditId":{Expression:"a.audit_id",Cast:"::uuid"},"occurredAt":{Expression:"a.occurred_at",Cast:"::timestamptz"},"ingestedAt":{Expression:"a.ingested_at",Cast:"::timestamptz"},
		"assetType":{Expression:"a.asset_type",Cast:"::text"},"auditType":{Expression:"a.audit_type",Cast:"::text"},"operation":{Expression:"a.operation",Cast:"::text"},"source":{Expression:"a.source",Cast:"::text"},
	},
}
