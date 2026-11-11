# Restricted controllers

Place internal-only HTTP controller interfaces here (package `restricted`).
With the service-level `@Routes /foundation`, generated endpoints start
with `/foundation/restricted/`.

Restricted is **not** a security boundary today. Before adding an actual
restricted endpoint, enforce internal-only access at the gateway or bind the
route to an internal-only listener. Do not put privileged maintenance or
data-correction handlers here while these URLs are externally reachable.
