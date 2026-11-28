# Core build tools

These tools are development and release executables, not runtime capabilities.
There is deliberately no tools.go or importable umbrella package.

- **clientgen/generate.sh**: standalone Docker-based OpenAPI client generator.
  The release contract bundle exports this exact script, so its entry point
  must not depend on sibling source files.
- **modulepack/main.go**: Go CLI for packaging a release candidate in standard
  module proxy format. Reusable command logic is private under
  internal/src/operations/pack, with mirrored tests.

Both existing command paths remain stable. Do not import CLI packages from
runtime service code or move service-owned contracts into these tools.

```bash
cd backend/libs/core
GOWORK=off go test ./tools/modulepack/...
bash tools/clientgen/internal/resources/tests/unit-test/generate_test.sh
```
