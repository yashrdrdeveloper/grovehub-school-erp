# School backend

Owns browser-facing School and Identity use-case APIs. Platform calls use the restricted authenticated generated SDK clients. No direct Platform database access.

Startup requires SCHOOL_DATABASE_URL for a separately provisioned PostgreSQL database and least-privilege role. No fallback to DATABASE_URL or other project credentials. GET /health/database checks only the owned database; full application readiness remains false until the remaining School dependencies and end-to-end checks are implemented. No School business migrations exist yet.
