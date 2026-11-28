from pathlib import Path
import subprocess

root = Path.cwd()
old = root / "backend/projects/platform/internal/session"
new = root / "backend/libs/platform/session"

def edit(path, before, after):
    p = root / path
    value = p.read_text()
    if before not in value:
        raise SystemExit("missing edit anchor: " + path + ": " + repr(before))
    p.write_text(value.replace(before, after))

for p in sorted(old.rglob("*.go")):
    relative = p.relative_to(old)
    destination = new / relative if relative == Path("session.go") else new / "internal/src" / relative
    destination.parent.mkdir(parents=True, exist_ok=True)
    subprocess.run(["git", "mv", str(p), str(destination)], check=True)

tracked = subprocess.check_output(["git", "ls-files", "-z"]).decode().split("\0")
for name in tracked:
    if not name.endswith(".go"):
        continue
    p = root / name
    if not p.is_file():
        continue
    source = p.read_text()
    if "grovehub/internal/session" not in source:
        continue
    if p.is_relative_to(new):
        source = source.replace("grovehub/internal/session/", "grovehub/libs/platform/session/internal/src/")
    elif "grovehub/internal/session/" in source:
        raise SystemExit("private session subpackage imported by consumer: " + name)
    source = source.replace("grovehub/internal/session", "grovehub/libs/platform/session")
    p.write_text(source)

