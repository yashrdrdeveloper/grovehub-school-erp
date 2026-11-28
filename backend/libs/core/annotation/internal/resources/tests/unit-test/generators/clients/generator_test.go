package clients_test

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"os"
	"path/filepath"
	"strings"

	"grovehub/libs/core/annotation/internal/src/generators/clients"
)

const sourceClient = `package client
import (
"context"
)
// @Client personClient
type Client interface {
	// @RequestLine POST /protected/persons
	CreatePerson(ctx context.Context, request CreateDTO) (Person, bool, error)
	// @RequestLine GET /protected/persons/{gid}
	GetPersonByGID(parent context.Context, request ByGID) (Person, error)
	// @RequestLine PATCH /protected/persons/{gid}
	UpdatePerson(ctx context.Context, request UpdateDTO) (Person, error)
	// @RequestLine POST /protected/persons/search
	SearchPersons(ctx context.Context, request SearchDTO) (SearchResult, error)
	// @RequestLine PATCH /protected/persons/{gid}/deactivate
	DeactivatePerson(ctx context.Context, request ByGID) (Person, error)
}
`

func fixture(t GinkgoTInterface, clientContent, prefix string) (string, string, string) {
	t.Helper()
	root := t.TempDir()
	owner := filepath.Join(root, "services", "foundation")
	client := filepath.Join(owner, "client", "person.go")
	routes := filepath.Join(owner, ".generated", "routes", "zz_routes.gen.go")
	output := filepath.Join(owner, ".generated", "client", "impl", "zz_clients.gen.go")
	for path, content := range map[string]string{
		client: clientContent,
		routes: `package routes
import "net/http"
func Register(mux *http.ServeMux) {
mux.Handle("POST ` + prefix + `/protected/persons", nil)
mux.Handle("GET ` + prefix + `/protected/persons/{gid}", nil)
mux.Handle("PATCH ` + prefix + `/protected/persons/{gid}", nil)
mux.Handle("POST ` + prefix + `/protected/persons/search", nil)
mux.Handle("PATCH ` + prefix + `/protected/persons/{gid}/deactivate", nil)
}
`,
	} {
		Expect(os.MkdirAll(filepath.Dir(path), 0755)).To(Succeed())
		Expect(os.WriteFile(path, []byte(content), 0644)).To(Succeed())
	}
	return root, client, output
}

var _ = Describe("GenerateAll", func() {
	It("Request Line Generates Full Controller Dispatch", func() {
		t := GinkgoT()
		root, _, out := fixture(t, sourceClient, "/foundation")
		Expect(clients.GenerateAll(root)).To(Succeed())
		data, err := os.ReadFile(out)
		Expect(err).NotTo(HaveOccurred())
		generated := string(data)
		for _, expected := range []string{
			"httpendpoint.CallStruct[Person](ctx, c.handler, \"POST\", \"/foundation/protected/persons\", request)",
			"httpendpoint.CallStruct[Person](parent, c.handler, \"GET\", \"/foundation/protected/persons/{gid}\", request)",
			"httpendpoint.CallStruct[Person](ctx, c.handler, \"PATCH\", \"/foundation/protected/persons/{gid}\", request)",
			"httpendpoint.CallStruct[SearchResult](ctx, c.handler, \"POST\", \"/foundation/protected/persons/search\", request)",
			"status == 201 && err == nil",
		} {
			Expect(generated).To(ContainSubstring(expected))
		}
		Expect(clients.GenerateAll(root)).To(Succeed())
		again, _ := os.ReadFile(out)
		Expect(string(again)).To(Equal(generated))

	})
})

var _ = Describe("GenerateAll", func() {
	It("Service Prefix Changes Only In Generated Routes", func() {
		t := GinkgoT()
		root, _, out := fixture(t, sourceClient, "/school")
		Expect(clients.GenerateAll(root)).To(Succeed())
		data, _ := os.ReadFile(out)
		Expect(string(data)).To(ContainSubstring("\"/school/protected/persons/{gid}\""))
		Expect(string(data)).NotTo(ContainSubstring("\"/foundation/"))

	})
})

