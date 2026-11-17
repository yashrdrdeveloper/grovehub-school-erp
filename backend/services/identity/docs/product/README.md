# Identity — Product Documentation

**Scope:** verified or source-supported Identity use cases on branch `phase-2.1-identity-account`. General engineering rules are in `../engineering/README.md`.

## Available documented workflow

- [Registration API workflow](registration.md) — registration request structure, account/person interaction, expected result, retries, conflicts and reproduction template.

## Documentation policy

Each public API gets: actor, permissions, prerequisites, method and route, authentication, headers, exact request and response schema, idempotency, validation, success and error cases, side effects and executable verification steps. Only use paths and responses verified from route wiring/OpenAPI and tests; label uncertainty clearly.

## Planned capability areas, not verified in this document

Login, access/refresh-token lifecycle, logout, token revocation, account recovery, password reset, MFA, OIDC, OAuth2, SAML, session management, RBAC and policy evaluation require independent implementation verification. Their presence on a roadmap is not evidence they are exposed by this branch.

## Engineering separation

Business behavior and API reproduction live here. Class diagrams, GORM internals, transaction propagation, data modeling, annotations, design patterns and ownership decisions live in `../engineering/`.
