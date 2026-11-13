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
	headercontext "grovehub/libs/core/headercontext"
	annotationdto "grovehub/libs/core/annotation/dto"
	foundationclientimpl "grovehub/services/foundation/client/impl"
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
	lastCtx context.Context
}

func (s *fakeService) Create(ctx context.Context, input foundationrequest.CreatePerson) (foundationresponse.Person, bool, error) {
	s.calls++
	s.lastCreate = input
	s.lastCtx = ctx
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

	It("routes generated @RequestLine CreatePerson through the controller and derives created from HTTP 201", func() {
		for _, created := range []bool{true, false} {
			service := &fakeService{created: created}
			client := foundationclientimpl.NewPersonClient(fixtureMux(service))
			ctx := headercontext.WithValues(context.Background(), headercontext.Values{IdempotencyKey:"an-internal-key"})
			person, wasCreated, err := client.CreatePerson(ctx, foundationrequest.CreatePerson{
				Person: foundationrequest.Person{FirstName:"  Arun ", LastName:"  Kumar "},
			})
			Expect(err).NotTo(HaveOccurred())
			Expect(wasCreated).To(Equal(created))
			Expect(person.FirstName).To(Equal("Arun"))
			Expect(service.calls).To(Equal(1))
			Expect(headercontext.IdempotencyKey(service.lastCtx)).To(Equal("an-internal-key"))
		}
	})

	It("maps generated Get, Update, Search and Deactivate clients through controller methods", func() {
		service := &fakeService{}
		client := foundationclientimpl.NewPersonClient(fixtureMux(service))
		path := foundationrequest.PersonByGID{GID: goodGID}
		person, err := client.GetPersonByGID(context.Background(), path)
		Expect(err).NotTo(HaveOccurred())
		Expect(person.GID).To(Equal(goodGID))

		firstName := "  Name  "
		person, err = client.UpdatePerson(context.Background(), path, foundationrequest.UpdatePerson{
			FirstName: annotationdto.Patch[string]{Set:true,Value:&firstName},
		})
		Expect(err).NotTo(HaveOccurred())
		Expect(person.GID).To(Equal(goodGID))
		Expect(*service.lastUpdate.FirstName.Value).To(Equal("Name"))
		Expect(service.lastPath.GID).To(Equal(goodGID))

		_, err = client.SearchPersons(context.Background(), coresearch.Request{})
		Expect(err).NotTo(HaveOccurred())

		person, err = client.DeactivatePerson(context.Background(), path)
		Expect(err).NotTo(HaveOccurred())
		Expect(person.GID).To(Equal(goodGID))
		Expect(service.calls).To(Equal(4))
	})

	It("propagates controller validation errors as typed client application errors", func() {
		service := &fakeService{}
		client := foundationclientimpl.NewPersonClient(fixtureMux(service))
		_, err := client.GetPersonByGID(context.Background(), foundationrequest.PersonByGID{GID:"bad-id"})
		var applicationError *coreerrors.Error
		Expect(errors.As(err, &applicationError)).To(BeTrue())
		Expect(applicationError.Definition().Code()).To(Equal(personerror.InvalidPersonGID.Code()))
		Expect(applicationError.Definition().HTTPStatus()).To(Equal(http.StatusBadRequest))
		Expect(service.calls).To(BeZero())

		_, _, err = client.CreatePerson(context.Background(), foundationrequest.CreatePerson{
			Person: foundationrequest.Person{FirstName:"   ",LastName:"Valid"},
		})
		Expect(errors.As(err, &applicationError)).To(BeTrue())
		Expect(applicationError.Definition().Code()).To(Equal(personerror.InvalidPersonRequest.Code()))
		Expect(service.calls).To(BeZero())
	})

	It("round trips service error metadata through the generated HTTP client", func() {
		service := &fakeService{failure:coreerrors.NewErrorWithMetadata(personerror.PersonContactConflict,map[string]string{"field":"email"})}
		client := foundationclientimpl.NewPersonClient(fixtureMux(service))
		_, err := client.GetPersonByGID(context.Background(), foundationrequest.PersonByGID{GID:goodGID})
		var applicationError *coreerrors.Error
		Expect(errors.As(err,&applicationError)).To(BeTrue())
		Expect(applicationError.Definition().Code()).To(Equal(personerror.PersonContactConflict.Code()))
		Expect(applicationError.Definition().HTTPStatus()).To(Equal(http.StatusConflict))
		Expect(applicationError.Metadata()).To(Equal(map[string]any{"field":"email"}))
	})

	It("does not reveal unexpected service failures", func() {
		service := &fakeService{failure: errors.New("postgres secret error")}
		response := call(fixtureMux(service), http.MethodGet, "/foundation/protected/persons/"+goodGID, "")
		Expect(response.Code).To(Equal(http.StatusInternalServerError))
		Expect(response.Body.String()).To(ContainSubstring("FOUNDATION_PERSON_READ_FAILED"))
		Expect(response.Body.String()).NotTo(ContainSubstring("postgres secret error"))
	})
})
