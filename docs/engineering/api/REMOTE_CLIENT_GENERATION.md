# Remote client generation

Service-owned controller annotations and DTOs are the source of truth. Generate versioned OpenAPI specifications from each owning project and generate consumer-local typed clients from pinned specifications. No shared Contracts module or dependency on Platform implementation is permitted. Core owns generic generation and authenticated HTTP transport. Existing in-process clients remain service-owned. Platform business operations must be restricted; application-facing wrappers own browser authentication and authorization. CI must check deterministic generation, schema compatibility and consumer compilation. Do not discover schemas at request time.

Implementation of the OpenAPI consumer generator and its integration is pending.
