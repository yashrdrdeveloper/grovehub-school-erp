# Release Guidelines

> **Status:** Mandatory engineering standard  
> **Location:** `docs/engineering/release/RELEASE_GUIDELINES.md`  
> **Applies to:** versions, Git tags, changelog, release notes, release branches, OCI images, Helm charts, Compose bundles, CLI/Agent binaries, checksums, SBOM, provenance/signatures, provider compatibility, upgrade notes, release gates, rollback, and air-gap bundles  
> **Versioning:** Semantic Versioning for core product releases  
> **Core principles:** one coherent product release, immutable artifacts, reproducibility, tested upgrades, explicit compatibility, verifiable distribution, and honest release evidence

---

## 1. Objective

A release is not just a Git tag.

It is a tested, immutable, documented set of compatible artifacts.

---

## 2. Product Release Unit

Initially the monorepo ships one product version:

```text
Perihellion vX.Y.Z
```

Core components from the same release are tested together.

---

## 3. Release Artifacts

May include:

```text
server OCI image
web OCI image
CLI binaries
Local Development Agent binaries
Helm chart
Compose bundle
migration assets
checksums
SBOMs
release manifest
release notes
air-gap archive
provider packages
```

---

## 4. Mandatory Rules

1. Published versions are immutable.
2. Git release tags are immutable.
3. Release artifacts are immutable.
4. SemVer governs product version.
5. Public breaking changes require major-version consideration.
6. Release notes describe user-visible changes.
7. Upgrade notes describe migrations/config/provider compatibility.
8. Artifacts record source commit.
9. Checksums are published.
10. SBOM is generated.
11. Release signing/provenance is required as the release pipeline matures and before enterprise trust claims.
12. Release credentials are isolated from untrusted PRs.
13. Release comes from protected reviewed source.
14. Release gate includes build, test, security, install, and upgrade evidence.
15. Do not claim a platform/profile/provider is supported unless included in its required release validation.
16. Release automation uses pinned toolchains.
17. Changelog is generated/curated from reviewed changes, not raw commit dump only.
18. Database migration compatibility is explicit.
19. Rollback/forward recovery limitations are explicit.
20. Air-gap bundle contains every runtime artifact.
21. OCI references use immutable digests in release manifest.
22. Provider SDK compatibility window is documented.
23. Release failure does not partially overwrite an existing immutable version.
24. Security fixes identify affected versions without exposing exploit detail prematurely.
25. Patch releases do not silently introduce breaking contracts.
26. Release notes link migration/operational actions.
27. Release evidence is retained.
28. Reproducibility is tested where practical.
29. Deprecations include removal horizon.
30. Artifact naming is consistent.

---

## 5. SemVer

```text
MAJOR.MINOR.PATCH
```

Use:

```text
MAJOR  incompatible public/product contract change
MINOR  backward-compatible capability/feature
PATCH  backward-compatible fix/security/operational change
```

Internal refactor alone does not force major version.

---

## 6. Pre-1.0

Before `1.0.0`, compatibility is still treated seriously.

Do not use `0.x` as permission for arbitrary breakage.

Document pre-1.0 stability expectations.

---

## 7. Prerelease

Use:

```text
v1.2.0-alpha.1
v1.2.0-beta.1
v1.2.0-rc.1
```

Prerelease artifacts are clearly marked and not silently promoted to GA.

---

## 8. Build Metadata

Use SemVer build metadata only when needed.

Artifact manifest separately records commit/build ID.

---

## 9. Tag Format

Canonical:

```text
v1.2.3
```

Do not move a published tag.

---

## 10. Branches

Normal releases can come from protected `main`.

Use maintenance branches only for supported older release lines:

```text
release/1.x
```

---

## 11. Changelog

Maintain:

```text
CHANGELOG.md
```

with user/administrator/provider-author relevant changes.

Do not fill it with every formatting/refactor commit.

---

## 12. Release Notes

Sections may include:

```text
Highlights
Features
Fixes
Security
Deprecations
Breaking changes
Upgrade notes
Provider compatibility
Known issues
```

---

