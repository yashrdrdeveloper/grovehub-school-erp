# Internal Tool backend

Independent operator identity, approvals, API templates and execution. No School session reuse or direct Platform database access.

Set INTERNAL_TOOL_DATABASE_URL and INTERNAL_TOOL_ALLOWED_ORIGIN. Provision a separate non-superuser database role. Run make migrate explicitly before startup. Bootstrap the initial L2 once using a password of 12–256 bytes; pipe it from a secret manager or use a shell's silent read, e.g. bash: read -rs OPERATOR_PASSWORD; printf '%s' "$OPERATOR_PASSWORD" | make bootstrap-operator EMAIL=operator@example.com; unset OPERATOR_PASSWORD. Do not pass a password via argv or commit it. Migrations do not create accounts.

Operator login, me and logout are exposed at /internal/public/auth/login and /internal/protected/auth/{me,logout}. Session and CSRF cookies use the distinct __Host-gh_internal_* namespace, Secure and SameSite=Strict. Use HTTPS ingress (localhost development behavior depends on browser). Operator sessions expire after eight hours; logout revokes the stored digest. Login has a persistent per-email attempt budget and bounded Argon2 concurrency. Configure ingress rate limits and periodic cleanup of expired sessions and old attempt records before deployment.

GET /health/database checks the owned database; /health/ready remains false while approvals, templates, execution authorization and full application readiness are unfinished. L1/L2 are the initial personas; only first L2 bootstrap is implemented, not operator management. No Platform execution credentials are granted yet.
