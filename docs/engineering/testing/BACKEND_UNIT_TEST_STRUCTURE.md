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

---

## 2. Mirroring Rule

Inside an owner's `tests/unit-test/`, mirror the production responsibility.

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

backend/services/identity/controller/
  ↓
backend/services/identity/tests/unit-test/controller/

backend/services/identity/constants/errorcode/
  ↓
backend/services/identity/tests/unit-test/constants/errorcode/
```

Do not flatten unrelated responsibilities into one test package.

---

## 3. Canonical Framework

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

## 4. `suite_test.go` Is Mandatory

Every backend unit-test package folder must contain:

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

## 5. Behavior Test Files

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

## 6. Example Suite

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
        // focused behavior only
    })
})
```

---

## 7. Package Boundaries

Prefer external test packages when practical:

```text
package routing_test
package controller_test
package response_test
```

Use the production package directly only when testing intentionally private behavior is unavoidable and justified.

Tests should validate public/owned behavior rather than lock implementation details unnecessarily.

---

## 8. Unit-Test Constraints

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

## 9. Repository-Wide Suites

Only backend unit tests move to owner-local `tests/unit-test/` folders.

Repository-wide suites remain:

```text
testing/sanity/
testing/regression/
```

Do not move product-wide sanity/regression into a single service.

Do not create new top-level `integration`, `smoke`, `critical`, or `contract` suites.

---

## 10. New Backend Owner Checklist

When adding a new backend owner or responsibility:

```text
[ ] Identify owning production path
[ ] Create tests/unit-test only when unit tests exist
[ ] Mirror the production responsibility under tests/unit-test
[ ] Add suite_test.go in every test package folder
[ ] Use Ginkgo v2 + Gomega
[ ] Put shared constants/functions/setup in suite_test.go
[ ] Keep behavior files focused on specs
[ ] Do not recreate backend/testing/unit-test
```

---

## 11. Agent Rule

Before creating or moving a backend test, determine its production owner first.

Routing examples:

```text
launcher behavior
  -> backend/launcher/tests/unit-test/...

identity behavior
  -> backend/services/identity/tests/unit-test/...

core library behavior
  -> backend/libs/core/tests/unit-test/...
```

If an older document shows `backend/testing/unit-test`, this newer GroveHub standard wins.
