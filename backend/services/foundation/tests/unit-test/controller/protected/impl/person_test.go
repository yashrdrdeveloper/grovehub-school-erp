package impl_test

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	coreerrors "grovehub/libs/core/errors"
	coresearch "grovehub/libs/core/search"
	personerror "grovehub/services/foundation/constants/error"
	controllerimpl "grovehub/services/foundation/controller/protected/impl"
	foundationrequest "grovehub/services/foundation/dto/request"
	foundationresponse "grovehub/services/foundation/dto/response"
	foundationroutes "grovehub/services/foundation/routes"
	validatorimpl "grovehub/services/foundation/validators/impl"
)

const goodGID = "018f5d04-315c-75a9-b890-df0d02e787c1"

type fakeService struct {
	calls int
	created bool
	failure error
	lastCreate foundationrequest.CreatePerson
	lastUpdate foundationrequest.UpdatePerson
	lastPath foundationrequest.PersonByGID
}

func (s *fakeService) Create(_ context.Context, input foundationrequest.CreatePerson) (foundationresponse.Person, bool, error) {
	s.calls++
	s.lastCreate = input
	return foundationresponse.Person{GID: goodGID, FirstName: input.FirstName}, s.created, s.failure
}
func (s *fakeService) GetByGID(_ context.Context, input foundationrequest.PersonByGID) (foundationresponse.Person, error) {
	s.calls++
	s.lastPath = input
	return foundationresponse.Person{GID: input.GID}, s.failure
}
func (s *fakeService) Update(_ context.Context, path foundationrequest.PersonByGID, input foundationrequest.UpdatePerson) (foundationresponse.Person, error) {
	s.calls++
	s.lastPath = path
	s.lastUpdate = input
	return foundationresponse.Person{GID: path.GID}, s.failure
}
func (s *fakeService) Search(_ context.Context, _ coresearch.Request) (coresearch.Result[foundationresponse.Person], error) {
	s.calls++
	return coresearch.Result[foundationresponse.Person]{}, s.failure
}
func (s *fakeService) Deactivate(_ context.Context, path foundationrequest.PersonByGID) (foundationresponse.Person, error) {
	s.calls++
	s.lastPath = path
	return foundationresponse.Person{GID: path.GID}, s.failure
}

func fixtureMux(service *fakeService) *http.ServeMux {
	mux := http.NewServeMux()
	foundationroutes.RegisterHTTP(mux, controllerimpl.NewPersonController(service, validatorimpl.NewPersonValidator()))
	return mux
}

func call(mux *http.ServeMux, method, path, body string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(method, path, strings.NewReader(body)))
	return recorder
}

var _ = Describe("typed Person controller pipeline", func() {
	It("normalizes Create body and writes 201 for created, 200 for existing", func() {
		for _, created := range []bool{true, false} {
			service := &fakeService{created: created}
			response := call(fixtureMux(service), http.MethodPost, "/foundation/protected/persons", `{"firstName":"  Arun ","lastName":" Kumar "}`)
			want := http.StatusOK
			if created { want = http.StatusCreated }
			Expect(response.Code).To(Equal(want))
			Expect(service.calls).To(Equal(1))
			Expect(service.lastCreate.FirstName).To(Equal("Arun"))
			Expect(service.lastCreate.LastName).To(Equal("Kumar"))
		}
	})

	It("rejects an invalid GID before executing the service", func() {
		service := &fakeService{}
		response := call(fixtureMux(service), http.MethodGet, "/foundation/protected/persons/not-a-gid", "")
		Expect(response.Code).To(Equal(http.StatusBadRequest))
		Expect(response.Body.String()).To(ContainSubstring("FOUNDATION_PERSON_INVALID_GID"))
		Expect(service.calls).To(BeZero())
	})

	It("normalizes PATCH body and forwards the path DTO", func() {
		service := &fakeService{}
		response := call(fixtureMux(service), http.MethodPatch, "/foundation/protected/persons/"+goodGID, `{"firstName":"  Arun "}`)
		Expect(response.Code).To(Equal(http.StatusOK))
		Expect(service.calls).To(Equal(1))
		Expect(service.lastPath.GID).To(Equal(goodGID))
		Expect(service.lastUpdate.FirstName.Set).To(BeTrue())
		Expect(*service.lastUpdate.FirstName.Value).To(Equal("Arun"))
	})

	It("preserves domain status and safe metadata from service errors", func() {
		service := &fakeService{failure: coreerrors.NewErrorWithMetadata(personerror.PersonContactConflict, map[string]string{"field":"email"})}
		response := call(fixtureMux(service), http.MethodGet, "/foundation/protected/persons/"+goodGID, "")
		Expect(response.Code).To(Equal(http.StatusConflict))
		Expect(response.Body.String()).To(ContainSubstring("FOUNDATION_PERSON_CONTACT_CONFLICT"))
		Expect(response.Body.String()).To(ContainSubstring("email"))
	})

	It("does not reveal unexpected service failures", func() {
		service := &fakeService{failure: errors.New("postgres secret error")}
		response := call(fixtureMux(service), http.MethodGet, "/foundation/protected/persons/"+goodGID, "")
		Expect(response.Code).To(Equal(http.StatusInternalServerError))
		Expect(response.Body.String()).To(ContainSubstring("FOUNDATION_PERSON_READ_FAILED"))
		Expect(response.Body.String()).NotTo(ContainSubstring("postgres secret error"))
	})
})