for path, content in {"backend/libs/platform/go.mod":"module grovehub/libs/platform\n\ngo 1.27.0\n\ntoolchain go1.27.3\n","backend/libs/platform/session/internal/resources/tests/unit-test/session_test.go":"package session_test\n\nimport (\n\t\"context\"\n\t\"errors\"\n\t\"testing\"\n\n\t\"grovehub/libs/platform/session\"\n)\n\nfunc TestPrincipalFacade(t *testing.T) {\n\tctx := context.Background()\n\n\tfor _, value := range []session.Principal{{}, {AccountID: \"account\"}, {SessionID: \"session\"}} {\n\t\tif _, ok := session.PrincipalFrom(session.WithPrincipal(ctx, value)); ok {\n\t\t\tt.Fatalf(\"accepted incomplete principal: %+v\", value)\n\t\t}\n\t}\n\n\tif _, ok := session.PrincipalFrom(ctx); ok {\n\t\tt.Fatal(\"accepted absent principal\")\n\t}\n\n\twant := session.Principal{\n\t\tAccountID: \"account\", SessionID: \"session\", TenantID: \"school\",\n\t\tPersonGID: \"person\", SessionGeneration: 3,\n\t}\n\tif got, ok := session.PrincipalFrom(session.WithPrincipal(ctx, want)); !ok || got != want {\n\t\tt.Fatalf(\"principal not preserved: %+v, %v\", got, ok)\n\t}\n}\n\nfunc TestTenantFacade(t *testing.T) {\n\tctx := context.Background()\n\tif session.ActiveTenantCheck(nil) != nil || session.WithTenantCheck(ctx, nil) != ctx {\n\t\tt.Fatal(\"absent policy should leave isolated context unchanged\")\n\t}\n\n\tif generation, err := session.TenantGeneration(ctx, \"\"); generation != 0 || err != nil {\n\t\tt.Fatalf(\"isolated generation: %d, %v\", generation, err)\n\t}\n\tif err := session.CheckTenantSession(ctx, \"school\", 1); err != nil {\n\t\tt.Fatal(err)\n\t}\n\n\tlookup := func(_ context.Context, id string) (session.TenantState, error) {\n\t\tswitch id {\n\t\tcase \"active\":\n\t\t\treturn session.TenantState{Generation: 3, Active: true}, nil\n\t\tcase \"disabled\":\n\t\t\treturn session.TenantState{Generation: 3}, nil\n\t\tcase \"zero\":\n\t\t\treturn session.TenantState{Active: true}, nil\n\t\tdefault:\n\t\t\treturn session.TenantState{}, errors.New(\"unavailable\")\n\t\t}\n\t}\n\tscoped := session.WithTenantCheck(ctx, session.ActiveTenantCheck(lookup))\n\n\tfor _, id := range []string{\"\", \"disabled\", \"zero\", \"unknown\"} {\n\t\tif err := session.CheckTenant(scoped, id); !errors.Is(err, session.ErrTenantDenied) {\n\t\t\tt.Fatalf(\"%q: expected tenant denial, got %v\", id, err)\n\t\t}\n\t}\n\n\tactive := \"active\"\n\tif err := session.CheckAccountTenant(scoped, &active); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := session.CheckAccountTenant(scoped, nil); !errors.Is(err, session.ErrTenantDenied) {\n\t\tt.Fatalf(\"missing account tenant: %v\", err)\n\t}\n\tif generation, err := session.TenantGeneration(scoped, active); generation != 3 || err != nil {\n\t\tt.Fatalf(\"generation: %d, %v\", generation, err)\n\t}\n\n\tfor _, issued := range []int64{0, 1, 4} {\n\t\tif err := session.CheckAccountGeneration(scoped, &active, issued); !errors.Is(err, session.ErrTenantDenied) {\n\t\t\tt.Fatalf(\"accepted incorrect generation %d: %v\", issued, err)\n\t\t}\n\t}\n\tif err := session.CheckAccountGeneration(scoped, &active, 3); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := session.CheckAccountGeneration(scoped, nil, 3); !errors.Is(err, session.ErrTenantDenied) {\n\t\tt.Fatalf(\"missing account generation tenant: %v\", err)\n\t}\n\tif err := session.CheckTenantSession(scoped, \"disabled\", 3); !errors.Is(err, session.ErrTenantDenied) {\n\t\tt.Fatalf(\"disabled session: %v\", err)\n\t}\n\n\tfor _, generation := range []int64{0, -1} {\n\t\tinvalid := session.WithTenantCheck(ctx, func(context.Context, string) (int64, error) {\n\t\t\treturn generation, nil\n\t\t})\n\t\tif _, err := session.TenantGeneration(invalid, active); !errors.Is(err, session.ErrTenantDenied) {\n\t\t\tt.Fatalf(\"accepted nonpositive generation %d: %v\", generation, err)\n\t\t}\n\t}\n}\n\nfunc TestSessionCookieContract(t *testing.T) {\n\tgot := []string{session.AccessCookie, session.RefreshCookie, session.CSRFName, session.ResetCookie, session.SetupCookie}\n\twant := []string{\"__Host-gh_access\", \"__Secure-gh_refresh\", \"gh_csrf\", \"__Secure-gh_reset\", \"__Secure-gh_setup\"}\n\tfor index := range want {\n\t\tif got[index] != want[index] {\n\t\t\tt.Fatalf(\"cookie %d = %q, want %q\", index, got[index], want[index])\n\t\t}\n\t}\n}\n","backend/libs/platform/session/internal/resources/tests/unit-test/operations/principal/principal_test.go":"package principal_test\n\nimport (\n\t\"context\"\n\t\"testing\"\n\n\t\"grovehub/libs/platform/session/internal/src/contracts\"\n\t\"grovehub/libs/platform/session/internal/src/operations/principal\"\n)\n\nfunc TestPrincipalContext(t *testing.T) {\n\tctx := context.Background()\n\tif _, ok := principal.PrincipalFrom(ctx); ok {\n\t\tt.Fatal(\"absent principal accepted\")\n\t}\n\tfor _, value := range []contracts.Principal{{AccountID: \"a\"}, {SessionID: \"s\"}, {AccountID: \"a\", SessionID: \"s\"}} {\n\t\tgot, ok := principal.PrincipalFrom(principal.WithPrincipal(ctx, value))\n\t\twant := value.AccountID != \"\" && value.SessionID != \"\"\n\t\tif got != value || ok != want {\n\t\t\tt.Fatalf(\"principal = %+v, %v; expected %+v, %v\", got, ok, value, want)\n\t\t}\n\t}\n}\n","backend/libs/platform/session/internal/resources/tests/unit-test/operations/tenant/tenant_test.go":"package tenant_test\n\nimport (\n\t\"context\"\n\t\"errors\"\n\t\"testing\"\n\n\t\"grovehub/libs/platform/session/internal/src/contracts\"\n\t\"grovehub/libs/platform/session/internal/src/operations/tenant\"\n)\n\nfunc TestTenantPolicy(t *testing.T) {\n\tctx := context.Background()\n\tif tenant.ActiveTenantCheck(nil) != nil || tenant.WithTenantCheck(ctx, nil) != ctx {\n\t\tt.Fatal(\"absent policy changed isolated context\")\n\t}\n\tif n, err := tenant.TenantGeneration(ctx, \"\"); n != 0 || err != nil {\n\t\tt.Fatalf(\"isolated generation: %d, %v\", n, err)\n\t}\n\tif err := tenant.CheckTenantSession(ctx, \"school\", 1); err != nil {\n\t\tt.Fatal(err)\n\t}\n\n\tgeneration := int64(1)\n\tactive := true\n\tcheck := tenant.ActiveTenantCheck(func(_ context.Context, id string) (contracts.TenantState, error) {\n\t\tif id != \"school\" {\n\t\t\treturn contracts.TenantState{}, errors.New(\"unknown tenant\")\n\t\t}\n\t\treturn contracts.TenantState{Generation: generation, Active: active}, nil\n\t})\n\tctx = tenant.WithTenantCheck(ctx, check)\n\tid := \"school\"\n\tif err := tenant.CheckAccountTenant(ctx, &id); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := tenant.CheckAccountGeneration(ctx, &id, 1); err != nil {\n\t\tt.Fatal(err)\n\t}\n\n\tactive = false\n\tgeneration++\n\tif err := tenant.CheckTenantSession(ctx, id, 1); !errors.Is(err, contracts.ErrTenantDenied) {\n\t\tt.Fatalf(\"disabled tenant accepted: %v\", err)\n\t}\n\tactive = true\n\tif err := tenant.CheckTenantSession(ctx, id, 1); !errors.Is(err, contracts.ErrTenantDenied) {\n\t\tt.Fatalf(\"stale credential restored: %v\", err)\n\t}\n\tif err := tenant.CheckTenantSession(ctx, id, 2); err != nil {\n\t\tt.Fatal(err)\n\t}\n\tif err := tenant.CheckTenant(ctx, \"unknown\"); !errors.Is(err, contracts.ErrTenantDenied) {\n\t\tt.Fatalf(\"unknown tenant accepted: %v\", err)\n\t}\n\tif err := tenant.CheckAccountTenant(ctx, nil); !errors.Is(err, contracts.ErrTenantDenied) {\n\t\tt.Fatalf(\"missing tenant accepted: %v\", err)\n\t}\n\tif err := tenant.CheckAccountGeneration(ctx, nil, 2); !errors.Is(err, contracts.ErrTenantDenied) {\n\t\tt.Fatalf(\"missing tenant generation accepted: %v\", err)\n\t}\n\n\tgeneration = 0\n\tif err := tenant.CheckTenant(ctx, id); !errors.Is(err, contracts.ErrTenantDenied) {\n\t\tt.Fatalf(\"zero lookup generation accepted: %v\", err)\n\t}\n\tctx = tenant.WithTenantCheck(ctx, func(context.Context, string) (int64, error) {\n\t\treturn 0, nil\n\t})\n\tif err := tenant.CheckTenant(ctx, id); !errors.Is(err, contracts.ErrTenantDenied) {\n\t\tt.Fatalf(\"zero policy generation accepted: %v\", err)\n\t}\n}\n","backend/libs/platform/session/internal/resources/tests/unit-test/constants/cookies_test.go":"package constants_test\n\nimport (\n\t\"testing\"\n\n\t\"grovehub/libs/platform/session/internal/src/constants\"\n)\n\nfunc TestCookieNames(t *testing.T) {\n\tgot := []string{constants.AccessCookie, constants.RefreshCookie, constants.CSRFName, constants.ResetCookie, constants.SetupCookie}\n\twant := []string{\"__Host-gh_access\", \"__Secure-gh_refresh\", \"gh_csrf\", \"__Secure-gh_reset\", \"__Secure-gh_setup\"}\n\tfor i := range want {\n\t\tif got[i] != want[i] {\n\t\t\tt.Fatalf(\"cookie %d = %q, want %q\", i, got[i], want[i])\n\t\t}\n\t}\n}\n","backend/libs/platform/session/internal/resources/tests/unit-test/contracts/principal_test.go":"package contracts_test\n\nimport (\n\t\"context\"\n\t\"testing\"\n\n\t\"grovehub/libs/platform/session/internal/src/contracts\"\n)\n\ntype verifier struct{ principal contracts.Principal }\n\nfunc (v verifier) VerifyAccess(context.Context, string) (contracts.Principal, error) {\n\treturn v.principal, nil\n}\n\nfunc (v verifier) AuthorizeProtectedRoute(context.Context, string) error {\n\treturn nil\n}\n\nvar (\n\t_ contracts.AccessVerifier = verifier{}\n\t_ contracts.ProtectedRouteAuthorizer = verifier{}\n)\n\nfunc TestPrincipalContracts(t *testing.T) {\n\twant := contracts.Principal{AccountID: \"account\", SessionID: \"session\", TenantID: \"school\"}\n\tvar access contracts.AccessVerifier = verifier{principal: want}\n\tgot, err := access.VerifyAccess(context.Background(), \"secret\")\n\tif err != nil || got != want {\n\t\tt.Fatalf(\"principal contract: %+v, %v\", got, err)\n\t}\n}\n","backend/libs/platform/session/internal/resources/tests/unit-test/contracts/tenant_test.go":"package contracts_test\n\nimport (\n\t\"context\"\n\t\"errors\"\n\t\"testing\"\n\n\t\"grovehub/libs/platform/session/internal/src/contracts\"\n)\n\nfunc TestTenantContracts(t *testing.T) {\n\tvar lookup contracts.TenantStateLookup = func(context.Context, string) (contracts.TenantState, error) {\n\t\treturn contracts.TenantState{Generation: 7, Active: true}, nil\n\t}\n\tvar check contracts.TenantCheck = func(ctx context.Context, id string) (int64, error) {\n\t\tif id == \"\" {\n\t\t\treturn 0, contracts.ErrTenantDenied\n\t\t}\n\t\tstate, err := lookup(ctx, id)\n\t\treturn state.Generation, err\n\t}\n\tif n, err := check(context.Background(), \"school\"); n != 7 || err != nil {\n\t\tt.Fatalf(\"tenant state contract: %d, %v\", n, err)\n\t}\n\tif _, err := check(context.Background(), \"\"); !errors.Is(err, contracts.ErrTenantDenied) {\n\t\tt.Fatalf(\"denial contract: %v\", err)\n\t}\n}\n","backend/libs/platform/tools/coverage/check.sh":"#!/usr/bin/env bash\nset -Eeuo pipefail\n\nmodule_dir=\"$(cd -- \"$(dirname -- \"${BASH_SOURCE[0]}\")/../..\" && pwd)\"\nexport GOWORK=off\ncd -- \"$module_dir\"\n\nprofile=\"${1:-$module_dir/.generated/coverage.out}\"\nmkdir -p -- \"$(dirname -- \"$profile\")\"\n# Include every production package and use all standalone test trees.\ngo test -race -count=1 -covermode=atomic -coverpkg=./session/... -coverprofile=\"$profile\" ./...\ngo tool cover -func=\"$profile\"\nawk '\n  NR == 1 { next }\n  { statements[$1] = $2; if ($3 > 0) covered[$1] = 1 }\n  END {\n    for (block in statements) {\n      total += statements[block]\n      if (covered[block]) hit += statements[block]\n      else missing += statements[block]\n    }\n    printf \"Platform library: %d/%d covered statements, %d uncovered\\n\", hit, total, missing\n    if (total == 0 || missing != 0) exit 1\n  }\n' \"$profile\"\n","backend/libs/platform/docs/engineering/README.md":"# Platform shared library\n\nModule: `grovehub/libs/platform`.\n\nThe Platform-specific shared library owns security contracts and session context shared by the Platform launcher and services. It does not own business services, HTTP routing, persistence, or workload orchestration. Session is its first capability.\n\nImport only `grovehub/libs/platform/session`. The capability's constants, contracts, and operations live under `session/internal/src`; other packages cannot import them directly. Principals must originate from an authenticated adapter, and tenant checks must originate from trusted registry policy. Authentication, cookies, and generation checks retain their existing semantics.\n\nThe module is standalone, uses only the Go standard library, and does not depend on Core or application projects. Its local module path is deliberately not a GitHub URL. Platform's project module uses a checkout-relative `replace`; no workspace is required. Remote module distribution is not configured.\n\nFrom this module:\n\n```sh\nGOWORK=off go build ./...\nGOWORK=off go test ./...\nbash tools/coverage/check.sh\n```\n\nTests are capability-owned and run without annotation generation or an application database. Launcher retains an annotated consumer-contract suite. The independent coverage gate rejects every uncovered production statement; Launcher coverage measures only Launcher-owned production packages.\n"}.items():
    p = root / path
    p.parent.mkdir(parents=True, exist_ok=True)
    p.write_text(content)

