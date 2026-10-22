# Security Guidelines

> **Status:** Mandatory Phase 0 engineering standard  
> **Location:** `docs/engineering/security/SECURITY_GUIDELINES.md`  
> **Applies to:** authentication, authorization, sessions/tokens, secrets, provider credentials, input/output handling, tenant isolation, APIs, frontend, backend, scripts, containers, supply chain, audit, dependencies, logging, storage, file handling, network calls, and secure development  
> **Core principles:** least privilege, deny by default, trusted-context derivation, defense in depth, secret minimization, safe defaults, explicit trust boundaries, auditable mutations, dependency hygiene, and no security-by-UI

---

## 1. Objective

Security is a cross-cutting product property.

No layer may assume another layer will "handle security later."

The core model is:

```text
authenticate
    ↓
derive trusted principal/context
    ↓
authorize action/resource
    ↓
execute least-privilege operation
    ↓
audit security-relevant outcome
```

---

## 2. Mandatory Rules

1. Backend authorization is authoritative.
2. Frontend permission checks are UX only.
3. Deny by default.
4. Authentication context is derived from validated credentials, not arbitrary headers.
5. Tenant/org/project scope is verified server-side.
6. Secrets are never committed.
7. Secrets are never logged.
8. Prefer secret references/workload identity over raw credential distribution.
9. Passwords/tokens/private keys are never stored in plaintext.
10. Transport security uses TLS where crossing untrusted/network boundaries.
11. SQL is parameterized.
12. Shell commands avoid injection and unsafe interpolation.
13. Untrusted HTML is not rendered without sanitization.
14. Unsafe redirects are rejected.
15. File paths are normalized and constrained.
16. Uploads are size/type/ownership controlled.
17. SSRF-sensitive outbound requests are restricted/validated.
18. Public errors do not leak internals.
19. Security-relevant mutations are audited.
20. Dependency/secret/vulnerability scanning runs in CI/release.
21. Generated artifacts are checksummed; release artifacts should become signed.
22. Least-privilege service accounts are used.
23. Provider credentials are scoped per provider instance/use.
24. Access to production/admin operations requires stronger controls than ordinary reads.
25. Security exceptions require documented owner/reason/expiry.
26. Do not invent cryptography.
27. Use maintained security/protocol libraries.
28. Sensitive configuration is not stored in URLs/query strings.
29. Logs/telemetry are treated as potentially sensitive data stores.
30. Cross-tenant isolation is explicitly tested.

---

## 3. Trust Boundaries

Document important trust transitions:

```text
browser -> public API
service -> service
core -> provider
agent -> control plane
internet -> webhook
user upload -> storage
marketplace package -> installation
```

Each boundary validates identity/input and narrows permissions.

---

## 4. Authentication

Use standards:

```text
OIDC
OAuth 2.x
workload identity
```

Do not create proprietary token schemes for ordinary authentication needs.

---

## 5. Session / Token Handling

Tokens must be:

```text
short-lived where practical
audience-scoped
issuer-validated
expiration-validated
signature-validated
```

Refresh/session material follows the IAM/security design.

---

## 6. Browser Token Storage

Do not persist long-lived sensitive tokens in unsafe browser storage by default.

Prefer secure session mechanisms consistent with the chosen auth architecture.

---

## 7. Workload Identity

Service/agent/provider access should prefer workload identity or short-lived credentials.

Avoid distributing one long-lived shared service secret to many workloads.

---

## 8. Authorization

Authorization evaluates:

```text
principal
resource
action
organization/project/environment context
policy/role
```

Controllers transport the request; business services enforce business authorization.

---

## 9. Permission Naming

Use stable machine-readable permissions.

Pattern:

```text
domain.resource.action
```

Examples:

```text
catalog.application.read
delivery.release.execute
provider.binding.update
```

---

## 10. Least Privilege

Grant only required permissions.

Avoid broad:

```text
admin
*
all
```

when a narrower role/capability exists.

---

## 11. Separation of Duties

