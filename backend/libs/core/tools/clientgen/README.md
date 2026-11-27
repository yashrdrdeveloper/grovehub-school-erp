# Consumer-local API clients

Use the pinned OpenAPI Generator CLI adapter to generate Go clients from a service-owned versioned OpenAPI specification. Run `generate.sh SPEC OUTPUT_DIR PACKAGE_NAME` with Docker available and the image pre-pulled. Output must be a fresh directory within the consumer project. Review generated changes before replacing an existing client. The consumer owns dependencies and injects its Core workload HTTP client into the generated configuration. Never use a browser token for Platform calls.

Existing Foundation, Identity and Integration specifications are migration inputs, not verified restricted contracts. Do not deploy generated clients until controller paths, security declarations and schemas agree. OpenAPI 3.1 nullable/union and PATCH omission semantics require explicit compatibility tests. The adapter does not yet automate annotation-to-spec generation or @Client integration. Audit and Infrastructure specifications still need inventory and completion.