# Keep the exact strict Launcher denominator owned by Launcher alone.
g = "backend/libs/core/annotation/internal/src/generators/testing/coverage/generator.go"
edit(g, '''	patterns := []string{item.pattern}
	// Platform's shared session context remains in Launcher coverage even
	// though Go's internal-package rules require a Platform-wide boundary.
	if item.name == "launcher" {
		patterns = append(patterns, "./internal/session/...")
	}
	cmd := exec.Command("go", append([]string{"list", "-f", "{{.ImportPath}}"}, patterns...)...)''', '''	cmd := exec.Command("go", "list", "-f", "{{.ImportPath}}", item.pattern)''')

# Exercise the external module dependency while verifying that Launcher only
# measures its own generated and maintained code.
test = "backend/libs/core/annotation/internal/resources/tests/unit-test/generators/testing/coverage/generator_test.go"
edit(test, "materializes physical paths and includes the private Platform session in Launcher coverage", "materializes physical paths and excludes a separately owned Platform module from Launcher coverage")
edit(test, '''		writeFixture(filepath.Join(root, "internal", "session", "session.go"), "package session\\nfunc Generation() int { return 7 }\\n")''', '''		writeFixture(filepath.Join(root, "platform-lib", "go.mod"), "module grovehub/libs/platform\\n\\ngo 1.20\\n")
		writeFixture(filepath.Join(root, "platform-lib", "session", "session.go"), "package session\\nfunc Generation() int { return 7 }\\n")
		writeFixture(filepath.Join(root, "go.mod"), "module example.com/coveragefixture\\n\\ngo 1.20\\n\\nrequire grovehub/libs/platform v0.0.0\\nreplace grovehub/libs/platform => ./platform-lib\\n")''')
