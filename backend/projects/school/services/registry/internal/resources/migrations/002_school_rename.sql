CREATE TABLE school_registry.rename_events (
 school_id uuid NOT NULL REFERENCES school_registry.schools(id),
 execution_id uuid PRIMARY KEY,
 workload_subject text NOT NULL,
 previous_name text NOT NULL,
 new_name text NOT NULL,
 occurred_at timestamptz NOT NULL DEFAULT clock_timestamp(),
 CHECK(previous_name<>new_name)
);
CREATE FUNCTION school_registry.reject_rename_event_mutation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 RAISE EXCEPTION 'school rename events are immutable';
END;
$$;
CREATE TRIGGER rename_events_immutable BEFORE UPDATE OR DELETE ON school_registry.rename_events
 FOR EACH ROW EXECUTE FUNCTION school_registry.reject_rename_event_mutation();
CREATE TRIGGER rename_events_no_truncate BEFORE TRUNCATE ON school_registry.rename_events
 FOR EACH STATEMENT EXECUTE FUNCTION school_registry.reject_rename_event_mutation();
