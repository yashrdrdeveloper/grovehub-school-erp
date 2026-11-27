# Shared library release candidates

The shared libraries remain monorepo-owned during development. No registry publication or external repository creation is implied by this change.

| Library | Source | Candidate |
|---|---|---|
| Go Core | backend/libs/core | v0.1.0-ci.<commit prefix> |
| Base | frontend/libs/base | 0.1.0-ci.<commit prefix> |
| Components | frontend/libs/components | 0.1.0-ci.<commit prefix> |
| Styles | frontend/libs/styles | 0.1.0-ci.<commit prefix> |

The shared-library-candidates workflow creates immutable commit-specific candidates, then independently resolves Core from a standard file-based Go proxy and builds a copied Internal Tool frontend with the three packed source dependencies. Frontend candidates exclude tests and package-local tooling and convert workspace dependencies to exact candidate versions. React remains a peer of Components. Consumers need TypeScript/Vite, CSS support, shared alias mappings to installed package roots and the shell-installed global translation function. This is a source distribution contract, not a compiled JS/declaration distribution. Third-party resolution in the standalone smoke test is not locked to the workspace lockfile.

Core's existing logical module path is different from its physical repository subdirectory. Do not create a misleading Go Git tag: a real release needs either a canonical path migration or an authenticated module proxy that serves the existing path. The file proxy verifies the candidate format and use without a replace directive; it does not make the module publicly reachable. Frontend candidates intentionally retain private:true. Registry ownership, visibility, credentials, stable versions, full independent School consumer validation and consumer lockfiles require a separate explicit rollout. Never publish automatically from a push. Current production consumers retain local development replacements/workspace links until that rollout.

## Independent application consumer gate

Candidate CI now exports provider-owned API contracts from the same commit into a verified archive (not a second editable contracts library). School and Internal Tool SDK generation accept explicit GROOVEHUB_CONTRACT_ROOT and GROOVEHUB_CLIENT_GENERATOR paths; monorepo defaults remain available for development. The isolated Go consumer gate copies each actual project into a temporary directory, removes the local Core replacement, pins the candidate version, regenerates clients from the provider artifact and captures go.mod/go.sum before running tests and builds with -mod=readonly. PostgreSQL-dependent tests are still exercised by the separate project database CI.

Frontend consumer verification copies both applications outside the workspace, references the candidate tarballs through relative paths, generates a candidate-specific lockfile, then builds after a frozen install. Internal Tool unit tests run in this gate; School source and shared-library test/architecture/browser suites remain in the existing monorepo frontend pipeline. The release artifacts retain lock snapshots, not yet committed application lockfiles. Initial resolution of third-party semver ranges can differ from the monorepo lockfile.

The existing Core logical module path is preserved using the private-proxy distribution approach. A public Git-backed module path rename is deferred rather than breaking all importers or publishing an inconsistent module tag. The file-based proxy is CI-only; production registry infrastructure and access configuration are not provisioned.
