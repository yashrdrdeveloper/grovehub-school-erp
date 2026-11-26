# Accounting service

Independent accounting module. Owns books, chart of accounts, journals, ledger,
periods and financial statements. Payment and Finance are separate domains.

This initial scaffold follows Identity's public entrypoint and private bootstrap
builder/runtime pattern. No financial routes or client operations are exposed yet.

Next: launcher and migration command registration, accounting entities and books,
RBAC, generated repositories, posting engine and integration tests.
