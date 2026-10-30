-- +goose Up
DROP SCHEMA IF EXISTS infrastructure CASCADE;

CREATE SCHEMA infrastructure;

CREATE TABLE infrastructure.dkim_keys (
    domain text NOT NULL,
    selector text NOT NULL,
    private_key_pem text NOT NULL,
    public_key_pem text NOT NULL,
    created_at timestamptz NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT pk_infrastructure_dkim_keys PRIMARY KEY (domain),
    CONSTRAINT chk_infrastructure_dkim_keys_domain CHECK (length(btrim(domain)) > 0),
    CONSTRAINT chk_infrastructure_dkim_keys_selector CHECK (length(btrim(selector)) > 0),
    CONSTRAINT chk_infrastructure_dkim_keys_private_key CHECK (length(btrim(private_key_pem)) > 0),
    CONSTRAINT chk_infrastructure_dkim_keys_public_key CHECK (length(btrim(public_key_pem)) > 0)
);

CREATE FUNCTION infrastructure.reject_dkim_key_mutation()
RETURNS trigger
LANGUAGE plpgsql
AS $$
BEGIN
    RAISE EXCEPTION 'infrastructure DKIM identities are immutable; create an explicit key-rotation workflow instead';
END;
$$;

CREATE TRIGGER trg_infrastructure_dkim_keys_immutable
BEFORE UPDATE OR DELETE ON infrastructure.dkim_keys
FOR EACH ROW
EXECUTE FUNCTION infrastructure.reject_dkim_key_mutation();

-- +goose Down
DROP SCHEMA IF EXISTS infrastructure CASCADE;
