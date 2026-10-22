# Git & Pull Request Guidelines

> **Status:** Mandatory Phase 0 engineering standard  
> **Location:** `docs/engineering/git/GIT_AND_PR_GUIDELINES.md`  
> **Applies to:** branches, commits, pull requests, CODEOWNERS, issue linking, reviews, generated code, merge policy, forks, dependency updates, release branches, and open-source contribution workflow  
> **Primary model:** protected `main` + short-lived branches + reviewed pull requests  
> **Core principles:** small reviewable changes, traceability, automated gates, ownership, clean history, no hidden generated diffs, and contributor-friendly open-source workflow

---

## 1. Objective

Git history should explain:

```text
what changed
why it changed
who reviewed it
what checks passed
what compatibility risk existed
```

Repository workflow must stay simple for external contributors.

---

## 2. Branch Model

Use:

```text
main
  ↑
short-lived feature/fix/docs/chore branches
```

Avoid long-lived environment branches.

Release branches are created only when maintaining an older supported release line requires them.

---

## 3. Mandatory Rules

1. `main` is protected.
2. Normal changes merge through PR.
3. Direct push to protected branches is restricted.
4. Branches are short-lived.
5. PRs have one coherent purpose.
6. Generated changes are separated/understandable.
7. Required CI checks pass before merge.
8. CODEOWNERS applies to high-risk areas.
9. Public contract/schema/security changes receive appropriate owner review.
10. Secrets never enter Git history.
11. Binary artifacts are not committed unless explicitly repository-owned.
12. Lockfiles are committed.
13. Formatting is automated.
14. Merge does not bypass failing tests without documented emergency process.
15. Force-push is acceptable on personal feature branches but not protected/shared branches.
16. Rewrite of published release tags is forbidden.
17. Tags are immutable.
18. Commit/PR messages are understandable.
19. Dependency update PRs still pass full relevant checks.
20. External contributors can follow documented setup/contribution steps.

---

## 4. Branch Naming

Recommended:

```text
feat/<short-name>
fix/<short-name>
refactor/<short-name>
docs/<short-name>
test/<short-name>
chore/<short-name>
security/<short-name>
```

Include issue ID where useful, not mandatory for every community contribution.

---

## 5. Branch Scope

A branch should represent one reviewable outcome.

Avoid combining:

```text
frontend redesign
database migration
unrelated cleanup
provider refactor
```

unless they are inseparable parts of one feature.

---

## 6. Commit Style

Use Conventional Commit style for clarity/automation.

Examples:

```text
feat(catalog): add application lifecycle status
fix(execution): prevent duplicate retry scheduling
docs(frontend): define shell composition boundary
test(provider): add cache conformance cases
chore(deps): update playwright
```

---

## 7. Commit Types

Recommended:

```text
feat
fix
refactor
test
docs
build
ci
chore
perf
security
revert
```

Do not invent many near-duplicate types.

---

## 8. Commit Scope

Scope maps to meaningful ownership:

```text
catalog
execution
frontend
provider
api
database
scripts
release
```

Avoid usernames or ticket-only scopes.

---

## 9. Commit Message Body

Use body when behavior/risk is non-obvious.

Explain:

```text
why
compatibility
migration
tradeoff
```

Do not narrate every line changed.

---

## 10. Atomic Commits

Commits should be logically coherent.

During PR review, fixup commits are okay; final merge strategy may squash.

---

## 11. Merge Strategy

Default recommendation:

```text
Squash merge PRs into main
```

This keeps a clean high-level history while allowing iterative review commits.

Critical/release branches may use another strategy if governance requires.

---

## 12. PR Title

PR title should follow the same semantic style because it may become the squash commit.

Example:

```text
feat(provider): add scoped capability binding resolver
```

---

## 13. PR Description

Include:

```text
Summary
Why
What changed
Compatibility/migration
Testing actually run
Screenshots for visual changes
Risk/rollback where relevant
Linked issue/design
```

