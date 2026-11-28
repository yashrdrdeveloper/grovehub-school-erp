CREATE TABLE internal_tool.execution_reconciliations (
 execution_id uuid PRIMARY KEY REFERENCES internal_tool.executions(id),
 reviewer_id uuid NOT NULL REFERENCES internal_tool.operators(id),
 outcome text NOT NULL CHECK(outcome IN ('CONFIRMED_SUCCEEDED','CONFIRMED_FAILED','UNRESOLVED')),
 evidence text NOT NULL CHECK(length(evidence) BETWEEN 1 AND 1000),
 created_at timestamptz NOT NULL DEFAULT now()
);
CREATE INDEX internal_tool_execution_running ON internal_tool.executions(created_at) WHERE status='RUNNING';
