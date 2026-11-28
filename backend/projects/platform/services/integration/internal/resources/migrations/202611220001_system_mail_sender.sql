-- +goose Up
CREATE TABLE integration.system_mail_senders (
    tenant_id UUID PRIMARY KEY,
    connection_id UUID NOT NULL REFERENCES integration.connections(id),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- +goose Down
DROP TABLE integration.system_mail_senders;