Sensitive workflows may require:

```text
requester != approver
deployment operator != policy administrator
```

Where required, enforce server-side and audit.

---

## 12. Just-in-Time Access

High-risk elevated access should support time-bounded/JIT models as enterprise capability evolves.

Do not make permanent admin elevation the default.

---

## 13. Tenant / Organization Isolation

Every tenant-scoped operation derives/validates scope from trusted identity/context.

Never allow an untrusted request header alone to select another organization's data.

---

## 14. IDOR Prevention

Knowing an opaque ID does not grant access.

Every resource lookup that crosses authorization scope verifies ownership/access.

---

## 15. Input Validation

Validate structure at ingress:

```text
type
length
format
allowed values
size
```

Business invariants remain in business services.

---

## 16. Output Encoding

Use framework-safe output encoding.

Never concatenate untrusted strings into HTML/JS/SQL/shell.

---

## 17. XSS

Frontend:

```text
avoid raw HTML
sanitize approved rich HTML
use CSP
treat external markdown/renderers as untrusted
```

---

## 18. CSRF

If browser authentication uses cookies/session credentials, implement appropriate CSRF protection based on cookie/SameSite/origin architecture.

Do not assume SPA = no CSRF risk.

---

## 19. CORS

CORS is explicit and minimal.

Never use wildcard origins with credentials.

Development relaxation must not leak into production.

---

## 20. Redirects

Redirect destinations must be allowlisted/validated.

Avoid open redirects from arbitrary `returnUrl`.

---

## 21. SQL Injection

All values are parameterized.

Dynamic identifiers use a fixed mapping/allowlist.

Never concatenate raw request input into SQL.

---

## 22. Command Injection

Scripts/backend command execution:

```text
prefer direct exec with argument arrays
validate executable/action
avoid shell interpolation
```

Never construct arbitrary shell from untrusted user input.

---

## 23. Path Traversal

Normalize and restrict filesystem access.

Never trust uploaded filename/path.

Use generated storage keys.

---

## 24. SSRF

Outbound URL/provider configuration can become SSRF.

Protect sensitive network zones using:

```text
scheme allowlist
host/network policy
DNS/IP validation where appropriate
egress control
explicit provider endpoints
```

---

## 25. File Uploads

Define:

```text
max size
allowed type
content inspection where needed
storage isolation
authorization
retention
malware scanning responsibility
```

Do not execute uploaded content.

---

## 26. Deserialization

Avoid unsafe polymorphic/deserialization mechanisms that instantiate arbitrary code/types.

Validate untrusted JSON/YAML structures.

---

## 27. YAML

Use safe parsers.

Do not allow arbitrary object construction/tags.

---

## 28. Secrets

Secret values should flow:

```text
Secret Provider
      ↓
short-lived runtime use
      ↓
redacted memory/log handling
```

not:

```text
DB config table
      ↓
API response
      ↓
frontend
```

---

## 29. Secret References

Store references such as:

```text
provider secret ID
vault path/handle
credential profile ID
```

rather than raw secret material.

---

## 30. Secret Rotation

Provider/integration design must tolerate credential rotation without requiring product redeploy where practical.

---

## 31. Encryption

Use established cryptographic libraries.

Keys have:

```text
ownership
rotation
storage
access control
audit
```

Do not hardcode cryptographic secrets.

---

## 32. Randomness

Security tokens/nonces use cryptographically secure randomness.

Do not use pseudo-random utilities intended for simulation.

---

## 33. Passwords

If Perihellion/IAM handles passwords, use a modern password hashing algorithm and policy defined by the IAM security design.

Never log password input.

---

## 34. API Keys

API keys are:

```text
scoped
revocable
rotatable
hashed at rest where verification permits
shown only at creation when possible
```

---

## 35. Webhooks

Verify:

```text
signature/authentication
timestamp/replay protection
source expectations where useful
idempotency
payload size
```

Do not trust source IP alone.

---

## 36. Provider Credentials

