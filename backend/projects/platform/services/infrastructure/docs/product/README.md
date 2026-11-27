# Infrastructure — Messaging Use Cases

This module exposes a **typed Go client**, not HTTP routes. These workflows describe the observable behavior for module developers and downstream consumers.

- [messaging_workflows.md](messaging_workflows.md) — direct event, durable outbox event, consuming, replay/retry, failure reproduction and acceptance criteria.
- Engineering: [../engineering/README.md](../engineering/README.md).

Do not invent REST endpoints for Messaging. Application-facing business workflows belong to the services that publish/consume messages; these are platform delivery workflows.
