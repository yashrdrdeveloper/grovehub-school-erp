CREATE TABLE school_registry.status_events (
 execution_id uuid PRIMARY KEY,
 school_id uuid NOT NULL REFERENCES school_registry.schools(id),
 workload_subject text NOT NULL,
 previous_status text NOT NULL CHECK(previous_status IN ('ACTIVE','DISABLED')),
 new_status text NOT NULL CHECK(new_status IN ('ACTIVE','DISABLED')),
 previous_version bigint NOT NULL CHECK(previous_version>0),
 correlation_id uuid,
 occurred_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 CHECK(previous_status<>new_status)
);
CREATE INDEX status_events_school ON school_registry.status_events(school_id,occurred_at DESC);
CREATE TRIGGER status_events_immutable BEFORE UPDATE OR DELETE ON school_registry.status_events
 FOR EACH ROW EXECUTE FUNCTION school_registry.reject_rename_event_mutation();
CREATE TRIGGER status_events_no_truncate BEFORE TRUNCATE ON school_registry.status_events
 FOR EACH STATEMENT EXECUTE FUNCTION school_registry.reject_rename_event_mutation();
