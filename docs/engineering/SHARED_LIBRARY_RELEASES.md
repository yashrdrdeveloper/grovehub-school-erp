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
