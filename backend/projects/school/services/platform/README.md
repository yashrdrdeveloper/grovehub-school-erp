# School business API wrappers

Explicit operation registry for Identity accounts, account roles, capabilities, roles, approvals and Audit search. The School checks browser access credential, Origin and CSRF; Platform verifies the credential against live Identity state and applies its existing authorization policy. Request and response JSON are transitional wire payloads, not generated typed DTOs. No arbitrary forwarding or query-string forwarding is supported. Configure Platform workload grants with RequireActor=true for these routes. Existing cutover restrictions remain enforced.
