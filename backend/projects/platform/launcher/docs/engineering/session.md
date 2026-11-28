# Platform gateway session boundary

Import `grovehub/libs/platform/session` from the standalone module at `backend/libs/platform`. Launcher and Platform services share one security boundary for principals, cookies, and tenant generations. The library has no dependency on Launcher or a business service.

Only the capability facade is public; contracts, constants and operations are private to `session/internal/src`. The Platform project requires the library with a checkout-relative module replacement.

Never derive principals or trusted tenant checks from HTTP input. Authentication and tenant-generation semantics are unchanged. The shared library has an independent strict coverage gate, while Launcher retains annotated consumer-contract tests.
