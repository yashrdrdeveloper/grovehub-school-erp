CREATE TABLE school_registry.schools (
 id uuid PRIMARY KEY,
 code text NOT NULL UNIQUE CHECK(code ~ '^[a-z][a-z0-9_-]{1,63}$'),
 name text NOT NULL CHECK(char_length(name) BETWEEN 1 AND 150),
 status text NOT NULL DEFAULT 'ACTIVE' CHECK(status IN ('ACTIVE','DISABLED')),
 created_at timestamptz NOT NULL DEFAULT now()
);
