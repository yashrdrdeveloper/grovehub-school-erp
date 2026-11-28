ALTER TABLE school_registry.schools ADD COLUMN version bigint NOT NULL DEFAULT 1 CHECK(version>0);
ALTER TABLE school_registry.rename_events ADD COLUMN previous_version bigint CHECK(previous_version>0);
CREATE FUNCTION school_registry.advance_school_version() RETURNS trigger LANGUAGE plpgsql AS $$
BEGIN
 IF (NEW.name,NEW.code,NEW.status) IS DISTINCT FROM (OLD.name,OLD.code,OLD.status) THEN
  NEW.version := OLD.version+1;
 ELSE
  NEW.version := OLD.version;
 END IF;
 RETURN NEW;
END;
$$;
CREATE TRIGGER school_version BEFORE UPDATE ON school_registry.schools
 FOR EACH ROW EXECUTE FUNCTION school_registry.advance_school_version();
