# Shell
The shell is a composition-only package. It constructs React Router, TanStack Query, auth provider and the workspace chrome. Feature screens are imported from `@Domain/<name>` public APIs. No shell file imports `@Business`, calls `fetch`, or performs a use case.
