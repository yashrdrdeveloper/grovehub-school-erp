# Public controllers

Place public HTTP controller interfaces in this folder (package `public`).
Use `// @Routes /<resource>` on the interface. With the service-level
`@Routes /foundation`, generated endpoints start with `/foundation/public/`.

Public endpoints are intended to be usable without authentication **once
gateway access policy is implemented**. Route generation itself does not
enforce authorization, authentication, or rate limiting.
