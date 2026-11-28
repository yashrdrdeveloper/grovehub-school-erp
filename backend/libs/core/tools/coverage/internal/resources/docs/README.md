# Core coverage policy

Every capability has two independent coverage scopes:

- Internal: all executable statements under internal/src, exercised by that
  capability's internal/resources/tests/unit-test suites.
- Public: all executable statements in the root facade, exercised by the
  external tests/unit-test/<capability> suites.

No combined averages, rounded percentages, missing-package exclusions or
allowlists for uncovered blocks. A scope with only declarations has no
executable statements and is reported as such, not as 100% fabricated coverage.
All tests must pass. CLI entry points are covered separately with instrumented
executables; modulepack implementation and shared test support have their own
scope. The standalone clientgen Bash adapter has command-level regression
tests; Go statement coverage does not measure shell code.

Run a capability locally:

```bash
bash backend/libs/core/tools/coverage/check.sh search /tmp/core-coverage
```

The matrix discovers each runtime capability and runs all scopes on every
Core change. All uncovered source blocks are printed, even where a rounded
percentage might display 100%. Failed jobs publish their available profiles.
