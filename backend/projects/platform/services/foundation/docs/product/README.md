# Foundation — Product Workflows

Foundation supports Person creation/lookup, search, partial update, contact management through the Person aggregate, and deactivation.

- [api_workflows.md](api_workflows.md) contains **all current public Person endpoints**, request examples, expected outcomes, negative cases and repeatable steps.
- Canonical contract: `internal/resources/api/openapi/v1/openapi.yaml`.
- Engineering implementation: [../engineering/README.md](../engineering/README.md).

**Security prerequisite:** API route names include `/protected/` but current routing alone does not enforce IAM. Run scenarios only in an isolated development environment behind appropriate authenticated access, never on an exposed production listener.

Use synthetic names/contact values; avoid real user data in test transcripts. Client examples below are illustrative and must be compared to the current OpenAPI when behavior changes.
