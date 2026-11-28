# Execution scope

The caller project, target project and tenant are separate concepts.
PLATFORM requires no target project or tenant, PROJECT requires a target project
and no tenant, and TENANT requires both. A project grant never implicitly
authorizes tenant operations. The environment is bound by deployment and
credential audience, not an untrusted request field.

Only trusted authentication adapters may create an execution context using
WithVerifiedExecution after checking the live workload credential, registration,
target ownership, grant, tenant state and any delegated session. This library
validates context shape and exact grant matching; it is not an authenticator.
System contexts are reserved for trusted in-process jobs and are not workload
grants. No HTTP headers or project-specific values are interpreted here.

Existing session and gateway callers are unchanged until the Registry and
gateway cutover is complete. No legacy authorization fallback should remain
after that cutover.
