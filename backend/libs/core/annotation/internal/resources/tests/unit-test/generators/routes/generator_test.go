package routes_test

import (
	"os"
	"path/filepath"
	"strings"
	. "github.com/onsi/ginkgo/v2"

	"github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/annotation/internal/src/generators/routes"
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

func writeTestFile(t GinkgoTInterface, filename, value string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(filename), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filename, []byte(value), 0644); err != nil {
		t.Fatal(err)
	}
}

func testBackend(t GinkgoTInterface) (string, string) {
	t.Helper()
	backend := t.TempDir()
	writeTestFile(t, filepath.Join(backend, "go.mod"), "module grovehub\n")
	owner := filepath.Join(backend, "services", "foundation")
	writeTestFile(t, filepath.Join(owner, "routes", "http.go"), serviceRoutes)
	writeTestFile(t, filepath.Join(owner, "controller", "protected", "person.go"), personController)
	return backend, filepath.Join(owner, ".generated", "routes", "zz_routes.gen.go")
}

func generatedRoutes(t GinkgoTInterface, backend, output string) string {
	t.Helper()
	if err := routes.GenerateAll(backend); err != nil {
		t.Fatal(err)
	}
	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatal(err)
	}
	return string(data)
}

var _ = Describe("GenerateAll", func() {
	It("Route Access Tiers And Service Prefix", func() {
		t := GinkgoT()
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

	})
})

var _ = Describe("GenerateAll", func() {
	It("Service Prefix Is Owned By Route Composition", func() {
		t := GinkgoT()
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

	})
})

var _ = Describe("GenerateAll", func() {
	It("Legacy Uncategorized Routes Stay Compatible", func() {
		t := GinkgoT()
	backend := t.TempDir()
	writeTestFile(t, filepath.Join(backend, "go.mod"), "module grovehub\n")
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

	})
})

var _ = Describe("GenerateAll", func() {
	It("Categorized Controllers Require Service Prefix", func() {
		t := GinkgoT()
	backend, _ := testBackend(t)
	filename := filepath.Join(backend, "services", "foundation", "routes", "http.go")
	writeTestFile(t, filename, `package routes
import "net/http"
func RegisterHTTP(mux *http.ServeMux) {}
`)
	if err := routes.GenerateAll(backend); err == nil || !strings.Contains(err.Error(), "requires a service @Routes") {
		t.Fatalf("expected missing service prefix error, got %v", err)
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Service Prefix Requires Categorized Controllers", func() {
		t := GinkgoT()
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
	if err := routes.GenerateAll(backend); err == nil || !strings.Contains(err.Error(), "must live in public/") {
		t.Fatalf("expected uncategorized controller error, got %v", err)
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Generator Rejects Duplicate Routes Within Tier", func() {
		t := GinkgoT()
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
	if err := routes.GenerateAll(backend); err == nil || !strings.Contains(err.Error(), "duplicate route") {
		t.Fatalf("expected duplicate route error, got %v", err)
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Service Prefix Annotation Only On Register HTTP", func() {
		t := GinkgoT()
	backend, _ := testBackend(t)
	filename := filepath.Join(backend, "services", "foundation", "routes", "http.go")
	writeTestFile(t, filename, strings.Replace(serviceRoutes, "func RegisterHTTP", "func Configure", 1))
	if err := routes.GenerateAll(backend); err == nil || !strings.Contains(err.Error(), "RegisterHTTP function") {
		t.Fatalf("expected misplaced service @Routes error, got %v", err)
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Generator Rejects Malformed Service Prefix", func() {
		t := GinkgoT()
	backend, _ := testBackend(t)
	filename := filepath.Join(backend, "services", "foundation", "routes", "http.go")
	writeTestFile(t, filename, strings.Replace(serviceRoutes, "/foundation", "/foundation//admin", 1))
	if err := routes.GenerateAll(backend); err == nil || !strings.Contains(err.Error(), "invalid route path") {
		t.Fatalf("expected invalid prefix error, got %v", err)
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Resource First Controller Package", func() {
		t := GinkgoT()
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

	})
})

var _ = Describe("GenerateAll", func() {
	It("Route Limit Generates Exact Method And Pattern", func() {
		t := GinkgoT()
	backend, output := testBackend(t)
	controller := filepath.Join(backend, "services", "foundation", "controller", "protected", "person.go")
	writeTestFile(t, controller, strings.Replace(personController,
		"// @POST /search",
		"// @POST /search\n\t// @RouteLimit 30", 1))
	content := generatedRoutes(t, backend, output)
	if !strings.Contains(content, "`POST /foundation/protected/persons/search`") &&
		!strings.Contains(content, "\"POST /foundation/protected/persons/search\": 30") {
		t.Fatal("annotated method policy missing from generated routes")
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Route Limit Rejects Malformed Limits", func() {
		t := GinkgoT()
	for _, value := range []string{"0", "-2", "abc", "100001"} {
		By(value)
		{
			backend, _ := testBackend(t)
			controller := filepath.Join(backend, "services", "foundation", "controller", "protected", "person.go")
			writeTestFile(t, controller, strings.Replace(personController,
				"// @POST /search",
				"// @POST /search\n\t// @RouteLimit "+value, 1))
			if err := routes.GenerateAll(backend); err == nil || !strings.Contains(err.Error(), "@RouteLimit") {
				t.Fatalf("expected invalid limit rejection; got %v", err)
			}
		}
	}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Routes Use Target Module Path", func() {
		t := GinkgoT()
	backend, output := testBackend(t)
	writeTestFile(t, filepath.Join(backend, "go.mod"), "module example.org/internal-platform\n")
	content := generatedRoutes(t, backend, output)
	if !strings.Contains(content, "example.org/internal-platform/services/foundation/controller/protected") {
		t.Fatal("generated controller import must use target module")
	}
	if strings.Contains(content, "grovehub/services/") {
		t.Fatal("generated controller import contains the school module")
	}

	})
})
