# GroveHub Infrastructure Service Guidelines

> **Status:** Mandatory GroveHub backend infrastructure ownership standard  
> **Applies to:** integrations with SMTP/mail engines, cache engines, object storage, queues, search engines, external infrastructure APIs, and similar technical infrastructure capabilities  
> **Owner:** `backend/services/infrastructure`

---

## 1. Single Infrastructure Ownership Boundary

All GroveHub **backend integration code for infrastructure capabilities** lives under:

```text
backend/services/infrastructure/
```

Business services must not create provider-specific infrastructure adapters inside their own service roots.

Examples:

```text
GOOD
backend/services/infrastructure/services/impl/smtp.go
backend/services/infrastructure/services/impl/dkim.go
backend/services/infrastructure/repository/impl/dkim.go

BAD
backend/services/mailProvider/transport/impl/haraka.go
backend/services/identity/redis.go
backend/services/foundation/s3.go
```

The business service owns the business use case. Infrastructure owns the technical integration.

Runtime/deployment assets still follow repository-wide deployment ownership. Dockerfiles, Compose files, deployment manifests, and developer shell orchestration may remain under `deployments/` and `scripts/`; the Go integration and infrastructure state boundary belongs under `backend/services/infrastructure`.

---

## 2. Canonical Infrastructure Structure

Infrastructure follows the same GroveHub service structure as every other backend module:

```text
backend/services/infrastructure/
├── infrastructure.go
├── client/
│   ├── <capability>.go
│   └── impl/
│       └── <capability>.go
├── config/
├── constants/
│   └── error/
├── dto/
│   ├── request/
│   └── response/
├── entity/
├── repository/
│   ├── <resource>.go
│   └── impl/
│       └── <resource>.go
├── services/
│   ├── <capability>.go
│   ├── impl/
│   │   └── <capability>.go
│   └── helpers/
│       └── <capability>.go
├── controller/
│   ├── <resource>.go
│   ├── impl/
│   │   └── <resource>.go
│   └── helpers/
├── routes/
├── migrations/
└── tests/
```

Create only folders that have a current responsibility.

---

## 3. Client Rule

Consumers depend only on the public Infrastructure client contract.

```text
Consuming service
    ↓
infrastructure/client
    ↓
infrastructure/client/impl
    ↓
infrastructure/services interface
    ↓
infrastructure/services/impl
```

For mail:

```text
MailProvider
    ↓
client.MailClient
    ↓
client/impl.MailClient
    ↓
services.SMTPService
    ↓
services/impl.SMTPService
    ↓
SMTP endpoint / Haraka
```

The client implementation is intentionally thin. It forwards `context.Context` and delegates to the owning service. It must not duplicate SMTP behavior, retry policy, DKIM logic, or provider-specific configuration.

Business services must not import:

```text
grovehub/services/infrastructure/services/impl
grovehub/services/infrastructure/repository/impl
```

They import the public `grovehub/services/infrastructure/client` package.

---

## 4. Service and Helper Rule

`services/<capability>.go` defines the service contract.

`services/impl/<capability>.go` owns technical use-case orchestration and external I/O required by the infrastructure capability.

`services/helpers/<capability>.go` contains **pure capability-specific helpers only**.

Good helper examples:

```text
normalize DKIM domain/selector
parse PEM data
format Message-ID headers
format References headers
convert a PEM public key to a DNS value
```

Do not put these in `services/helpers/`:

```text
database access
environment reads
network calls
SMTP connections
filesystem mutation
Docker execution
HTTP controller logic
```

Environment/configuration loading belongs under `config/`. Persistence belongs under `repository/impl/`. HTTP transport belongs under `controller/` and `routes/`.

---

## 5. Provider-Agnostic Naming

Public contracts use capability names, not implementation/vendor names.

Use:

```text
MailClient
SMTPService
StorageClient
CacheClient
QueueClient
```

Avoid public names such as:

```text
HarakaClient
RedisClient
S3Client
KafkaClient
```

Provider-specific names may exist only inside a concrete implementation when the distinction is genuinely required. A business service must not know whether SMTP is currently delivered through Haraka, Postfix, SES, or another implementation.

---

## 6. DKIM State Ownership

DKIM identity is infrastructure state, not MailProvider message/domain business state.

The durable source of truth is the Infrastructure database schema:

```text
infrastructure.dkim_keys
```

The persisted identity includes the domain, selector, private key material, public key material, and creation time.

Current rule:

```text
first setup
    -> import an existing valid filesystem key when present
    -> otherwise generate once
    -> persist once in infrastructure.dkim_keys

later startup/recreate/rebuild
    -> read the same persisted DB identity
    -> materialize the same key for the SMTP/DKIM runtime
    -> never silently generate a replacement key
```

DKIM rows are immutable. Rotation must be an explicit future workflow; normal startup/domain preparation must never rotate a DKIM identity.

The private key must never be returned through the Infrastructure HTTP/API response. Only the public key required for DNS publication may be returned.

A future secret-provider/KMS implementation may replace direct private-key storage with a secret reference without changing the public Infrastructure client/service ownership model.

---

## 7. HTTP/Internal Operations

If infrastructure exposes an internal HTTP operation, it still follows the normal GroveHub flow:

```text
Route
  ↓
Controller
  ↓
Service
  ↓
Repository / external infrastructure
```

Do not put HTTP handlers directly in `infrastructure.go`.

The module file is the composition root only: construct repositories, services, clients, controllers, and register routes.

---

## 8. Hard Rules

```text
YES all backend infrastructure integration code is owned by backend/services/infrastructure
YES consumers use infrastructure/client contracts
YES client implementations delegate to infrastructure services
YES provider-specific behavior stays behind infrastructure services
YES services/helpers contains pure helper logic only
YES infrastructure owns durable technical state such as DKIM identity
YES DKIM startup reuses persisted identity and never silently rotates it

NO Haraka-specific client exposed to MailProvider
NO business service importing infrastructure impl packages
NO SMTP implementation inside MailProvider
NO database/network/filesystem work inside services/helpers
NO HTTP handlers directly inside infrastructure.go
NO silent DKIM key replacement during restart/recreate/rebuild
```
