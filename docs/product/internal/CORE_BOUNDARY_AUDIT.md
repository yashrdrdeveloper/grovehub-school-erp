# Backend Core extraction audit

Inventory from [isolated generated-overlay build and tests](https://github.com/yashrdrdeveloper/grovehub-school-erp/actions/runs/45828175362), candidate `688a1ca2d0b74db6ce074255a253ddf1c1937eb4`, transferred to `049df4e8e6e68dc36350c0220ecb7b42a5814a8c`. **104 Go packages**, including generated test package imports, zero package listing errors and zero school/launcher imports after removing the product DTO test dependency and school session contracts. This proves import-boundary separation, not complete semantic neutrality or a standalone build. Core is still physically under `backend/libs/core`.

## Disposition

- Neutral mechanics: retain candidate runtime packages, review public contracts and telemetry defaults; move together after standalone module tests.
- Tooling: annotation processors and testing/overlay/coverage commands remain conditional. Generated service import strings use the target go.mod; Core import prefix is centralized. Explicit application/core roots, overlay ownership, test inventory/coverage and shell build scripts still need migration.
- Test infrastructure: preserve Core fixtures and unit tests without product imports. Original Foundation DTO tag tests now live in Foundation.
- School adapter: `backend/libs/schoolgateway` remains in the school module, including the legacy admin flag, PersonGID and existing browser cookie names. No privilege or cookie behavior was changed. Core retains only generic route tiers/workload verification.

## Complete direct import inventory

Includes production, test and external-test imports from go list with the generated overlay. Package paths below are relative to `grovehub/libs/core`; imports are full Go paths. Stdlib imports are included. Classifications describe extraction intent, not final approval to move.

| Package | Classification | Direct imports |
|---|---|---|
| `annotation` | Runtime candidate | `errors`, `fmt`, `grovehub/libs/core/annotation/chain`, `grovehub/libs/core/annotation/composition`, `grovehub/libs/core/annotation/configuration`, `grovehub/libs/core/annotation/contract`, `grovehub/libs/core/annotation/dto`, `grovehub/libs/core/annotation/entity`, `grovehub/libs/core/annotation/event`, `grovehub/libs/core/annotation/external`, `grovehub/libs/core/annotation/http`, `grovehub/libs/core/annotation/mapping`, `grovehub/libs/core/annotation/messaging`, `grovehub/libs/core/annotation/middleware`, `grovehub/libs/core/annotation/observability`, `grovehub/libs/core/annotation/patterns`, `grovehub/libs/core/annotation/persistence`, `grovehub/libs/core/annotation/security`, `grovehub/libs/core/annotation/server`, `grovehub/libs/core/annotation/strategy`, `grovehub/libs/core/annotation/testing`, `grovehub/libs/core/annotation/utilities`, `grovehub/libs/core/annotation/validation`, `sort`, `strings`, `unicode` |
| `annotation/chain` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/composition` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/configuration` | Runtime candidate | `fmt`, `grovehub/libs/core/annotation/contract`, `os`, `reflect`, `strconv`, `strings`, `time` |
| `annotation/contract` | Runtime candidate | None |
| `annotation/dto` | Runtime candidate | `bytes`, `encoding/json`, `grovehub/libs/core/annotation/contract` |
| `annotation/entity` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/event` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/external` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/generation/check` | Tooling — verify roots | `fmt`, `go/ast`, `go/parser`, `go/token`, `grovehub/libs/core/annotation`, `grovehub/libs/core/annotation/generation/parser`, `io/fs`, `os`, `path/filepath`, `reflect`, `strconv`, `strings`, `testing` |
| `annotation/generation/clients` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `grovehub/libs/core/annotation/generation/layout`, `os`, `path/filepath`, `sort`, `strconv`, `strings`, `testing` |
| `annotation/generation/entities` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `os`, `path`, `path/filepath`, `sort`, `strconv`, `strings` |
| `annotation/generation/events` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `io/fs`, `os`, `path/filepath`, `regexp`, `sort`, `strconv`, `strings` |
| `annotation/generation/expand` | Tooling — verify roots | `fmt`, `go/ast`, `go/parser`, `go/token`, `os`, `path/filepath`, `strings` |
| `annotation/generation/gormtable` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `os`, `path/filepath`, `regexp`, `sort`, `strings` |
| `annotation/generation/layout` | Tooling — verify roots | `fmt`, `golang.org/x/mod/modfile`, `golang.org/x/mod/module`, `grovehub/libs/core/annotation/generation/layout`, `os`, `path/filepath`, `strings`, `testing` |
| `annotation/generation/mappers` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `grovehub/libs/core/annotation/generation/layout`, `os`, `path`, `path/filepath`, `sort`, `strconv`, `strings`, `testing` |
| `annotation/generation/overlay` | Tooling — verify roots | `bytes`, `encoding/json`, `fmt`, `io/fs`, `os`, `path/filepath`, `strings` |
| `annotation/generation/parser` | Tooling — verify roots | `fmt`, `go/ast`, `go/parser`, `go/token`, `grovehub/libs/core/annotation` |
| `annotation/generation/patchdto` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `os`, `path`, `path/filepath`, `reflect`, `sort`, `strconv`, `strings`, `testing` |
| `annotation/generation/routes` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `grovehub/libs/core/annotation/generation/layout`, `os`, `path/filepath`, `sort`, `strconv`, `strings`, `testing` |
| `annotation/generation/search` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `grovehub/libs/core/annotation/generation/layout`, `os`, `path/filepath`, `sort`, `strconv`, `strings`, `testing` |
| `annotation/generation/testing` | Tooling — verify roots | `bytes`, `fmt`, `go/ast`, `go/format`, `go/parser`, `go/token`, `grovehub/libs/core/annotation`, `io/fs`, `os`, `os/exec`, `path/filepath`, `sort`, `strings`, `unicode` |
| `annotation/generation/testing/coverage` | Tooling — verify roots | `bufio`, `context`, `encoding/json`, `errors`, `flag`, `fmt`, `github.com/onsi/gomega`, `io`, `os`, `os/exec`, `path/filepath`, `sort`, `strconv`, `strings`, `testing` |
| `annotation/generation/testing/inventory` | Tooling — verify roots | `encoding/json`, `fmt`, `go/ast`, `go/parser`, `go/token`, `io/fs`, `os`, `path/filepath`, `sort`, `strings` |
| `annotation/http` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/mapping` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/messaging` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/middleware` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/observability` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/patterns` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/persistence` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/security` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/server` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/strategy` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/testing` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/utilities` | Runtime candidate | `grovehub/libs/core/annotation/contract` |
| `annotation/validation` | Runtime candidate | `fmt`, `grovehub/libs/core/annotation/contract`, `grovehub/libs/core/annotation/dto`, `grovehub/libs/core/normalize`, `reflect`, `strconv`, `strings`, `unicode/utf8` |
| `database/postgres` | Runtime candidate | `gorm.io/driver/postgres`, `gorm.io/gorm` |
| `database/transaction` | Runtime candidate | `context`, `errors`, `gorm.io/gorm` |
| `diff` | Runtime candidate | `bytes`, `encoding/json`, `fmt`, `reflect`, `sort`, `strconv` |
| `errors` | Runtime candidate | `grovehub/libs/core/errors/contracts`, `grovehub/libs/core/errors/helpers` |
| `errors/contracts` | Runtime candidate | `net/http` |
| `errors/helpers` | Runtime candidate | `grovehub/libs/core/errors/contracts` |
| `events` | Runtime candidate | `context`, `grovehub/libs/core/telemetry`, `time` |
| `events/pipeline` | Runtime candidate | `context`, `errors`, `testing` |
| `gateway` | Runtime candidate | `context`, `net/http` |
| `headercontext` | Runtime candidate | `context`, `errors`, `grovehub/libs/core/events`, `grovehub/libs/core/http/headers`, `grovehub/libs/core/idempotency`, `net/http`, `strings` |
| `helpers` | Runtime candidate | `testing` |
| `http/contracts` | Runtime candidate | `grovehub/libs/core/errors` |
| `http/endpoint` | Runtime candidate | `bytes`, `context`, `encoding/json`, `errors`, `fmt`, `grovehub/libs/core/errors`, `grovehub/libs/core/headercontext`, `grovehub/libs/core/http/headers`, `grovehub/libs/core/request`, `io`, `net/http`, `net/url`, `reflect`, `regexp`, `strconv`, `strings` |
| `http/headers` | Runtime candidate | None |
| `http/logging` | Runtime candidate | `log/slog`, `net/http` |
| `http/request` | Runtime candidate | `encoding/json`, `errors`, `io`, `net/http` |
| `http/response` | Runtime candidate | `context`, `encoding/json`, `go.opentelemetry.io/otel/attribute`, `go.opentelemetry.io/otel/codes`, `go.opentelemetry.io/otel/trace`, `grovehub/libs/core/errors`, `grovehub/libs/core/http/contracts`, `net/http` |
| `http/tracing` | Runtime candidate | `go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp`, `go.opentelemetry.io/otel/attribute`, `go.opentelemetry.io/otel/codes`, `go.opentelemetry.io/otel/trace`, `grovehub/libs/core/http/headers`, `grovehub/libs/core/telemetry`, `net/http`, `strings`, `time` |
| `id` | Runtime candidate | `grovehub/libs/core/id/helpers` |
| `id/helpers` | Runtime candidate | `github.com/google/uuid` |
| `idempotency` | Runtime candidate | `grovehub/libs/core/idempotency/contracts`, `grovehub/libs/core/idempotency/helpers`, `grovehub/libs/core/idempotency/validators` |
| `idempotency/contracts` | Runtime candidate | `errors` |
| `idempotency/helpers` | Runtime candidate | `crypto/sha256`, `encoding/hex`, `encoding/json`, `fmt`, `grovehub/libs/core/idempotency/contracts` |
| `idempotency/validators` | Runtime candidate | `grovehub/libs/core/idempotency/contracts`, `strings` |
| `logging` | Runtime candidate | `grovehub/libs/core/logging/contracts`, `grovehub/libs/core/logging/helpers`, `log/slog` |
| `logging/contracts` | Runtime candidate | `log/slog` |
| `logging/helpers` | Runtime candidate | `context`, `fmt`, `grovehub/libs/core/logging/contracts`, `log/slog`, `strings` |
| `module` | Runtime candidate | `net/http` |
| `normalize` | Runtime candidate | `grovehub/libs/core/normalize/helpers` |
| `normalize/helpers` | Runtime candidate | `errors`, `github.com/google/uuid`, `net/mail`, `regexp`, `strings`, `time` |
| `pipeline/rest` | Runtime candidate | `context`, `errors`, `fmt`, `grovehub/libs/core/annotation/validation`, `grovehub/libs/core/errors`, `grovehub/libs/core/http/request`, `grovehub/libs/core/http/response`, `grovehub/libs/core/request`, `net/http`, `net/url`, `reflect`, `strconv`, `strings` |
| `pipeline/service` | Runtime candidate | `context`, `errors`, `grovehub/libs/core/database/transaction`, `testing` |
| `request` | Runtime candidate | `fmt`, `reflect` |
| `search` | Runtime candidate | `fmt`, `gorm.io/gorm`, `grovehub/libs/core/search/contracts`, `grovehub/libs/core/search/helpers`, `grovehub/libs/core/search/mappers`, `grovehub/libs/core/search/validators`, `math`, `strings`, `testing` |
| `search/contracts` | Runtime candidate | None |
| `search/helpers` | Runtime candidate | `errors`, `fmt`, `grovehub/libs/core/search/contracts`, `strings` |
| `search/mappers` | Runtime candidate | `encoding/base64`, `encoding/json`, `grovehub/libs/core/search/contracts` |
| `search/validators` | Runtime candidate | `errors`, `grovehub/libs/core/normalize`, `grovehub/libs/core/search/contracts`, `grovehub/libs/core/search/mappers`, `strings`, `time` |
| `telemetry` | Runtime candidate | `context`, `go.opentelemetry.io/otel/metric`, `go.opentelemetry.io/otel/trace`, `grovehub/libs/core/telemetry/contracts`, `grovehub/libs/core/telemetry/helpers`, `grovehub/libs/core/telemetry/mappers`, `log/slog`, `time` |
| `telemetry/contracts` | Runtime candidate | `context` |
| `telemetry/helpers` | Runtime candidate | `context`, `errors`, `fmt`, `go.opentelemetry.io/contrib/instrumentation/runtime`, `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/attribute`, `go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp`, `go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp`, `go.opentelemetry.io/otel/metric`, `go.opentelemetry.io/otel/propagation`, `go.opentelemetry.io/otel/sdk/metric`, `go.opentelemetry.io/otel/sdk/resource`, `go.opentelemetry.io/otel/sdk/trace`, `go.opentelemetry.io/otel/trace`, `grovehub/libs/core/telemetry/contracts`, `grovehub/libs/core/telemetry/mappers`, `grovehub/libs/core/telemetry/validators`, `log/slog`, `net/http`, `os`, `strings`, `sync`, `time` |
| `telemetry/mappers` | Runtime candidate | `grovehub/libs/core/telemetry/contracts`, `strings` |
| `telemetry/validators` | Runtime candidate | `fmt`, `strconv`, `strings`, `time` |
| `tests/helpers/clock` | Test infrastructure | `time` |
| `tests/helpers/http` | Test infrastructure | `bytes`, `net/http`, `net/http/httptest` |
| `tests/mocks/controller` | Test infrastructure | `go.uber.org/mock/gomock` |
| `tests/mocks/database` | Test infrastructure | `database/sql`, `github.com/DATA-DOG/go-sqlmock`, `gorm.io/driver/postgres`, `gorm.io/gorm` |
| `tests/unit-test/annotation/dto` | Test infrastructure | `encoding/json`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/annotation/dto`, `testing` |
| `tests/unit-test/annotation/validation` | Test infrastructure | `encoding/json`, `errors`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/annotation/dto`, `grovehub/libs/core/annotation/validation`, `strings`, `testing` |
| `tests/unit-test/database/transaction` | Test infrastructure | `context`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/database/transaction`, `testing` |
| `tests/unit-test/diff` | Test infrastructure | `encoding/json`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/diff`, `testing` |
| `tests/unit-test/entity` | Test infrastructure | `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/annotation/entity`, `testing` |
| `tests/unit-test/errors` | Test infrastructure | `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/errors`, `net/http`, `testing` |
| `tests/unit-test/events` | Test infrastructure | `context`, `encoding/json`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `go.opentelemetry.io/otel/trace`, `grovehub/libs/core/events`, `testing`, `time` |
| `tests/unit-test/headercontext` | Test infrastructure | `context`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/events`, `grovehub/libs/core/headercontext`, `net/http`, `net/http/httptest`, `testing` |
| `tests/unit-test/http/endpoint` | Test infrastructure | `context`, `encoding/json`, `errors`, `fmt`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/errors`, `grovehub/libs/core/headercontext`, `grovehub/libs/core/http/endpoint`, `grovehub/libs/core/http/response`, `grovehub/libs/core/pipeline/rest`, `grovehub/libs/core/request`, `net/http`, `strings`, `testing` |
| `tests/unit-test/http/logging` | Test infrastructure | `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/http/logging`, `log/slog`, `net/http`, `testing` |
| `tests/unit-test/http/response` | Test infrastructure | `context`, `encoding/json`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/attribute`, `go.opentelemetry.io/otel/codes`, `go.opentelemetry.io/otel/sdk/trace`, `go.opentelemetry.io/otel/sdk/trace/tracetest`, `go.opentelemetry.io/otel/trace`, `grovehub/libs/core/errors`, `grovehub/libs/core/http/contracts`, `grovehub/libs/core/http/response`, `grovehub/libs/core/telemetry`, `net/http`, `net/http/httptest`, `testing` |
| `tests/unit-test/http/tracing` | Test infrastructure | `context`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/propagation`, `go.opentelemetry.io/otel/sdk/trace`, `go.opentelemetry.io/otel/sdk/trace/tracetest`, `go.opentelemetry.io/otel/trace`, `grovehub/libs/core/http/headers`, `grovehub/libs/core/http/tracing`, `grovehub/libs/core/telemetry`, `io`, `net/http`, `net/http/httptest`, `strings`, `testing` |
| `tests/unit-test/id` | Test infrastructure | `github.com/google/uuid`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/id`, `testing` |
| `tests/unit-test/idempotency` | Test infrastructure | `grovehub/libs/core/idempotency`, `testing` |
| `tests/unit-test/logging` | Test infrastructure | `bytes`, `context`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/logging`, `log/slog`, `testing` |
| `tests/unit-test/normalize` | Test infrastructure | `errors`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/normalize`, `testing` |
| `tests/unit-test/pipeline/rest` | Test infrastructure | `context`, `errors`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `grovehub/libs/core/annotation/validation`, `grovehub/libs/core/errors`, `grovehub/libs/core/http/request`, `grovehub/libs/core/pipeline/rest`, `grovehub/libs/core/request`, `net/http`, `net/http/httptest`, `strings`, `testing` |
| `tests/unit-test/request` | Test infrastructure | `grovehub/libs/core/request`, `testing` |
| `tests/unit-test/telemetry` | Test infrastructure | `bytes`, `context`, `encoding/json`, `github.com/onsi/ginkgo/v2`, `github.com/onsi/gomega`, `go.opentelemetry.io/otel`, `go.opentelemetry.io/otel/sdk/trace`, `go.opentelemetry.io/otel/trace`, `grovehub/libs/core/telemetry`, `log/slog`, `testing` |

## Remaining gates

Select canonical Go module identity, generalize generation roots and script invocations, validate package ownership and generated overlays across both app targets, migrate Go imports atomically, generate standalone go.mod/go.sum and go.work, verify GOWORK=off build/release dependency resolution and all school tests. No school role/permission policy belongs in Core. See the [internal plan](INTERNAL_PLATFORM_IMPLEMENTATION_PLAN.md).
