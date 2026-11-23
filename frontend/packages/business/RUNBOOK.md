# Business — typed contracts, routes and transport

`packages/business` is reusable domain logic and API access. It must not import React presentation or Shell. Every backend call uses the shared credentialed `apiRequest` envelope; opaque access/refresh cookies are browser/API-owned, and transport endpoints are **never** browser navigation paths.

## Generic ownership standard

`<domain>/constants/` holds browser route IDs/config or genuinely domain-wide values; `<domain>/helpers/` holds domain-level route and pure business helpers. These folders are flat, with adjacent `types/` and `__tests__/` and **no index barrels**. `<domain>/types/` owns entity/contracts and an intentional type-only `index.ts` when required. **No `<domain>/types.ts`** at the domain root.

`<domain>/services/constants/` owns operation enums, transport endpoint metadata and proof schemes. `<domain>/services/helpers/` owns endpoint lookup, authorization header composition and typed encoded-resource helpers. Operations stay in named `services/<operation>/` units with `types/<operation>.request.ts`, response/callable contracts, unit tests and public `index.ts`. Reusable generic HTTP methods belong in `api/constants/`; runtime `ApiError` lives under `api/errors/`, not in the types directory.

All Identity public/protected endpoint URLs, proof-header schemes and operation IDs are centralized under `auth/services/`. Its domain root keeps **only** browser navigation metadata and the public type contract directory. Accounts, Foundation and Integrations similarly use `services/constants` and tested typed encoded-URL helpers.

Do not mock pure URL, route or mapping helpers in tests. Verify HTTP method, URL, body, CSRF/credential handling and authorization headers independently of presentation. Keep cross-domain imports on the owning package public entrypoint.

## Text and state ownership

Business entities, HTTP methods, backend operation/status codes, route metadata, error identifiers, schemas and validation data remain **raw technical values** or their own typed enums. Never run these through the global `__()` and never import the browser ShellContainer/i18n installer into Business. A domain presentation maps business error/status codes to stable human-readable copy in its own `constants/<owner>.copy.ts`, then calls `__()` while rendering. Test pure mapping and business utilities without mocking them.
