# Integration — Product and API Workflows

This directory explains **business use cases and API execution**, not internal algorithms.

- [api_workflows.md](api_workflows.md): provider catalog, connection create/list/get/disconnect, OAuth start/callback, tests, watch registration, webhook, and actions.
- [mail_actions.md](mail_actions.md): provider-specific mail read/send workflows, inputs and expected outcomes.
- OpenAPI: `internal/resources/api/openapi/v1/openapi.yaml`.
- Engineering: [../engineering/README.md](../engineering/README.md).

**Security:** Connection operations require a **verified owner/principal context**. Do not make production calls by forging identity headers. Local development may have an explicitly loopback-restricted mode; it is not an authorization strategy for deployment. Use provider sandbox/testing mailboxes and never paste real OAuth tokens in issue comments or docs.
