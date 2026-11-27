# Platform workload configuration

Set PLATFORM_WORKLOAD_CONFIG to a JSON object with Keys and Routes. Each key has KeyID, PublicKey (unpadded base64url Ed25519), Issuer, Subject, Scopes and AllowDelegation. Routes maps exact method + restricted route template to Scope and RequireActor. Audience is platform. Empty configuration denies workload calls except the separately verified Gmail webhook.

The provider-owned restrictedPolicies table is authoritative: configured routes must match its exact scope and actor requirement. A configuration may omit a route to disable it. One signed scope is accepted per request. Login and forgot-password use only the school workload; refresh, recovery and setup use route-specific proofs validated by their controllers. Other Identity operations and Audit require live access-proof verification and protected-route authorization. Integration business operations require live access-proof verification; controllers enforce tenant, ownership and system-mail admin rules. OAuth callback uses the school workload and one-time state, not a browser credential.

Foundation uses the separate school-foundation subject and explicit foundation.persons.create/get/update/search/deactivate scopes; no School browser wrapper or ordinary school workload is allowed. Gmail push cannot be configured as a school workload operation.

School constructs separate Identity, Audit and Integration HTTP clients, each restricted to its own operation policy. For development these clients share the existing School signing identity and key; this is client-side policy isolation, not independent cryptographic identities. Keep private keys in consumers. Existing environment configuration must use the canonical scopes in the table. Do not deploy until full application verification passes.

## Internal Tool catalog access

Optional InternalToolKeys in PLATFORM_WORKLOAD_CONFIG is a separate list of public keys. Each must have Issuer and Subject equal to internal-tool, AllowDelegation=false, and exactly one scope internal.integration.catalog.read. Key IDs and public keys may not overlap School keys. It authorizes only GET /integrations/restricted/providers, without cookies, identity proof, actor claims or query parameters. The existing School route policy still requires live actor authorization. No tenant-specific operations or mutations are enabled. Configure keys separately for each Platform environment.
