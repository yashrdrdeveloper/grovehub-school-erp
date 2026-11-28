# First administrator bootstrap (operator only)

The normal `POST /identity/protected/accounts` route requires a verified administrator. On a new deployment, there is no administrator to perform the first registration. This command is an **offline, privileged, one-time exception**, not a public registration route. The command is not attached to the Launcher, router, Swagger or a network listener.

## Prerequisites

- An authorized operator with controlled execution access to the backend and a PostgreSQL credential capable of writing Foundation and Identity tables.
- The **existing** Foundation and Identity migrations applied. Do **not** reset/recreate schemas to bootstrap an administrator.
- `DATABASE_URL` configured in the repository root `.env`; its secrets should not be committed.
- An existing tenant identifier to assign to the initial admin. The command does **not** create a tenant or assign RBAC roles.
- No rows in `identity.accounts`. If any account already exists, the command refuses to run. Provisioning an admin into an existing database requires a separate audited recovery procedure.

## Run from repository root

```bash
make identity-bootstrap-admin EMAIL=admin@example.com FIRST_NAME=Initial LAST_NAME=Admin TENANT_ID=school-001
```

An optional `PERSON_GID=<UUID>` links an **existing ACTIVE Foundation Person** with the same email contact instead of creating a new person. If no Person GID is supplied, the command reuses the matching Foundation email contact or creates a new Person and primary email in the same PostgreSQL transaction.

The launcher prompts twice for the password without echoing it. Do not put passwords in command-line arguments or `.env`. Passwords must be 12–256 bytes and are hashed with Identity's existing Argon2id implementation.

## Grant admin allowlist access

On successful bootstrap, the CLI prints the new **Account ID**. Update the managed secret/configuration `IDENTITY_ADMIN_ACCOUNT_IDS` to include this ID (comma-separated if needed), then **restart** the backend to reload this configuration. The first account is ACTIVE and can log in with its email and bootstrap password. It will be an administrator only after the allowlist is loaded. Its tenant scope must match the scope used for future account management.

```bash
# Add the exact printed UUID to the securely managed runtime environment:
IDENTITY_ADMIN_ACCOUNT_IDS=<printed-account-uuid>
make dev-stop
make dev-start
```

Do not confuse `identity.accounts.account_id` with `foundation.persons.gid`: only the Account ID belongs in `IDENTITY_ADMIN_ACCOUNT_IDS`.

## Safeguards and limitations

- No HTTP endpoint; this command must run from an operator-controlled environment.
- Requires an explicit `--confirm-bootstrap` flag and a password over stdin; the interactive wrapper provides both.
- A transaction-scoped PostgreSQL advisory lock serializes simultaneous bootstrap commands. Within the transaction, the command checks that `identity.accounts` is empty before creating the person, account and credential.
- A failed transaction rolls back account and person writes. This assumes Foundation and Identity use the same PostgreSQL database for this controlled bootstrap.
- The command does not transmit mail, print a password, add permanent bypass routes, or grant privileges to subsequent accounts.
- Provisioning privileges are **not** automatic: `IDENTITY_ADMIN_ACCOUNT_IDS` must be controlled by deployment operators.
- Operator DB credentials are powerful: secure SSH/CI access, secret storage, PostgreSQL authorization, and deployment logs. Take an audit record of the change through your deployment control plane.
- This initial exception uses a password chosen at the terminal; it does **not** issue a one-time invitation. Standard accounts still use invitation-based setup.

## Verification

1. Run `make dev-start`.
2. POST `/identity/public/auth/login` with `{"identifier":"admin@example.com","password":"<your-password>"}` over an authorized local/dev client. Login should issue secure HttpOnly cookies.
3. Verify `GET /identity/protected/auth/me` using the issued session.
4. With CSRF/Origin protection satisfied, use `POST /identity/protected/accounts` for subsequent invitations.

Never publish the bootstrap password or a complete connection string in issue trackers, logs, or Swagger.


The single executable is built from `backend/projects/platform/cmd`; use its `identity-bootstrap-admin` subcommand. The privileged service-local facade lives at `backend/projects/platform/services/identity/internal/app/services/bootstrapadmin/operations/provision`. Do not recreate `services/identity/cmd`.