var _ = Describe("GenerateAll", func() {
	It("Request Line Missing Controller Route Fails Closed", func() {
		t := GinkgoT()
		root, _, _ := fixture(t, strings.Replace(sourceClient, "GET /protected/persons/{gid}", "DELETE /protected/persons/{gid}", 1), "/foundation")
		{
			err := clients.GenerateAll(root)
			Expect(err).To(MatchError(ContainSubstring("no matching generated controller route")))
		}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Request Line Rejects Ambiguous Matches", func() {
		t := GinkgoT()
		root, _, _ := fixture(t, sourceClient, "/foundation")
		routes := filepath.Join(root, "services", "foundation", ".generated", "routes", "zz_routes.gen.go")
		data, _ := os.ReadFile(routes)
		src := strings.Replace(string(data), "mux.Handle(\"GET /foundation/protected/persons/{gid}\", nil)", "mux.Handle(\"GET /foundation/protected/persons/{gid}\", nil)\nmux.Handle(\"GET /new/protected/persons/{gid}\", nil)", 1)
		Expect(os.WriteFile(routes, []byte(src), 0644)).To(Succeed())
		{
			err := clients.GenerateAll(root)
			Expect(err).To(MatchError(ContainSubstring("matches multiple routes")))
		}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Request Line Rejects Invalid Signatures", func() {
		t := GinkgoT()
		cases := []struct{ source, errorFragment string }{
			{strings.Replace(sourceClient, "(Person, bool, error)", "(Person, string, error)", 1), "second result must be bool"},
			{strings.Replace(sourceClient, "request ByGID", "request []ByGID", 1), "second @RequestLine parameter must be a named request struct"},
			{strings.Replace(sourceClient, "// @RequestLine GET /protected/persons/{gid}", "// @RequestLine GET /protected/persons/{gid}\n\t// @Proxy service.Get", 1), "exactly one of"},
			{strings.Replace(sourceClient, "// @RequestLine GET /protected/persons/{gid}", "// @RequestLine GET /protected/persons/{gid}\n\t// @RequestLine GET /protected/persons/{gid}", 1), "duplicate @RequestLine"},
		}
		for _, item := range cases {
			root, _, _ := fixture(t, item.source, "/foundation")
			{
				err := clients.GenerateAll(root)
				Expect(err).To(MatchError(ContainSubstring(item.errorFragment)))
			}
		}

	})
})

var _ = Describe("GenerateAll", func() {
	It("Legacy Proxy Methods Remain Supported", func() {
		t := GinkgoT()
		client := `package client
import "context"
// @Client existingClient
type Client interface {
	// @Proxy personService.Get
	Get(ctx context.Context, id string) (string,error)
}
`
		root, _, out := fixture(t, client, "/foundation")
		Expect(clients.GenerateAll(root)).To(Succeed())
		b, _ := os.ReadFile(out)
		Expect(string(b)).To(ContainSubstring("return c.personService.Get(ctx, id)"))
		Expect(string(b)).NotTo(ContainSubstring("httpendpoint"))

	})
})

var _ = Describe("GenerateAll", func() {
	It("Resource First Client Generation", func() {
		t := GinkgoT()
		root, original, _ := fixture(t, sourceClient, "/foundation")
		owner := filepath.Join(root, "services", "foundation")
		resource := filepath.Join(owner, "client", "person", "person.go")
		Expect(os.MkdirAll(filepath.Dir(resource), 0755)).To(Succeed())
		Expect(os.Rename(original, resource)).To(Succeed())
		Expect(clients.GenerateAll(root)).To(Succeed())
		output := filepath.Join(owner, ".generated", "client", "person", "impl", "zz_clients.gen.go")
		data, err := os.ReadFile(output)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(ContainSubstring("httpendpoint.CallStruct[Person]("))
		{
			_, err := os.Stat(filepath.Join(owner, ".generated", "client", "impl", "zz_clients.gen.go"))
			Expect(os.IsNotExist(err)).To(BeTrue())
		}

	})
})