edit(test, "example.com/coveragefixture/internal/session", "grovehub/libs/platform/session")
edit(test, "Expect(total).To(Equal(int64(3)))\n\t\tExpect(covered)", "Expect(total).To(Equal(int64(2)))\n\t\tExpect(covered)")
edit(test, 'Expect(string(data)).To(ContainSubstring("internal/session/session.go"))', 'Expect(string(data)).NotTo(ContainSubstring("grovehub/libs/platform/session"))')

# Application wiring and developer entrypoints.
edit("scripts/lib/paths.sh", 'CORE_DIR="${REPO_ROOT}/backend/libs/core"', 'CORE_DIR="${REPO_ROOT}/backend/libs/core"\nPLATFORM_LIB_DIR="${REPO_ROOT}/backend/libs/platform"\nPLATFORM_LIB_MODULE="grovehub/libs/platform"')
edit("scripts/lib/go.sh", '  (\n    cd -- "${BACKEND_DIR}"', '  (\n    cd -- "${PLATFORM_LIB_DIR}"\n    GOWORK=off go mod tidy\n  )\n  (\n    cd -- "${BACKEND_DIR}"')
p = root / "scripts/backend/test.sh"
p.write_text(p.read_text() + '\n# The Platform-specific shared library is an independent Go module.\n(\n  cd -- "${PLATFORM_LIB_DIR}"\n  GOWORK=off go test ./...\n)\n')
edit("scripts/backend/format.sh", '  integration)\n', '  platform-library)\n    path="${PLATFORM_LIB_DIR}"\n    ;;\n  integration)\n')
edit("scripts/backend/format.sh", 'Expected all, core, foundation, foundation-repository, or integration.', 'Expected all, core, platform-library, foundation, foundation-repository, or integration.')
edit("scripts/backend/format.sh", 'if [[ "$scope" == all ]]; then go run "${SCRIPT_DIR}/groovefmt/main.go" -check "${CORE_DIR}"; fi', 'if [[ "$scope" == all ]]; then\n      go run "${SCRIPT_DIR}/groovefmt/main.go" -check "${CORE_DIR}"\n      go run "${SCRIPT_DIR}/groovefmt/main.go" -check "${PLATFORM_LIB_DIR}"\n    fi')
edit("scripts/backend/format.sh", 'if [[ "$scope" == all ]]; then go run "${SCRIPT_DIR}/groovefmt/main.go" "${CORE_DIR}"; fi', 'if [[ "$scope" == all ]]; then\n      go run "${SCRIPT_DIR}/groovefmt/main.go" "${CORE_DIR}"\n      go run "${SCRIPT_DIR}/groovefmt/main.go" "${PLATFORM_LIB_DIR}"\n    fi')
p = root / "Makefile"
p.write_text(p.read_text() + '\n.PHONY: platform-lib-test platform-lib-coverage\nplatform-lib-test:\n\t@cd backend/libs/platform && GOWORK=off go test ./...\n\nplatform-lib-coverage:\n\t@bash ./backend/libs/platform/tools/coverage/check.sh\n')

