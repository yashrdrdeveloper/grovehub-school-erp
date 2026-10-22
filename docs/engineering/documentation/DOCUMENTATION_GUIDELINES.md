# Documentation Guidelines

> **Status:** Mandatory engineering standard  
> **Location:** `docs/engineering/documentation/DOCUMENTATION_GUIDELINES.md`  
> **Applies to:** README files, engineering standards, runbooks, architecture docs, ADRs, API docs, provider docs, troubleshooting, examples, diagrams, release/upgrade docs, contributor docs, code comments, and documentation ownership  
> **Core principles:** documentation close to ownership, one source of truth, explicit audience, executable examples, version awareness, searchable structure, link integrity, and no duplicate contradictory specifications

---

## 1. Objective

Documentation should let:

```text
a new user run Perihellion
a contributor find the correct architecture boundary
an operator recover from failure
a provider author implement a capability
an AI coding agent gather correct context
```

without tribal knowledge.

---

## 2. Documentation Types

```text
Product docs
Engineering standards
Architecture/design
Runbooks
API/reference
Provider/SDK docs
Contributor docs
Operational troubleshooting
Release/upgrade
ADRs
Examples/tutorials
```

Each has a different audience/purpose.

---

## 3. Mandatory Rules

1. Every document states purpose/audience where non-obvious.
2. One concept has one authoritative source.
3. Other docs link to it instead of copying large blocks.
4. Docs live near the code/domain they govern when useful.
5. Repository-wide standards live under `docs/engineering`.
6. Package/service runbooks live with the package/service.
7. Public API reference derives from machine-readable contracts.
8. Generated API docs are not manually edited.
9. Examples must remain valid.
10. Commands use stable Make/script interfaces when available.
11. Never place real secrets in examples.
12. Use realistic but obviously fake IDs/domains/credentials.
13. Architecture diagrams reflect current design, not aspirational features unless clearly labeled.
14. Future design is labeled future/proposed.
15. Breaking behavior updates docs in the same PR.
16. New public component/provider/service has corresponding documentation where required.
17. Runbooks include verification and failure recovery.
18. ADRs record significant irreversible/cross-cutting decisions.
19. Documentation ownership is explicit.
20. Links are validated in CI where practical.
21. Markdown style is consistent.
22. Headings are meaningful and searchable.
23. Avoid giant catch-all documents when a topic has clear subdomains.
24. Do not create `misc.md`.
25. Documentation should be concise where possible but complete where contract/operation requires it.
26. Screenshots are supplemental, not the only instruction.
27. Accessibility text/alt descriptions are provided for meaningful diagrams/images.
28. Version-specific docs identify version applicability.
29. Deprecated instructions are removed or clearly marked.
30. AI-agent docs must point to canonical standards, not duplicate them.

---

## 4. Repository Documentation Tree

Recommended:

```text
docs/
├── AGENT.md
├── architecture/
├── concepts/
├── getting-started/
├── installation/
├── operations/
├── providers/
├── api/
├── contributing/
├── decisions/
└── engineering/
```

Engineering standards are separate from user/operator guides.

---

## 5. Engineering Standards

Canonical structure:

```text
docs/engineering/
├── README.md
├── backend/
├── frontend/
├── testing/
├── scripts/
├── api/
├── database/
├── providers/
├── security/
├── observability/
├── deployment/
├── release/
├── git/
└── documentation/
```

---

## 6. `docs/engineering/README.md`

Acts as index.

It should state:

```text
which standard applies
when to read it
dependencies between standards
```

Do not repeat full rules from each file.

---

## 7. Root README

The root `README.md` is product-facing.

First screen answers:

```text
What is Perihellion?
Why would I use it?
How do I run it?
Where are docs?
How do I contribute?
```

Do not start with internal architecture details.

---

## 8. Quick Start

Use stable commands:

```bash
make dev-start
```

or released installer workflow.

Avoid a 20-step manual dependency setup when scripts own it.

---

## 9. CONTRIBUTING

Contains:

```text
development prerequisites
bootstrap
test commands
branch/PR process
engineering standards link
DCO/CLA policy
code of conduct
```

---

## 10. SECURITY.md

Contains public vulnerability reporting process.

Do not publish sensitive exploit workflow.

---

## 11. CODE_OF_CONDUCT

Use an established open-source code of conduct chosen by governance.

---

## 12. Runbooks

A runbook is operational/development guidance for a bounded component/package.

