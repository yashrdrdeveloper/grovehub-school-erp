# Internal Tool backend

Independent operator authentication, approvals, API templates and execution. No School session reuse or direct Platform database access.

Startup requires INTERNAL_TOOL_DATABASE_URL for a separately provisioned PostgreSQL database and least-privilege role. No fallback to DATABASE_URL or School credentials. GET /health/database checks the owned pool. Full readiness remains false: operator authentication, business migrations and services are not implemented yet.
