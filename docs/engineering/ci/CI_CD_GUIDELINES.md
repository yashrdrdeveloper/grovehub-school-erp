# CI/CD Engineering Guidelines

> **Status:** Mandatory engineering standard  
> **Location:** `docs/engineering/ci/CI_CD_GUIDELINES.md`  
> **Applies to:** Woodpecker CI, pull-request pipelines, main/nightly/release pipelines, build caching, test orchestration, generated-code checks, security scanning, artifact publishing, deployment promotion, environments, approvals, secrets, and CI observability  
> **Reference CI:** Woodpecker CI  
> **Core principles:** same commands locally and in CI, deterministic pipelines, least-privilege secrets, fast feedback, parallel safe jobs, immutable promotion, explicit evidence, and no hidden CI-only architecture

---

## 1. Objective

CI/CD should automate the engineering contract already defined in repository scripts/guidelines.

It must not become a second undocumented build system.

---

## 2. Primary Flow

```text
PR
 ↓
format/lint/static
 ↓
unit/component
 ↓
integration/contract/conformance
 ↓
smoke/critical regression
 ↓
merge

main/nightly
 ↓
broader regression/security/performance

release
 ↓
build once
 ↓
verify
 ↓
publish immutable artifacts
 ↓
install/upgrade validation
 ↓
promote
```

---

## 3. Mandatory Rules

1. Woodpecker is the reference CI, but repository commands remain CI-provider-neutral.
2. CI calls Make/scripts where practical.
3. No duplicate hidden CI-only install logic.
4. Toolchain versions are pinned.
5. Dependencies use lockfiles.
6. Jobs are deterministic.
7. CI secrets are least privilege.
8. Untrusted fork PRs never receive release/production secrets.
9. Artifacts are immutable after publishing.
10. Release builds originate from protected source.
11. Caches improve speed but never correctness.
12. Cache miss must still produce valid build.
13. Generated-code freshness is checked.
14. Architecture/import rules are checked.
15. Security scans run automatically.
16. Test reports/artifacts are retained appropriately.
17. CI failures are actionable.
18. Required checks map to repository quality gates.
19. Flaky tests are not normalized through retries.
20. Heavy suites can move to scheduled/release lanes, but critical gates remain.
21. Deployment promotion reuses the same built artifact.
22. CI does not directly mutate production outside governed deployment workflow.
23. Manual approvals are used for high-risk environments when policy requires.
24. CI logs redact secrets.
25. Build metadata includes commit/run/version.
26. Pipeline definitions are code-reviewed.
27. `main` should remain releasable.
28. Pipeline concurrency/cancellation avoids wasting obsolete PR runs.
29. CI resource usage is bounded.
30. Release evidence is preserved.

---

## 4. Pipeline Classes

```text
PR
main
nightly/scheduled
release
security
provider integration
```

Keep responsibilities explicit.

---

## 5. PR Pipeline

Recommended stages:

```text
bootstrap validation
format/lint
type/static
unit
frontend component
backend integration
API contract
provider conformance
migration validation
smoke/critical regression
security scans
```

Parallelize independent stages.

---

## 6. Fast Feedback

Put cheapest/high-signal checks first.

Example:

```text
format/type/lint before long E2E
```

---

## 7. Main Pipeline

Adds:

```text
build distributable artifacts
broader integration
full regression where affordable
snapshot/dev image for internal use if policy allows
```

Do not publish GA version from every main commit.

---

## 8. Nightly

Good for:

```text
full browser matrix
full regression
provider matrix
race detector
performance baseline
dependency scan refresh
flakiness repetition
```

---

## 9. Release Pipeline

Uses release guideline.

Required:

```text
clean protected commit/tag
all release gates
immutable build
checksums
SBOM
signing/provenance
install/upgrade validation
publication
```

---

## 10. Commands

Prefer:

```bash
make test-frontend
make test-backend
make test-regression-smoke
```

over copied npm/go commands.

Internal jobs may call scripts directly when necessary.

---

## 11. Bootstrap

CI uses non-interactive bootstrap.

No prompts.

No system mutation requiring sudo.

---

## 12. Caching

May cache:

```text
Go modules/build
pnpm store
Playwright browsers
toolchain archives
container layers
```

Cache keys include relevant lock/tool versions.

---

## 13. Cache Poisoning

Untrusted PR cache writes must not poison trusted release cache.

Use scope/separation as CI supports.

---

## 14. Artifacts

Retain:

```text
test reports
coverage
Playwright traces/screenshots/videos on failure
built candidate images/binaries
SBOM
scan results
release manifest
```

---

## 15. Test Reports

Use machine-readable formats where supported:

```text
JUnit
coverage formats
Playwright HTML
```

---

## 16. Logs

Pipeline output is concise.

Detailed logs/artifacts available for failed job.

Never echo secret environment.

---

## 17. Secrets

Use CI secret store.

Scope by:

```text
repository
event type
branch/tag
environment
job
```

---

## 18. Fork PRs

Run safe tests without privileged secrets.

Provider live integration requiring credentials is skipped or uses safe public/local emulation.

---

## 19. Generated Code Check

Workflow:

```text
run generator
git diff --exit-code
```

or equivalent.

Never merge stale generated output.

---

## 20. Architecture Checks

Enforce:

```text
frontend deep-import rules
domain/shell boundaries
backend forbidden imports
service boundaries
provider capability rules
circular dependency checks
```

---

## 21. Go Checks

As applicable:

```text
gofmt
go test
go vet
golangci-lint
go test -race selected
```

---

## 22. Frontend Checks

As applicable:

```text
format
ESLint
TypeScript
Jest/RTL
Storybook build/test
dependency boundary checks
```