## 13. Breaking Changes

Require:

```text
explicit label/approval
migration path
version decision
consumer impact
documentation
test evidence
```

---

## 14. Deprecation

A deprecation has:

```text
introduced version
replacement
warning
planned removal release/window
```

---

## 15. Release Manifest

Publish one machine-readable manifest containing:

```text
product version
source commit
artifact names
digests/checksums
architectures
Helm chart version
provider SDK/protocol version
minimum supported upgrade version
```

---

## 16. OCI Images

Publish immutable images.

Record digest for:

```text
server
web
```

and any release component.

---

## 17. CLI/Agent Binaries

Publish per supported OS/architecture:

```text
darwin/arm64
darwin/amd64 where needed
linux/amd64
linux/arm64
windows/amd64 if native support is shipped
```

Do not publish an untested target.

---

## 18. Checksums

Publish SHA-256 or approved stronger checksum file.

Air-gap/import tooling verifies it.

---

## 19. SBOM

Generate SBOM for:

```text
OCI images
binaries where tooling supports
air-gap bundle
```

Use an industry-supported format such as SPDX/CycloneDX.

---

## 20. Provenance

Release pipeline should produce build provenance/attestation.

This strengthens supply-chain verification.

---

## 21. Signing

Sign release artifacts/images/tags according to chosen trust model.

Signing keys/identity are protected from PR workloads.

---

## 22. Helm Chart

Chart version follows Helm SemVer.

`appVersion` identifies Perihellion product version.

---

## 23. Compose Bundle

Release bundle contains pinned image versions/digests and sample environment/config documentation.

Do not ship `latest` in a release Compose file.

---

## 24. Database Migrations

Release identifies:

```text
schema changes
migration duration/risk
backup requirement
rollback limitations
```

---

## 25. Upgrade Path

Define supported direct upgrade path.

Example policy:

```text
latest patch of previous minor -> current
```

or a broader range after testing.

Do not assume arbitrary N-year leap upgrades work.

---

## 26. Upgrade Test

Release gate performs:

```text
install supported previous version
create representative data
upgrade
run migrations
verify critical workflows
```

---

## 27. Fresh Install Test

Release gate installs target version from empty environment and runs critical smoke.

---

## 28. Rollback Test

When rollback is documented as supported, test it.

If schema/data makes rollback unsafe, document forward recovery instead.

---

## 29. Configuration Migration

Deprecated/renamed configuration includes migration guidance.

Startup should provide actionable error/warning.

---

## 30. Provider Compatibility

Release notes state:

```text
Provider SDK/protocol version
minimum/maximum supported provider package compatibility if applicable
reference provider versions
```

---

## 31. Marketplace Compatibility

Marketplace packages declare compatible core/provider SDK ranges.

Core should reject incompatible package cleanly.

---

## 32. Frontend Compatibility

Server and web from one release are the tested pairing.

If independent serving/caching can create mixed versions, APIs remain compatible across supported rollout window.

---

## 33. API Compatibility

Run OpenAPI/protobuf breaking checks.

Breaking changes are never discovered only after publishing SDKs.

---

## 34. Migration Compatibility

Run DB migration tests from supported previous schema.

---

## 35. Regression Gate

Run:

```text
smoke
critical
full required regression
upgrade regression
install regression
```

according to release profile.

---

## 36. Cross-Browser Gate

Run required browser matrix from regression guideline before GA.

---

## 37. Security Gate

Run:

```text
secret scan
dependency scan
container scan
SBOM generation
license/policy checks where applicable
```

Block unresolved release-critical findings.

---

## 38. Vulnerability Exception

Requires:

```text
risk assessment
affected component
mitigation
owner
expiry
```

---

## 39. Performance Gate

Maintain baseline/performance smoke for critical APIs/workflows.

Large benchmark regressions require review.

---

## 40. Release Candidate

Use RC to validate:

```text
installation
upgrade
air-gap
docs
providers
operational runbooks
```

where release maturity warrants it.

---

## 41. Known Issues

Publish known important defects/limitations.

Do not hide them to make release notes look cleaner.

---

