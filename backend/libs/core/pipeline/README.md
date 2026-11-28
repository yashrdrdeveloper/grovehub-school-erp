# Core Pipeline

The root `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline` facade exposes typed REST execution and
service orchestration. Internal operations preserve the existing behavior.

- REST: RunREST, RunStruct, ReadBody, Config, StructConfig, Writer, Result, Input,
  Failure, OK, Created and CreatedOrOK.
- Service: RunService, Chain, Next, Return, Fail, Decide, ResolveResult,
  Lookup, Sequence, Transactional and Recover, with their typed contracts.

The former rest.Run maps to RunREST, service.Run maps to RunService, and
service.Result maps to ResolveResult to avoid root name collisions. The former
public rest and service packages are temporary forwarding wrappers until
all consumers are migrated. Business-specific policies remain in services.

```text
pipeline/
├── pipeline.go
└── internal/
    ├── src/operations/{rest,service}/
    └── resources/tests/unit-test/operations/{rest,service}/
```

```bash
cd backend/libs/core
GOWORK=off go test ./pipeline/... ./tests/unit-test/pipeline/...
```
