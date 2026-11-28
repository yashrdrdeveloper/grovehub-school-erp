-- +goose Up
CREATE INDEX school_mutation_receipts_school_order ON audit.school_mutation_receipts(school_id,occurred_at DESC,execution_id DESC);
CREATE INDEX school_mutation_receipts_operation_order ON audit.school_mutation_receipts(operation,occurred_at DESC,execution_id DESC);
-- +goose Down
DROP INDEX audit.school_mutation_receipts_operation_order;
DROP INDEX audit.school_mutation_receipts_school_order;
