-- +goose Up
CREATE INDEX school_mutation_receipts_school_operation_order ON audit.school_mutation_receipts(school_id,operation,occurred_at DESC,execution_id DESC);
-- +goose Down
DROP INDEX audit.school_mutation_receipts_school_operation_order;
