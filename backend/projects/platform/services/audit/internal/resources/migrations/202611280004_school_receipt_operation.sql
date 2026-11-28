-- +goose Up
ALTER TABLE audit.school_mutation_receipts ADD COLUMN operation text NOT NULL DEFAULT 'RENAME' CHECK(operation IN ('RENAME','ACTIVATE','DISABLE'));
-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
 RAISE EXCEPTION 'school mutation receipt operations must be preserved';
END;
$$;
-- +goose StatementEnd