It includes:

```text
purpose
ownership
dependencies
normal operation
commands
health
common failures
diagnostics
recovery
escalation/links
```

---

## 13. Frontend Package Runbooks

Required:

```text
frontend/packages/base/RUNBOOK.md
frontend/packages/business/RUNBOOK.md
frontend/packages/components/RUNBOOK.md
frontend/packages/shell/RUNBOOK.md
frontend/packages/domain/<domain>/RUNBOOK.md
```

according to frontend guideline.

---

## 14. Backend Service Runbook

Each significant backend service should eventually document:

```text
responsibility
public contracts
owned data
dependencies
health
operations
common failures
extraction notes when relevant
```

---

## 15. Provider Documentation

Each provider documents:

```text
capabilities
configuration
secret references
permissions
health
limits
known semantic differences
installation
conformance status
```

---

## 16. ADRs

Architecture Decision Records are for significant decisions.

Location:

```text
docs/decisions/
```

Naming:

```text
0001-use-postgresql-as-primary-store.md
0002-provider-capability-model.md
```

---

## 17. ADR Template

```text
Title
Status
Date
Context
Decision
Consequences
Alternatives considered
References
```

Do not use ADRs for trivial code choices.

---

## 18. ADR Status

Use:

```text
Proposed
Accepted
Superseded
Deprecated
Rejected
```

Superseded ADR links replacement.

---

## 19. Architecture Docs

Describe current component relationships and boundaries.

Use diagrams plus prose.

Clearly mark:

```text
Current
Planned
Experimental
```

---

## 20. Mermaid

Mermaid is preferred for diagrams that benefit from diffable source.

Do not use unreadable 200-node diagrams.

---

## 21. Text Diagrams

ASCII trees/flows are encouraged for simple ownership/dependency rules.

They render everywhere and match existing engineering-standard style.

---

## 22. Screenshots

Use for UI guidance when useful.

Keep them updated.

Never make critical setup possible only by interpreting screenshot pixels.

---

## 23. API Reference

Generated from:

```text
OpenAPI
Protobuf
provider schema/contracts
```

Manual docs explain semantics/examples, not duplicate every field.

---

## 24. Examples

Examples should use current public interfaces.

When possible, test/compile examples in CI.

---

## 25. Shell Commands

Prefer:

```bash
make backend-start
make test-regression-smoke
```

over internal commands unless documenting the internal layer specifically.

---

## 26. Environment Variables

Document:

```text
name
purpose
default
required
secret?
example
scope
```

Never include real credentials.

---

## 27. Configuration Reference

Generate where possible from typed configuration/schema.

Manual text adds semantic explanation.

---

## 28. Troubleshooting

Troubleshooting flow:

```text
symptom
checks
expected result
diagnostic command
likely cause
safe remediation
when to escalate
```

Avoid "restart everything" as default answer.

---

## 29. Error Codes

Public stable error codes should have searchable reference.

Map:

```text
code
meaning
common causes
client/operator action
```

---

## 30. Runbook Safety

Destructive commands are visibly labeled.

Require explicit placeholders/confirmation.

Do not provide copy-paste commands that can destroy all clusters/databases accidentally.

---

## 31. Versioned Docs

For GA releases, docs should be selectable/versioned when incompatible changes occur.

`main` docs describe development/current-next only if clearly labeled.

---

## 32. Deprecation Docs

Mark deprecated command/API/config:

```text
Deprecated since
Replacement
Removal target
Migration
```

---

## 33. Searchability

Use exact product concepts in titles/headings.

Good:

```text
Provider Binding Resolution
Mission Runtime Retry Semantics
Frontend Shell Registration
```

Bad:

```text
Other notes
Misc
Stuff
```

---

## 34. Terminology

Maintain a glossary for platform vocabulary.

Use one canonical name:

```text
Application
Environment
Mission
Provider
Capability
Binding
Evidence
Microcluster
Workbench
```

Do not rename concepts casually between docs/code/UI.

---

## 35. File Naming

Use descriptive `UPPER_SNAKE_CASE.md` for repository-wide standards where current convention uses it.

Use readable lowercase/slug names for user docs/ADRs as appropriate.

Consistency within a folder matters.

---

## 36. Heading Structure

Engineering standards use:

```text
# Title
status block
---
## numbered sections
```

and finish with:

