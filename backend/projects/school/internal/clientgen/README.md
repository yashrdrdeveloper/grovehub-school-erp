# School generated Platform clients

Audit is the first runtime-integrated generated Go client. Its consumer-local
specification is pinned under specs/audit.yaml from Platform commit c22512f14b144eb02617e57768ba9227fd25ce05.
Run make generate-clients (Docker required) before direct go commands. The
School Makefile and independent backend CI generate it automatically. Generated
files live under internal/generated/audit and are excluded from Git; the pinned
specification and generator version determine their content.

Audit uses the generated SearchRequest and SearchAuditActivities operation
through the existing origin-bound signed workload transport. Response size
limits and original JSON preservation remain enforced at the application
boundary. Identity, Integration and Foundation runtime SDK integration is not
yet performed. Updating the Audit contract requires an explicit snapshot
refresh and contract regression tests. The local Core generator dependency
remains until independent library publication.