Each provider instance has the minimum credential scope needed for declared capabilities.

Do not reuse platform-admin credentials for ordinary provider operations.

---

## 37. Local Development Agent

The agent must use:

```text
device/workload identity
scoped grants
outbound connection where practical
TTL sessions
explicit filesystem/command permissions
```

It must not become an unrestricted remote shell.

---

## 38. Interactive Sessions

Terminal/database/debug sessions require:

```text
authorization
scope
TTL
audit
session identity
resource context
```

No permanent hidden backdoor.

---

## 39. Audit

Audit security-relevant actions.

Minimum safe metadata:

```text
actor/principal
action
target
result
time
org/project/environment
request/correlation ID
```

Do not put secrets in audit payloads.

---

## 40. Logging

Never log:

```text
password
access token
refresh token
private key
secret value
full Authorization header
session cookie
```

Use centralized redaction.

---

## 41. PII / Sensitive Data

Classify sensitive fields.

Logs/metrics/traces should use IDs rather than raw personal data where possible.

Retention/access follows policy.

---

## 42. Error Handling

External errors are safe and normalized.

Detailed internal cause belongs in protected logs/traces.

Do not expose stack traces in production API responses.

---

## 43. Rate Limiting

Apply to abuse-sensitive/public endpoints.

Examples:

```text
authentication
token issuance
password reset
webhooks
expensive search
provider configuration tests
```

Rate limits are scoped and observable.

---

## 44. Denial-of-Service Controls

Use:

```text
request size limits
timeouts
bounded concurrency
bounded queues
connection budgets
pagination
rate limits
```

Avoid unbounded user-controlled work.

---

## 45. Timeouts

All network/database/provider operations have bounded contexts.

Unbounded calls are a reliability and security risk.

---

## 46. Resource Exhaustion

Protect:

```text
DB connections
goroutines
worker concurrency
file descriptors
object upload size
provider clients
browser stream subscriptions
```

---

## 47. Dependency Security

Use automated vulnerability scanning for:

```text
Go modules
npm/pnpm dependencies
container images
OS packages
```

Critical known vulnerabilities require triage.

---

## 48. Lockfiles

Commit package lockfiles.

CI uses reproducible install mode.

Do not silently update lockfiles during unrelated builds.

---

## 49. Supply Chain

Release process should produce:

```text
checksums
SBOM
provenance/signature as maturity increases
immutable OCI digests
```

---

## 50. Generated Code

Generate from trusted source definitions/tools.

Pin generator versions.

Do not execute untrusted code generators from PR input without sandboxing.

---

## 51. CI Secrets

Pull requests from untrusted forks must not receive production/release secrets.

Separate trusted release workflows.

---

## 52. Branch Protection

Main/release branches require protected review/checks according to Git guidelines.

Security-sensitive CODEOWNERS may require specialist approval.

---

## 53. Container Security

Containers:

```text
run non-root
minimal image
no embedded secrets
read-only FS where practical
drop capabilities
resource limits
```

Detailed rules live in Docker/Kubernetes guidelines.

---

## 54. Kubernetes Security

Use:

```text
dedicated service accounts
least RBAC
securityContext
network policy where appropriate
secret references
```

No default cluster-admin.

---

## 55. Database Security

Runtime DB roles use least privileges for owned schemas.

Migration privileges can be separated in production.

TLS/storage encryption follow deployment profile.

---

## 56. Backup Security

Backups are encrypted and access-controlled.

Restore operations are privileged/audited.

Do not place backup credentials in repository scripts.

---

## 57. Marketplace Extensions

Extensions/providers are untrusted until verified.

Eventually require:

```text
publisher identity
signature/digest
declared capabilities
permission review
SDK compatibility
```

---

## 58. AI / Model Providers

Prompts/tools may contain sensitive product data.

AI gateway/provider configuration must define:

```text
data handling
retention
tenant scope
allowed models
redaction
tool permissions
```

Agents never bypass ordinary authorization.

---

## 59. MCP / Tool Calls

