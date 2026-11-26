-- +goose Up
CREATE TABLE accounting.accounting_entities (
 entity_id uuid PRIMARY KEY,
 tenant_id text NOT NULL CHECK (btrim(tenant_id) <> ''),
 code text NOT NULL CHECK (btrim(code) <> ''),
 name text NOT NULL CHECK (btrim(name) <> ''),
 created_at timestamptz NOT NULL DEFAULT now(),
 UNIQUE (tenant_id, entity_id),
 UNIQUE (tenant_id, code)
);
CREATE TABLE accounting.accounting_books (
 book_id uuid PRIMARY KEY,
 tenant_id text NOT NULL,
 entity_id uuid NOT NULL,
 code text NOT NULL CHECK (btrim(code) <> ''),
 name text NOT NULL CHECK (btrim(name) <> ''),
 base_currency text NOT NULL CHECK (base_currency = 'INR'),
 currency_scale smallint NOT NULL DEFAULT 2 CHECK (currency_scale = 2),
 accounting_basis text NOT NULL DEFAULT 'ACCRUAL' CHECK (accounting_basis = 'ACCRUAL'),
 created_at timestamptz NOT NULL DEFAULT now(),
 FOREIGN KEY (tenant_id, entity_id) REFERENCES accounting.accounting_entities (tenant_id, entity_id),
 UNIQUE (tenant_id, book_id),
 UNIQUE (tenant_id, entity_id, code)
);
-- +goose Down
DROP TABLE accounting.accounting_books;
DROP TABLE accounting.accounting_entities;