## 42. Support Matrix

Document:

```text
supported browsers
supported Kubernetes versions
supported architectures
supported PostgreSQL range
supported upgrade origin
```

Only include tested claims.

---

## 43. Release Promotion

Build once where possible, then promote immutable artifacts across environments.

Do not rebuild different bits for production after staging verification.

---

## 44. Reproducible Build

Same source/toolchain should yield functionally equivalent artifacts.

Record toolchain versions and source.

---

## 45. Release Failure

Publishing should be transactional where possible:

```text
stage artifacts
verify
publish manifest/release
```

Do not overwrite an existing version.

---

## 46. Patch Release

Patch includes focused fixes.

Avoid unrelated large refactors/dependency majors.

---

## 47. Security Release

Coordinate disclosure and release.

Rotate compromised credentials separately from code release.

---

## 48. Hotfix

Still runs proportional required tests/security checks.

Emergency bypass is documented and followed by postmortem/follow-up.

---

## 49. Air-Gap Bundle

Bundle contains:

```text
OCI images or image archive
Helm chart
Compose where supported
CLI/agent binaries
checksums
SBOM
provider packages
docs/install instructions
release manifest
```

---

## 50. Air-Gap Validation

Test install in network-isolated environment.

Searching logs for failed public downloads is not enough; block egress during validation.

---

## 51. Download Naming

Consistent names:

```text
perihellion_<version>_<os>_<arch>.tar.gz
perihellion-agent_<version>_<os>_<arch>.tar.gz
```

Exact final naming can be frozen before GA.

---

## 52. Release Documentation

Docs are versioned with release or published with version selector.

Upgrade docs remain available for supported versions.

---

## 53. Release Evidence

Retain:

```text
commit
CI run
test results
security results
artifact digests
SBOM/provenance
upgrade result
approvals
```

---

## 54. Release Checklist

```text
[ ] Version decided
[ ] Changelog/release notes
[ ] Breaking/deprecation reviewed
[ ] API compatibility checks
[ ] DB migration checks
[ ] Unit/integration/contract/conformance pass
[ ] Playwright release regression pass
[ ] Fresh install pass
[ ] Supported upgrade pass
[ ] Security scans pass/triaged
[ ] Images/binaries built
[ ] Checksums generated
[ ] SBOM generated
[ ] Signatures/provenance generated when enabled
[ ] Air-gap bundle verified when supported
[ ] Docs published
[ ] Tag/artifacts immutable
```

---

## 55. Definition of Done

```text
[ ] Release can be installed
[ ] Release can upgrade from supported prior version
[ ] Artifacts are identifiable/verifiable
[ ] Compatibility claims have evidence
[ ] Migrations are understood
[ ] Release notes are accurate
[ ] Recovery/known issues documented
[ ] Release evidence retained
```

---

## 56. Phase 0 Deliverables

```text
SemVer policy
release tag convention
CHANGELOG.md
release manifest format
server/web image publishing proof
CLI/agent artifact naming plan
checksums
SBOM proof
fresh-install smoke
upgrade-test harness skeleton
release CI skeleton
```

---

## 57. Agent Rules

Before release workflow changes inspect:

```text
docs/AGENT.md
this file
Git/PR guideline
deployment guideline
Docker/Kubernetes guideline
security guideline
testing/regression guideline
```

Never claim support based on unexecuted validation.

---

## 58. Final Mandatory Rules

1. One coherent product release.
2. SemVer.
3. Immutable tags/artifacts.
4. Source commit recorded.
5. Checksums and SBOM.
6. Compatibility checks.
7. Fresh install and upgrade tests.
8. Migration/recovery notes.
9. Security gate.
10. Air-gap bundle is complete/verifiable when supported.
11. Release evidence retained.
12. Support claims are tested.

---

## 59. Guiding Principle

```text
Reviewed Source
      +
Reproducible Build
      +
Immutable Artifacts
      +
Compatibility Evidence
      +
Upgrade Proof
      =
Trustworthy Release
```

> **A Perihellion release is something an operator can verify, install, upgrade, and recover from—not merely a version number.**
