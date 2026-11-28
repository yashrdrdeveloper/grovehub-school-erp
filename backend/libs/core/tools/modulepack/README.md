# Module pack CLI

Run from the Core module:

```bash
GOWORK=off go run ./tools/modulepack SOURCE PROXY VERSION RFC3339_TIME
```

The command is a build-time executable, not a runtime library. The entry
point main.go handles arguments, stderr and exit status. The internal
operations/pack package validates the source module, creates its versioned
ZIP, writes .mod and .info metadata, and writes the proxy list.

The Go proxy layout, argument order, and refusal to overwrite an existing
archive are unchanged. Use a clean source directory and a fresh proxy
destination for each release candidate. If packaging fails after the archive
is created, discard the incomplete destination before retrying.

Tests use temporary directories and real Go module archives; no network or
module registry is required.