---

## 23. API Checks

```text
OpenAPI validation
Buf lint
Buf breaking
client generation freshness
```

---

## 24. Database Checks

```text
migration ordering
clean migrate
upgrade migrate
real PostgreSQL integration
```

---

## 25. Provider Checks

```text
provider SDK compile
conformance suite
reference provider tests
protocol compatibility
```

---

## 26. Security Checks

```text
secret scan
dependency vulnerabilities
container scan
SAST where useful
license policy
SBOM on release
```

---

## 27. Regression

PR:

```text
Chromium smoke/critical
```

Nightly/release:

```text
broader regression/browser matrix
```

according to regression guideline.

---

## 28. Performance

Do not block every PR on heavy load.

Use baseline smoke for sensitive paths and scheduled/release deeper runs.

---

## 29. Matrix

Use matrix builds when dimensions genuinely matter:

```text
OS/arch
browser
provider
PostgreSQL supported version
Kubernetes supported version
```

Avoid combinatorial explosion.

---

## 30. Conditional Jobs

Condition based on reliable changed-path/dependency rules.

Critical global changes run broad gates.

---

## 31. Parallelism

Tests must own independent state.

CI parallelism is not allowed to create flaky shared-resource races.

---

## 32. Cancellation

Cancel superseded PR runs when safe.

Release runs are not auto-cancelled after publication begins.

---

## 33. Timeouts

Every job has bounded timeout.

A hung test does not consume runner indefinitely.

---

## 34. Retries

Infrastructure fetch/retry may be bounded.

Test job retries do not hide deterministic failures/flakiness.

---

## 35. Runners

Use isolated runners/containers.

Do not assume mutable state from previous job.

---

## 36. Privileged Runners

Avoid privileged Docker-in-Docker unless required.

Prefer safer build mechanisms/host runners with scoped permissions.

---

## 37. Container Build

BuildKit/buildx or approved builder.

Do not expose registry credentials to untrusted steps.

---

## 38. Build Once Promote

Release artifact:

```text
build
 ↓
scan/test
 ↓
publish immutable digest
 ↓
promote/deploy same digest
```

Do not rebuild separately for staging and production.

---

## 39. Environments

Model environments explicitly:

```text
dev
test
staging
production
```

CI credentials/approvals are scoped.

---

## 40. Deployment Approval

Production may require manual approval/change governance.

Approval records become evidence.

---

## 41. GitOps

When GitOps is used, CI updates desired state or produces artifacts; Flux reconciles.

Avoid direct kubectl mutation outside controlled workflows.

---

## 42. Rollback

CD workflow supports redeploying prior compatible artifact.

DB compatibility determines whether rollback is safe.

---

## 43. Release Failure

Publishing is staged.

If one artifact fails verification, do not declare release complete.

---

## 44. CI Observability

Track:

```text
pipeline duration
queue time
failure category
flaky test rate
cache hit
runner saturation
release lead time
```

Do not expose secret metadata.

---

## 45. Required Status Checks

Branch protection references stable job/check names.

Renaming required checks is coordinated.

---

## 46. CI Configuration Ownership

Cross-cutting pipeline config has explicit owners.

Domain jobs can be modularized but must not create incompatible conventions.

---

## 47. Pipeline DRY

Reuse templates/functions carefully.

Avoid one giant generic pipeline framework that hides behavior.

---

## 48. Pipeline Security

Treat PR code as untrusted.

Do not execute untrusted code in a context holding production/release credentials.

---

## 49. Dependency Bots

Automated update PRs use normal CI.

Group low-risk updates thoughtfully.

---

## 50. CI PR Checklist

```text
[ ] Uses stable Make/script entry point where possible
[ ] No duplicate bootstrap
[ ] Tool versions pinned
[ ] Secrets scoped
[ ] Fork behavior safe
[ ] Cache correctness independent
[ ] Job timeout defined
[ ] Artifacts retained appropriately
[ ] Required checks stable
[ ] Security scans present
[ ] Release artifact immutable
[ ] Same artifact promoted
```

---

## 51. Definition of Done

```text
[ ] Local equivalent command exists
[ ] Pipeline passes from clean runner
[ ] Secrets not exposed
[ ] Caches can be deleted without changing correctness
[ ] Failures are diagnosable
[ ] Required evidence retained
[ ] Release path does not rebuild after verification
```

---

## 52. Phase 0 Deliverables

```text
Woodpecker PR pipeline
frontend/backend/test jobs
real PostgreSQL service
provider conformance job
OpenAPI/Buf checks
generated freshness
secret/dependency scan
container build/scan
Playwright smoke
artifact/report retention
main/release skeleton
```

---

## 53. Agent Rules

Before CI change inspect:

```text
docs/AGENT.md
this file
Makefile/scripts
testing guideline
Git/PR guideline
security guideline
release guideline if publishing/deployment changes
```

Do not fix CI by bypassing repository standards.

---

## 54. Final Mandatory Rules

1. CI reuses repository commands.
2. Deterministic clean runners.
3. Least-privilege secrets.
4. Safe fork PRs.
5. Cache is optimization only.
6. Architecture/contracts/security are automated gates.
7. Artifacts are immutable.
8. Build once, promote same artifact.
9. Release evidence retained.
10. Failures actionable.
11. Main stays releasable.
12. CI is not a second architecture.

---

## 55. Guiding Principle

```text
Same Local Commands
      +
Clean Deterministic Runner
      +
Automated Architecture Gates
      +
Least-Privilege Secrets
      +
Immutable Promotion
      =
Trustworthy Delivery Pipeline
```

> **CI should prove the repository contract, not create a different one.**
