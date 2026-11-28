# Core Events

Consumers import only `github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/events`.
The root events.go facade exposes Envelope, Metadata, WithCorrelationID,
WithCausationID, NewEnvelope, Rules, Run and ErrIncompleteRules. Generic type
identity, JSON fields, optional IDs and function signatures are preserved.
The former public events/pipeline import is removed; use events instead.

```text
events/
├── events.go
├── README.md
└── internal/
    ├── src/
    │   ├── contracts/
    │   │   ├── envelope.go
    │   │   ├── metadata.go
    │   │   └── rules.go
    │   └── operations/
    │       ├── causation/causation.go
    │       ├── correlation/correlation.go
    │       ├── envelope/envelope.go
    │       └── pipeline/pipeline.go
    └── resources/tests/unit-test/  # exact mirror of src/
```

Correlation precedence: explicit metadata, explicit context, valid trace ID.
Causation precedence: explicit metadata, explicit context. Unknown IDs are
omitted on the wire. Envelope construction does not generate or validate domain
metadata. The pipeline checks all callbacks before executing eligibility,
construction and enqueue in order. Ineligible values do not build or enqueue;
errors pass through unchanged. Message types and delivery policy belong to
consumers. No stateful interface or impl package is necessary for these operations.

```bash
cd backend/libs/core
GOWORK=off go test ./events/... ./tests/unit-test/events/...
GOWORK=off go test -covermode=atomic -coverpkg=./events/... -coverprofile=/tmp/events-coverage.out ./events/internal/resources/tests/unit-test/... ./tests/unit-test/events/...
go tool cover -func=/tmp/events-coverage.out
```