---

## 14. PR Size

Prefer reviewable PRs.

When a feature is large, split by safe boundaries:

```text
contract
backend implementation
frontend
regression
```

without merging broken intermediate states into main.

---

## 15. Draft PRs

Use draft PR early for:

```text
architecture visibility
cross-team feedback
CI proof
```

Draft does not excuse broken secrets/security policy.

---

## 16. Required Reviews

At least one qualified review for ordinary code.

Additional CODEOWNER/security/database/provider/API review for sensitive paths.

---

## 17. CODEOWNERS

Recommended ownership:

```text
/backend/services/<domain>/
frontend/packages/domain/<domain>/
docs/engineering/security/
docs/engineering/providers/
migrations/
deployment/
```

Ownership can evolve with teams.

---

## 18. Self-Approval

Authors should not be the only required reviewer for protected high-risk changes.

Emergency rules must be explicit/audited.

---

## 19. Review Focus

Review:

```text
correctness
architecture
security
compatibility
tests
operability
readability
scope
```

Formatting should be automated.

---

## 20. Architecture Review

Reviewers check applicable engineering guideline.

Do not approve a shortcut that creates forbidden dependency direction merely because it works.

---

## 21. API Review

Public contract changes must show:

```text
compatibility result
client impact
version decision
contract tests
```

---

## 22. Database Review

Schema PRs show:

```text
migration
lock/backfill risk
rollback/forward recovery
index impact
real PostgreSQL tests
```

---

## 23. Frontend Review

Check:

```text
base/business/components/domain/shell placement
public index.ts
no deep imports
story/test coverage
service/reader rules
SRP/OCP/DRY
```

---

## 24. Provider Review

Check:

```text
capability contract
binding/resolver usage
conformance
secret handling
telemetry
no vendor leak
```

---

## 25. Security Review

Use security checklist for triggered areas.

Security reviewer approval may be required for auth/secrets/remote execution/tenant boundaries.

---

## 26. Generated Code

PR must clearly show source contract changes that produced generated code.

Never approve hand-edited generated output.

---

## 27. Lockfiles

Dependency changes include lockfile updates generated by pinned package manager.

Unexpected broad lockfile churn should be investigated.

---

## 28. Dependency Updates

Automated dependency PRs should be small and grouped carefully.

Major upgrades require explicit compatibility/test review.

---

## 29. Secrets

Secret scan runs before merge.

If a secret is committed:

```text
revoke/rotate first
then remove from history where required
```

Deleting the line alone does not make a leaked secret safe.

---

## 30. Large Files

Use appropriate artifact/package/object storage.

Do not commit build outputs, test videos, node_modules, binaries, DB dumps.

---

## 31. `.gitignore`

Must cover:

```text
.tools
.runtime
.state
.env.local
auth states/tokens
build output
coverage
Playwright artifacts
IDE/OS noise where appropriate
```

Do not ignore source/config that must be reproducible.

---

## 32. Git LFS

Use only if the project genuinely requires large versioned assets.

Do not add LFS as default complexity.

---

## 33. Issues

Bug reports should include:

```text
version/commit
environment
reproduction
expected/actual
logs/artifacts with secrets removed
```

Feature requests describe user outcome, not only implementation preference.

---

## 34. PR Linking

Link issues/design/ADR where available.

Do not require an issue for trivial docs/typo/community improvements.

---

## 35. Changelog Labels

PR labels may support release notes:

```text
breaking
feature
bug
security
documentation
provider
deployment
```

Keep label taxonomy small.

---

## 36. CI Checks

Required baseline may include:

```text
format
lint
typecheck
unit
integration
contract
generated freshness
migration
provider conformance
security scans
smoke regression
```

Only run relevant expensive suites where selection is reliable.

---

## 37. Merge Queue

A merge queue may be introduced when concurrency causes main-branch race failures.

