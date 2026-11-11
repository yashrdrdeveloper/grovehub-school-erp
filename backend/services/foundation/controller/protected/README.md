# Protected controllers

`PersonController` is exposed at `/foundation/protected/persons` through
the `@Routes /foundation` service prefix plus the package's `protected`
category and the interface's `@Routes /persons` resource prefix.

Currently the `protected` category does **not** enforce authentication.
Identity/gateway integration must implement and test the corresponding
policy before exposing these endpoints to untrusted clients.
