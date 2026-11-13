package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const personController = `package protected

import "net/http"

// @Routes /persons
type PersonController interface {
	// @POST /
	Create(http.ResponseWriter, *http.Request)
	// @GET /{gid}
	GetByGID(http.ResponseWriter, *http.Request)
	// @PATCH /{gid}
	Update(http.ResponseWriter, *http.Request)
	// @POST /search
	Search(http.ResponseWriter, *http.Request)
	// @PATCH /{gid}/deactivate
	Deactivate(http.ResponseWriter, *http.Request)
}
`

const serviceRoutes = `package routes
import "net/http"

// @Routes /foundation
func RegisterHTTP(mux *http.ServeMux) {}
`

func writeTestFile(t *testing.T, filename, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte(value), 0644); err != nil {
		t.Fatal(err)
	}
}

func testBackend(t *testing.T) (string, string) {
	t.Helper()
	backend := t.TempDir()
	owner := filepath.Join(backend, "services", "foundation")
	writeTestFile(t, filepath.Join(owner, "routes", "http.go"), serviceRoutes)
	writeTestFile(t, filepath.Join(owner, "controller", "protected", "person.go"), personController)
	return backend, filepath.Join(owner, ".generated", "routes", "zz_routes.gen.go")
}

func generatedRoutes(t *testing.T, backend, output string) string {
	t.Helper()
	if err := generateAll(backend); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

func TestRouteAccessTiersAndServicePrefix(t *testing.T) {
	backend, output := testBackend(t)
	owner := filepath.Join(backend, "services", "foundation")
	writeTestFile(t, filepath.Join(owner, "controller", "public", "health.go"), `package public
import "net/http"
// @Routes /health
type HealthController interface {
	// @GET /
	Status(http.ResponseWriter, *http.Request)
}
`)
	writeTestFile(t, filepath.Join(owner, "controller", "restricted", "maintenance.go"), `package restricted
import "net/http"
// @Routes /maintenance
type MaintenanceController interface {
	// @POST /repair
	Repair(http.ResponseWriter, *http.Request)
}
`)
	content := generatedRoutes(t, backend, output)
	for _, want := range []string{
		"POST /foundation/protected/persons",
		"GET /foundation/protected/persons/{gid}",
		"PATCH /foundation/protected/persons/{gid}",
		"POST /foundation/protected/persons/search",
		"PATCH /foundation/protected/persons/{gid}/deactivate",
		"GET /foundation/public/health",
		"POST /foundation/restricted/maintenance/repair",
		"func registerProtectedPersonRoutes(",
		"func registerPublicHealthRoutes(",
		"func registerRestrictedMaintenanceRoutes(",
		"protected.PersonController",
		"public.HealthController",
		"restricted.MaintenanceController",
	} {
		if !strings.Contains(content, want) {
			t.Errorf("generated routes missing %q", want)
		}
	}
	if strings.Contains(content, "GET /foundation/persons/") {
		t.Fatal("generated routes retain old unclassified person path")
	}
	if next := generatedRoutes(t, backend, output); next != content {
		t.Fatal("generator is not deterministic")
	}
}

func TestServicePrefixIsOwnedByRouteComposition(t *testing.T) {
	backend, output := testBackend(t)
	_ = generatedRoutes(t, backend, output)
	filename := filepath.Join(backend, "services", "foundation", "routes", "http.go")
	writeTestFile(t, filename, strings.Replace(serviceRoutes, "/foundation", "/school", 1))
	content := generatedRoutes(t, backend, output)
	if !strings.Contains(content, "GET /school/protected/persons/{gid}") {
		t.Fatal("changing service @Routes must update every generated endpoint")
	}
	if strings.Contains(content, "GET /foundation/") || strings.Contains(content, "POST /foundation/") || strings.Contains(content, "PATCH /foundation/") {
		t.Fatal("stale service prefix retained")
	}
}

func TestLegacyUncategorizedRoutesStayCompatible(t *testing.T) {
	backend := t.TempDir()
	owner := filepath.Join(backend, "services", "identity")
	writeTestFile(t, filepath.Join(owner, "routes", "http.go"), `package routes
import "net/http"
func RegisterHTTP(mux *http.ServeMux) {}
`)
	writeTestFile(t, filepath.Join(owner, "controller", "auth.go"), `package controller
import "net/http"
// @Routes /auth
type AuthController interface {
	// @POST /register
	Register(http.ResponseWriter, *http.Request)
}
`)
	output := filepath.Join(owner, ".generated", "routes", "zz_routes.gen.go")
	content := generatedRoutes(t, backend, output)
	if !strings.Contains(content, "POST /auth/register") || !strings.Contains(content, "registerAuthRoutes(") {
		t.Fatal("legacy Identity routes must remain unchanged")
	}
	if strings.Contains(content, "POST /identity/") {
		t.Fatal("legacy service prefix must not be inferred")
	}
}

func TestCategorizedControllersRequireServicePrefix(t *testing.T) {
	backend, _ := testBackend(t)
	filename := filepath.Join(backend, "services", "foundation", "routes", "http.go")
	writeTestFile(t, filename, `package routes
import "net/http"
func RegisterHTTP(mux *http.ServeMux) {}
`)
	if err := generateAll(backend); err == nil || !strings.Contains(err.Error(), "requires a service @Routes") {
		t.Fatalf("expected missing service prefix error, got %v", err)
	}
}

func TestServicePrefixRequiresCategorizedControllers(t *testing.T) {
	backend, _ := testBackend(t)
	filename := filepath.Join(backend, "services", "foundation", "controller", "legacy.go")
	writeTestFile(t, filename, `package controller
import "net/http"
// @Routes /legacy
type LegacyController interface {
	// @GET /
	Get(http.ResponseWriter, *http.Request)
}
`)
	if err := generateAll(backend); err == nil || !strings.Contains(err.Error(), "must live in public/") {
		t.Fatalf("expected uncategorized controller error, got %v", err)
	}
}

func TestGeneratorRejectsDuplicateRoutesWithinTier(t *testing.T) {
	backend, _ := testBackend(t)
	filename := filepath.Join(backend, "services", "foundation", "controller", "protected", "duplicate.go")
	writeTestFile(t, filename, `package protected
import "net/http"
// @Routes /persons
type DuplicateController interface {
	// @POST /
	Create(http.ResponseWriter, *http.Request)
}
`)
	if err := generateAll(backend); err == nil || !strings.Contains(err.Error(), "duplicate route") {
		t.Fatalf("expected duplicate route error, got %v", err)
	}
}

func TestServicePrefixAnnotationOnlyOnRegisterHTTP(t *testing.T) {
	backend, _ := testBackend(t)
	filename := filepath.Join(backend, "services", "foundation", "routes", "http.go")
	writeTestFile(t, filename, strings.Replace(serviceRoutes, "func RegisterHTTP", "func Configure", 1))
	if err := generateAll(backend); err == nil || !strings.Contains(err.Error(), "RegisterHTTP function") {
		t.Fatalf("expected misplaced service @Routes error, got %v", err)
	}
}

func TestGeneratorRejectsMalformedServicePrefix(t *testing.T) {
	backend, _ := testBackend(t)
	filename := filepath.Join(backend, "services", "foundation", "routes", "http.go")
	writeTestFile(t, filename, strings.Replace(serviceRoutes, "/foundation", "/foundation//admin", 1))
	if err := generateAll(backend); err == nil || !strings.Contains(err.Error(), "invalid route path") {
		t.Fatalf("expected invalid prefix error, got %v", err)
	}
}

func TestResourceFirstControllerPackage(t *testing.T) {
	backend,output:=testBackend(t)
	owner:=filepath.Join(backend,"services","foundation")
	legacy:=filepath.Join(owner,"controller","protected","person.go")
	if err:=os.Remove(legacy);err!=nil {t.Fatal(err)}
	writeTestFile(t,filepath.Join(owner,"controller","person","protected","person.go"),personController)
	content:=generatedRoutes(t,backend,output)
	for _,required:=range []string{
		"personprotected \"grovehub/services/foundation/controller/person/protected\"",
		"func registerProtectedPersonRoutes(",
		"c personprotected.PersonController",
		"GET /foundation/protected/persons/{gid}",
	}{
		if !strings.Contains(content,required){t.Errorf("missing resource route %q:\n%s",required,content)}
	}
}