Not required in Phase 0.

---

## 38. Broken Main

Fix/revert quickly.

Do not stack new unrelated work on known broken main.

---

## 39. Reverts

Use explicit revert commits/PRs for published main changes when rollback is needed.

Do not rewrite shared history.

---

## 40. Release Tags

Tags follow release guideline.

Example:

```text
v1.2.3
```

Never move/delete a published release tag casually.

---

## 41. Release Branches

Only for supported maintenance lines.

Example:

```text
release/1.x
```

Do not create release branches for every normal release if tags from main are sufficient.

---

## 42. Hotfix

Hotfix still receives review/tests proportional to risk.

Emergency bypass requires after-action follow-up.

---

## 43. Open-Source Fork Workflow

External contributor:

```text
fork
branch
make bootstrap / targeted setup
make relevant tests
PR
```

Do not require access to private infrastructure for ordinary contributions.

---

## 44. Contributor License / DCO

Project governance should choose one documented contribution attestation mechanism.

A lightweight DCO-style process is recommended if legal/governance approves it.

Do not silently mix incompatible CLA/DCO expectations.

---

## 45. Commit Signing

Signed commits/tags may be required for maintainers/releases.

Do not block casual external contributors unless governance explicitly requires it.

Release tags/artifacts should become cryptographically verifiable.

---

## 46. Review Conversations

Resolve substantive comments explicitly.

Do not mark unresolved security/compatibility concerns resolved without addressing them.

---

## 47. AI-Generated Changes

AI-assisted code follows the same standards.

Author remains responsible for:

```text
understanding diff
tests
license/provenance concerns
security
correctness
```

Do not merge generated code that nobody reviewed.

---

## 48. PR Template

Recommended:

```text
## Summary

## Why

## Changes

## Compatibility / Migration

## Validation
- [ ] ...

## Screenshots / Evidence

## Risk / Rollback
```

---

## 49. Reviewer Checklist

```text
[ ] Scope coherent
[ ] Applicable guideline followed
[ ] No forbidden dependency
[ ] Public contract compatible
[ ] Security reviewed
[ ] Data migration safe
[ ] Tests meaningful
[ ] Generated code legitimate
[ ] Observability/operability adequate
[ ] Docs/runbooks updated
[ ] Actual validation results provided
```

---

## 50. Definition of Done

```text
[ ] PR explains outcome
[ ] Required reviewers approved
[ ] Required checks pass
[ ] No secrets
[ ] Public compatibility reviewed
[ ] Migration/release impact documented
[ ] Tests/docs updated
[ ] Merge strategy preserves clean history
```

---

## 51. Phase 0 Deliverables

Phase 0 repository governance includes:

```text
protected main
CODEOWNERS
PR template
bug/feature issue templates
branch naming guidance
conventional commit guidance
required CI checks
secret scan
dependency update automation
CONTRIBUTING.md
CODE_OF_CONDUCT.md
SECURITY.md
```

---

## 52. Agent Rules

Before changing repo governance/workflows inspect:

```text
docs/AGENT.md
this file
.github/
CONTRIBUTING.md
CI configuration
release guideline if tags/releases affected
```

Do not create a second competing workflow.

---

## 53. Final Mandatory Rules

1. Protected main.
2. Short-lived branches.
3. Reviewed PRs.
4. Small coherent changes.
5. Conventional semantic titles/commits.
6. Required checks before merge.
7. CODEOWNERS for sensitive areas.
8. No secrets/binaries/build outputs.
9. Generated code traces to source.
10. Release tags immutable.
11. External contribution path remains simple.
12. AI code receives normal human review.

---

## 54. Guiding Principle

```text
Small Change
    +
Clear Intent
    +
Automated Evidence
    +
Qualified Review
    =
Trustworthy Main Branch
```

> **Git history should be a reliable engineering record, not merely a storage mechanism for diffs.**
