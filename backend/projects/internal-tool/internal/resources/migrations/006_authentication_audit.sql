ALTER TABLE internal_tool.workflow_events ADD COLUMN correlation_id uuid;
CREATE INDEX workflow_events_correlation ON internal_tool.workflow_events(correlation_id) WHERE correlation_id IS NOT NULL;
-- Authentication occurrence counts for attempts with no verified actor.
-- No identifier, email hash, IP address, credential or token is retained.
CREATE TABLE internal_tool.authentication_counters (
 bucket timestamptz NOT NULL,
 outcome text NOT NULL CHECK(outcome IN ('REJECTED','LIMITED')),
 attempts bigint NOT NULL CHECK(attempts>0),
 PRIMARY KEY(bucket,outcome)
);
