-- +goose Up
ALTER TABLE audit.operator_activities ADD COLUMN correlation_id uuid;
CREATE INDEX idx_operator_activity_correlation ON audit.operator_activities(correlation_id) WHERE correlation_id IS NOT NULL;
-- +goose Down
-- +goose StatementBegin
DO $$
BEGIN
 RAISE EXCEPTION 'operator audit correlation history must be preserved';
END;
$$;
-- +goose StatementEnd