p = root / "backend/projects/platform/launcher/docs/engineering/session.md"
p.write_text("""# Platform gateway session boundary

Import `grovehub/libs/platform/session` from the standalone module at `backend/libs/platform`. Launcher and Platform services share one security boundary for principals, cookies, and tenant generations. The library has no dependency on Launcher or a business service.

Only the capability facade is public; contracts, constants and operations are private to `session/internal/src`. The Platform project requires the library with a checkout-relative module replacement.

Never derive principals or trusted tenant checks from HTTP input. Authentication and tenant-generation semantics are unchanged. The shared library has an independent strict coverage gate, while Launcher retains annotated consumer-contract tests.
""")
edit("backend/projects/platform/launcher/docs/engineering/architecture.md", "session/ is the only public session contract boundary. All runtime logic is under internal/src, tests under internal/resources/tests. Session types remain public so code generation cannot leak a Go internal type into another service's mock.", "shared security contracts are imported from the standalone `grovehub/libs/platform/session` capability. All Launcher runtime logic is under internal/src, tests under internal/resources/tests. Session library internals are private to their owning capability.")
edit("backend/projects/platform/launcher/docs/engineering/architecture.md", "includes all Launcher production packages and the shared Platform-private session capability, and rejects every uncovered statement.", "includes all Launcher-owned production packages and rejects every uncovered statement. The standalone Platform library has its own strict coverage gate.")
edit("docs/engineering/backend/libraries.md", "The Platform launcher gateway owns its verified session principal, AccessVerifier and school cookie names under `backend/projects/platform/launcher/session`, backed by gateway contracts, constants and operations.", "The independently owned Platform library at `backend/libs/platform` (module `grovehub/libs/platform`) exposes its verified session principal, AccessVerifier and cookie names through `grovehub/libs/platform/session`. Its contracts, constants and operations live under the capability's private `internal/src`. Platform services and Launcher import the facade; Core does not depend on it.")
edit("backend/deployment/README.md", "The Platform, School, Internal Tool and Core modules build with GOWORK=off.", "The Platform, School, Internal Tool, Core and Platform shared-library modules build with GOWORK=off. The Platform project requires `grovehub/libs/platform` using a checkout-relative replacement for `backend/libs/platform`.")

