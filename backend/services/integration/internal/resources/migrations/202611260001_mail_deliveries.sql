-- +goose Up
CREATE TABLE integration.mail_deliveries (
    tenant_id text NOT NULL CHECK (tenant_id <> ''),
    delivery_id uuid NOT NULL,
    request_hash char(64) NOT NULL,
    state text NOT NULL CHECK (state IN ('DISPATCHING', 'ACCEPTED')),
    reserved_at timestamptz NOT NULL DEFAULT clock_timestamp(),
    accepted_at timestamptz,
    PRIMARY KEY (tenant_id, delivery_id),
    CHECK ((state = 'ACCEPTED' AND accepted_at IS NOT NULL) OR (state = 'DISPATCHING' AND accepted_at IS NULL))
);
-- +goose Down
-- Dropping dispatch reservations could cause duplicate mail after a rollback.
-- +goose StatementBegin
DO $$
BEGIN
    IF EXISTS (SELECT 1 FROM integration.mail_deliveries) THEN
        RAISE EXCEPTION 'mail delivery history must be preserved';
    END IF;
END;
$$;
-- +goose StatementEnd
DROP TABLE integration.mail_deliveries;
