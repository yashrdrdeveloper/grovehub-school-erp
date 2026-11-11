package routes_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	foundationroutes "grovehub/services/foundation/routes"
)

// The route adapter is intentionally inert: these tests exercise only
// generated dispatch, independent of business services and databases.
type personHandler struct{}

func (personHandler) Create(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, "Create")
}
func (personHandler) GetByGID(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "GetByGID:"+request.PathValue("gid"))
}
func (personHandler) Update(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, "Update")
}
func (personHandler) Search(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, "Search")
}
func (personHandler) Deactivate(writer http.ResponseWriter, _ *http.Request) {
	fmt.Fprint(writer, "Deactivate")
}

var _ = Describe("generated Foundation route dispatch", func() {
	var mux *http.ServeMux

	BeforeEach(func() {
		mux = http.NewServeMux()
		foundationroutes.RegisterHTTP(mux, personHandler{})
	})

	It("dispatches all five Person methods under the protected service prefix", func() {
		cases := []struct {
			method, path, body string
		}{
			{"POST", "/foundation/protected/persons", "Create"},
			{"GET", "/foundation/protected/persons/01900000-0000-7000-8000-000000000001", "GetByGID:01900000-0000-7000-8000-000000000001"},
			{"PATCH", "/foundation/protected/persons/01900000-0000-7000-8000-000000000001", "Update"},
			{"POST", "/foundation/protected/persons/search", "Search"},
			{"PATCH", "/foundation/protected/persons/01900000-0000-7000-8000-000000000001/deactivate", "Deactivate"},
		}
		for _, testcase := range cases {
			request := httptest.NewRequest(testcase.method, testcase.path, strings.NewReader("{}"))
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, request)
			Expect(response.Code).To(Equal(http.StatusOK), testcase.method+" "+testcase.path)
			Expect(response.Body.String()).To(Equal(testcase.body), testcase.method+" "+testcase.path)
		}
	})

	It("does not retain the old unclassified Foundation URLs", func() {
		for _, path := range []string{
			"/foundation/persons",
			"/foundation/persons/search",
			"/foundation/persons/01900000-0000-7000-8000-000000000001",
		} {
			response := httptest.NewRecorder()
			mux.ServeHTTP(response, httptest.NewRequest(http.MethodGet, path, nil))
			Expect(response.Code).To(Equal(http.StatusNotFound), path)
		}
	})
})
