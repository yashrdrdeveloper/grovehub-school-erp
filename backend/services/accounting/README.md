# Accounting service

Independent accounting module. Owns books, chart of accounts, journals, ledger,
periods and financial statements. Payment and Finance are separate domains.

This initial scaffold follows Identity's public entrypoint and private bootstrap
builder/runtime pattern. No financial routes or client operations are exposed yet.

Registered in the platform launcher and migration command (accounting and all).
No financial HTTP routes are exposed.

Next: accounting entities and books,
RBAC, generated repositories, posting engine and integration tests.
