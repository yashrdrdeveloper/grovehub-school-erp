# GroveHub Backend Unit Test Structure

> **Status:** Mandatory GroveHub engineering standard  
> **Applies to:** backend Go unit tests for launcher, services, shared libraries, and providers  
> **Framework:** Go `testing` + Ginkgo v2 + Gomega  
> **Precedence:** This file is authoritative for backend unit-test placement when an older generic guideline conflicts with it.

---

## 1. Ownership Rule

Backend unit tests live with the backend owner they validate.

Do not use the old centralized structure:

```text
backend/testing/unit-test/...
```

Canonical ownership is:

```text
backend/launcher/tests/unit-test/...
backend/services/<service>/tests/unit-test/...
backend/libs/core/tests/unit-test/...
backend/providers/<provider>/tests/unit-test/...
```

Create only the owner paths that actually have tests.

Do not create tests merely to assert that plain data-only structs contain fields. Pure data-only entities do not require meaningless unit tests.

---

## 2. Mirroring Rule

Inside an owner's `tests/unit-test/`, mirror the flat production responsibility.

Examples:

```text
backend/launcher/routing/
  ↓
backend/launcher/tests/unit-test/routing/

backend/launcher/constants/errorcode/
  ↓
backend/launcher/tests/unit-test/constants/errorcode/

backend/libs/core/http/response/
  ↓
backend/libs/core/tests/unit-test/http/response/

backend/libs/core/id/
  ↓
backend/libs/core/tests/unit-test/id/

backend/services/identity/controller/
  ↓
backend/services/identity/tests/unit-test/controller/

backend/services/foundation/services/
  ↓
backend/services/foundation/tests/unit-test/services/

backend/services/foundation/repository/postgres/
  ↓
backend/services/foundation/tests/unit-test/repository/postgres/
```

There is no `src/` segment in the current GroveHub service tree or its mirrored tests.

Do not flatten unrelated responsibilities into one test package.

---

## 3. What to Test by Layer

Follow the GroveHub layer responsibilities from `docs/engineering/backend/GROVEHUB_BACKEND_STRUCTURE.md`.

### Entity

`entity/` is data-only.

Do not add unit tests for plain fields, aliases, or trivial constants unless there is a real compatibility contract that needs protection.

If behavior appears in an entity and requires business-rule tests, first verify that the behavior should actually live in `services/`.

### Services

Business rules and use cases belong in `services/`, so this is where focused business unit tests normally live.

Test, as applicable:

```text
validation decisions
business defaults
status/lifecycle transitions
ID-generation orchestration
version/timestamp decisions
authorization decisions
repository/client orchestration
business error behavior
```

Use fake/stub repository contracts and controlled technical dependencies.

### Repository

Unit-test pure repository mapping/error translation only when it can be proven without real infrastructure.

Do not unit-test business decisions in repositories because repositories must not own business decisions.

Real PostgreSQL SQL/constraints/transactions/locking/index-dependent semantics belong to repository-wide sanity tests.

### Controller

Test transport behavior only:

```text
request mapping
service invocation
response mapping
HTTP status/error mapping
headers
```

Do not duplicate service business-rule tests in controllers.

---

## 4. Canonical Framework

Backend unit tests use:

```go
import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)
```

The ordinary Go `testing` package is used by `suite_test.go` to register the Ginkgo suite.

Do not add a second backend unit-test framework.

---

## 5. `suite_test.go` Is Mandatory

Every backend unit-test package folder that exists must contain:

```text
suite_test.go
```

Example:

```text
backend/services/identity/tests/unit-test/controller/
├── suite_test.go
└── health_test.go
```

`suite_test.go` owns suite-level setup and reusable test support for that folder.

It contains, as appropriate:

```text
RunSpecs registration
shared constants
shared test types
fixtures
builders
fake/stub definitions
helper functions
common setup/cleanup
controlled clocks/IDs
shared test tracer/provider setup
```

Behavior files should remain focused on the test cases.

---

## 6. Behavior Test Files

Files other than `suite_test.go` should primarily contain Ginkgo behavior declarations:

```go
var _ = Describe("...", func() {
    Context("...", func() {
        It("...", func() {
            // arrange through suite helpers
            // act
            // assertions
        })
    })
})
```

Do not repeatedly redefine the same helper, fixture, constant, fake, builder, or setup function across behavior files in the same package.

Promote shared test support to `suite_test.go` first.

If `suite_test.go` becomes too large, split coherent support into explicit files such as:

```text
fixtures_test.go
builders_test.go
fakes_test.go
assertions_test.go
```

Do not create generic `utils_test.go` or `helpers_test.go` dumping grounds without a clear responsibility.

---

## 7. Example Suite

```go
package controller_test

import (
    "testing"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

const healthPath = "/identity/health"

func TestIdentityControllerUnit(t *testing.T) {
    RegisterFailHandler(Fail)
    RunSpecs(t, "Identity Controller Unit Tests")
}
```

Behavior file:

```go
package controller_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
)

var _ = Describe("Health controller", func() {
    It("returns a healthy response", func() {
        // focused transport behavior only
    })
})
```

---

## 8. Package Boundaries

Prefer external test packages when practical:

```text
package services_test
package repository_test
package controller_test
package response_test
```

Use the production package directly only when testing intentionally private behavior is unavoidable and justified.

Tests should validate public/owned behavior rather than lock implementation details unnecessarily.

---

## 9. Unit-Test Constraints

Unit tests must be:

```text
deterministic
isolated
fast
network-independent
Docker-independent
real-clock independent when time matters
order-independent
```

Unit tests do not require a live Grafana, PostgreSQL, Redis, SaaS provider, or public network.

Real infrastructure semantics belong to the repository-wide `testing/sanity/` suite.

---

## 10. Repository-Wide Suites

Only backend unit tests live in owner-local `tests/unit-test/` folders.

Repository-wide suites remain:

```text
testing/sanity/
testing/regression/
```

Do not move product-wide sanity/regression into a single service.

Do not create new top-level `integration`, `smoke`, `critical`, or `contract` suites.

---

## 11. New Backend Owner Checklist

When adding a new backend owner or responsibility:

```text
[ ] Identify owning production path
[ ] Confirm the code is in the correct layer first
[ ] Create tests/unit-test only when meaningful unit tests exist
[ ] Mirror the production responsibility under tests/unit-test
[ ] Do not add entity tests for plain data-only structs
[ ] Add suite_test.go in every test package folder that exists
[ ] Use Ginkgo v2 + Gomega
[ ] Put shared constants/functions/setup in suite_test.go
[ ] Keep behavior files focused on specs
[ ] Do not recreate backend/testing/unit-test
[ ] Do not introduce a src/ segment
```

---

## 12. Agent Rule

Before creating or moving a backend test, determine its production owner and layer first.

Routing examples:

```text
launcher behavior
  -> backend/launcher/tests/unit-test/...

foundation business behavior
  -> backend/services/foundation/tests/unit-test/services/...

foundation PostgreSQL mapping behavior
  -> backend/services/foundation/tests/unit-test/repository/postgres/...

identity controller behavior
  -> backend/services/identity/tests/unit-test/controller/...

core library behavior
  -> backend/libs/core/tests/unit-test/...
```

If an older document shows `backend/testing/unit-test`, a `src/` test path, or business behavior tested in `entity/`, this newer GroveHub standard wins.
