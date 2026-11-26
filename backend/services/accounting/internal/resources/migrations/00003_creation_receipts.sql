-- +goose Up
CREATE TABLE accounting.creation_receipts (
 tenant_id text NOT NULL,
 operation text NOT NULL CHECK (operation IN ('ENTITY_CREATE','BOOK_CREATE')),
 idempotency_key text NOT NULL CHECK (length(idempotency_key) BETWEEN 1 AND 128),
 request_hash text NOT NULL,
 asset_id uuid,
 PRIMARY KEY (tenant_id, operation, idempotency_key)
);
-- +goose Down
DROP TABLE accounting.creation_receipts;
