ALTER TABLE school_registry.schools
 ADD COLUMN session_generation bigint NOT NULL DEFAULT 1 CHECK(session_generation > 0);

-- The generation advances atomically with a disable, including direct SQL
-- status updates. Activation and metadata changes cannot restore a generation.
CREATE FUNCTION school_registry.advance_session_generation() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF OLD.status = 'ACTIVE' AND NEW.status = 'DISABLED' THEN
  NEW.session_generation := OLD.session_generation + 1;
 ELSE
  NEW.session_generation := OLD.session_generation;
 END IF;
 RETURN NEW;
END;
$$;
CREATE TRIGGER school_session_generation BEFORE UPDATE ON school_registry.schools
 FOR EACH ROW EXECUTE FUNCTION school_registry.advance_session_generation();