Tool definitions have explicit permissions/input schemas.

Tool execution is subject to normal authorization and evidence/audit.

Do not grant "god mode" because the caller is an AI agent.

---

## 60. Frontend Security

Frontend must:

```text
not contain secrets
not trust hidden UI state as authorization
sanitize untrusted rich content
avoid unsafe eval
validate redirects
redact logs
```

---

## 61. Backend Security

Backend must:

```text
derive trusted context
authorize resources/actions
validate inputs
use parameterized SQL
normalize errors
bound resources
audit mutations
```

---

## 62. Script Security

Scripts:

```text
quote variables
verify downloads
avoid curl|bash toolchain bootstrap
never silently sudo
avoid broad rm/pkill
use explicit paths
```

Follow script guideline.

---

## 63. Security Testing

Minimum:

```text
authorization allow/deny
cross-tenant isolation
input boundary tests
secret redaction tests
dependency scan
secret scan
container scan
```

---

## 64. Threat Modeling

Security-sensitive features should document:

```text
assets
actors
trust boundaries
abuse cases
controls
residual risk
```

Keep threat model near design/ADR.

---

## 65. Security Review Triggers

Mandatory review for:

```text
auth/session/token
permission model
secret handling
remote command/terminal
file upload/download
provider credential flow
public webhook
crypto
tenant isolation
marketplace package execution
AI tool permission changes
```

---

## 66. Security Exceptions

Exception record includes:

```text
risk
scope
owner
mitigation
expiry/review date
```

No permanent "temporary" bypass.

---

## 67. Incident Readiness

Security events must be observable and traceable.

Preserve evidence according to policy.

Do not make production debugging require exposing secrets.

---

## 68. Security PR Checklist

```text
[ ] Trust boundary understood
[ ] Authentication source trusted
[ ] Authorization server-side
[ ] Tenant/org scope verified
[ ] Inputs bounded/validated
[ ] No injection path
[ ] Secrets referenced/redacted
[ ] Logs/errors safe
[ ] Timeouts/resource bounds set
[ ] Dependencies scanned
[ ] Audit added where needed
[ ] Allow + deny tests added
[ ] Threat review triggered if required
```

---

## 69. Definition of Done

```text
[ ] No security boundary relies only on UI
[ ] Least privilege applied
[ ] Secrets not persisted/logged casually
[ ] External inputs validated
[ ] Sensitive outputs redacted
[ ] Authz negative tests pass
[ ] Cross-tenant tests pass where relevant
[ ] Security scans pass/triaged
[ ] Audit/observability present
```

---

## 70. Phase 0 Deliverables

Phase 0 establishes:

```text
OIDC/auth integration baseline
principal/context model
RBAC baseline
secret/config handling rule
structured redaction
audit envelope
dependency scan
secret scan
container scan
security headers baseline
cross-tenant/permission negative tests
security review checklist
```

---

## 71. Agent Rules

Before security-sensitive changes read:

```text
docs/AGENT.md
this file
owning backend/frontend guideline
public contract
authorization model
existing tests
```

Never invent a security shortcut to make a test pass.

---

## 72. Final Mandatory Rules

1. Backend authorization is authoritative.
2. Trusted context comes from validated identity.
3. Deny by default.
4. Least privilege everywhere.
5. Secrets are references/redacted whenever possible.
6. Inputs are validated and resource-bounded.
7. No SQL/shell/HTML/path injection.
8. Cross-tenant isolation is tested.
9. Provider/agent credentials are scoped.
10. Security mutations are audited.
11. Dependencies/artifacts are scanned.
12. Security exceptions expire.
13. Agents/AI never bypass normal policy.

---

## 73. Guiding Principle

```text
Trusted Identity
      +
Least Privilege
      +
Validated Boundaries
      +
Secret Minimization
      +
Auditable Actions
      =
Defensible Platform
```

> **Perihellion must remain secure even when UI, providers, extensions, or automation are wrong or compromised at one layer.**