```text
PR/checklist
Definition of Done
Phase 0 deliverables when relevant
Agent Rules
Final Mandatory Rules
Guiding Principle
```

This matches Perihellion's existing standard style.

---

## 37. Good / Bad Examples

Use examples when a rule is easy to misunderstand.

Prefer concrete code/trees over abstract prose.

---

## 38. Normative Language

Use:

```text
must
must not
should
may
```

consistently.

"Mandatory" sections are enforceable standards.

---

## 39. Source of Truth

When two docs overlap:

```text
specialized guideline owns detail
higher-level document links to it
```

Example:

```text
TESTING_GUIDELINES -> general testing
REGRESSION_TESTING -> Playwright specifics
```

---

## 40. Duplicate Content

Small summaries are acceptable.

Do not duplicate hundreds of lines; they drift.

---

## 41. Code Comments

Comments explain:

```text
why
non-obvious invariant
workaround
security/compatibility rationale
```

Do not comment obvious syntax.

---

## 42. TODOs

TODO includes issue/reference when meaningful.

Do not use docs as an unowned indefinite TODO dump.

---

## 43. Generated Docs

Generated files clearly state:

```text
GENERATED — DO NOT EDIT
source:
generation command:
```

---

## 44. Documentation Build

If a docs site is introduced, source Markdown remains accessible in repository.

Docs build must fail on broken references/schema where reasonable.

---

## 45. Link Style

Prefer relative repo links for repository docs.

External links should point to authoritative sources.

---

## 46. Link Checking

CI can check:

```text
broken internal links
missing anchors
selected external links
```

Avoid flaky blocking checks on arbitrary internet availability.

---

## 47. Spell/Style Checking

May be automated, but product-specific vocabulary has an allowlist.

Do not make style tooling hostile to contributors.

---

## 48. Localization

Do not localize engineering source docs initially unless product documentation strategy requires it.

Keep canonical technical docs in English.

---

## 49. Accessibility

Images/diagrams have meaningful alt/description.

Avoid color-only meaning in diagrams.

---

## 50. Security

Docs/examples never contain:

```text
real tokens
production hostnames with sensitive topology
customer data
private keys
internal exploit details
```

---

## 51. Documentation Review

Docs change is reviewed by code/domain owner when it changes behavioral contract.

Documentation-only typo fixes can be lightweight.

---

## 52. Documentation Ownership

CODEOWNERS can cover:

```text
docs/engineering/security/
docs/engineering/providers/
docs/engineering/api/
```

---

## 53. Documentation PR Checklist

```text
[ ] Audience/purpose clear
[ ] Canonical terminology used
[ ] No duplicate source of truth
[ ] Commands are current
[ ] Examples safe/current
[ ] Future vs current clearly labeled
[ ] Links valid
[ ] Destructive commands clearly marked
[ ] Version applicability clear
[ ] Related code/contract docs updated
```

---

## 54. Definition of Done

```text
[ ] New behavior documented at correct level
[ ] Public contract docs generated/current
[ ] Runbook updated when operation changed
[ ] Upgrade/deprecation docs updated when needed
[ ] Links/examples validated
[ ] No secrets/customer data
```

---

## 55. Phase 0 Deliverables

```text
root README
CONTRIBUTING
SECURITY
CODE_OF_CONDUCT
docs/AGENT.md
docs/engineering/README.md
engineering guideline tree
architecture overview
glossary
ADR template + first foundational ADRs
package/service runbook templates
link-check baseline
```

---

## 56. Agent Rules

Before documentation changes inspect:

```text
docs/AGENT.md
this file
canonical source document
actual code/contract being documented
```

Do not "fix" a documented contract from general knowledge when source code/contracts say otherwise; flag the mismatch.

---

## 57. Final Mandatory Rules

1. One source of truth.
2. Docs close to ownership.
3. Contracts generated from machine-readable source.
4. Runbooks describe operation/recovery.
5. ADRs capture significant decisions.
6. Examples remain executable/current.
7. Current vs future is explicit.
8. Terminology is consistent.
9. No secrets/customer data.
10. Engineering standards use consistent structure.
11. Docs change with behavior.
12. AI agents are routed to canonical docs.

---

## 58. Guiding Principle

```text
Correct Source
     +
Clear Audience
     +
Searchable Structure
     +
Executable Examples
     +
Owned Runbooks
     =
Useful Documentation
```

> **A Perihellion document should reduce dependency on tribal knowledge, not create another place for it to hide.**