# Add independent library build, formatting, dependency and exact coverage CI.
w = ".github/workflows/backend-annotations.yml"
edit(w, 'echo "Use Platform internal/session."', 'echo "Use grovehub/libs/platform/session."')
edit(w, '          test ! -d backend/projects/platform/libs\n', '''          test ! -d backend/projects/platform/libs
          test ! -d backend/projects/platform/internal/session
          if git grep -nF '"grovehub/internal/session' -- '*.go'; then
            echo "Use grovehub/libs/platform/session."
            exit 1
          fi
''')
edit(w, 'jobs:\n  annotation-import-audit:', '''jobs:
  platform-library:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v6
        with:
          go-version: "1.27.3"
          cache: false
      - name: Verify standalone Platform library boundary
        working-directory: backend/libs/platform
        shell: bash
        run: |
          test "$(go list -m)" = "grovehub/libs/platform"
          test "$(go list -m all)" = "grovehub/libs/platform"
          test -z "$(gofmt -l session)"
          go build ./...
          go vet ./...
      - name: Enforce 100 percent Platform library statement coverage
        run: bash backend/libs/platform/tools/coverage/check.sh "$RUNNER_TEMP/platform-library-coverage.out"
      - uses: actions/upload-artifact@v4
        with:
          name: platform-library-coverage
          path: ${{ runner.temp }}/platform-library-coverage.out
          if-no-files-found: error
  annotation-import-audit:''')
edit("backend/projects/platform/launcher/docs/engineering/architecture.md", "Generated routes and facade methods are included.", "Generated routes and facade methods are included. Run `make platform-lib-coverage` independently for the shared library.")

