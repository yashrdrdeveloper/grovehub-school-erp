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
	identityerror "grovehub/services/identity/constants/error"
	controllerimpl "grovehub/services/identity/controller/impl"
	identityrequest "grovehub/services/identity/dto/request"
	identityresponse "grovehub/services/identity/dto/response"
	identityroutes "grovehub/services/identity/routes"
	validatorimpl "grovehub/services/identity/validators/impl"
)

type fakeAuthService struct {
	calls int
	last identityrequest.Register
	failure error
}

func (service *fakeAuthService) Register(_ context.Context, input identityrequest.RegisterRequest) (identityresponse.Register, error) {
	service.calls++
	service.last = input.RequestBody
	return identityresponse.Register{Email: input.RequestBody.Email}, service.failure
}

func serveRegister(service *fakeAuthService, body string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	identityroutes.RegisterHTTP(mux, controllerimpl.NewAuthController(service, validatorimpl.NewAuthValidator()))
	recorder := httptest.NewRecorder()
	mux.ServeHTTP(recorder, httptest.NewRequest(http.MethodPost, "/auth/register", strings.NewReader(body)))
	return recorder
}

const validRegistration = `{"firstName":"  Sima ","middleName":"  Devi ","lastName":" Bharti ","email":" SIMA@EXAMPLE.COM ","password":"securePassword123"}`

var _ = Describe("typed Identity Register pipeline", func() {
	It("normalizes a valid body and returns 201 with registration response", func() {
		service := &fakeAuthService{}
		response := serveRegister(service, validRegistration)
		Expect(response.Code).To(Equal(http.StatusCreated))
		Expect(service.calls).To(Equal(1))
		Expect(service.last.FirstName).To(Equal("Sima"))
		Expect(*service.last.MiddleName).To(Equal("Devi"))
		Expect(service.last.Email).To(Equal("sima@example.com"))
	})

	It("does not invoke the service for invalid DTO tags or strict JSON", func() {
		for _, body := range []string{
			`{"firstName":" ","middleName":"ok","lastName":"Bharti","email":"sima@example.com","password":"securePassword123"}`,
			`{"firstName":"Sima","middleName":"ok","lastName":"Bharti","email":"sima@example.com","password":"securePassword123","unknown":true}`,
		} {
			service := &fakeAuthService{}
			response := serveRegister(service, body)
			Expect(response.Code).To(Equal(http.StatusBadRequest))
			Expect(service.calls).To(BeZero())
		}
	})

	It("uses semantic errors and safe fallback without duplicating status maps", func() {
		service := &fakeAuthService{failure: coreerrors.NewError(identityerror.RegistrationAlreadyExists)}
		response := serveRegister(service, validRegistration)
		Expect(response.Code).To(Equal(http.StatusConflict))
		Expect(response.Body.String()).To(ContainSubstring("IDENTITY_REGISTRATION_ALREADY_EXISTS"))

		service.failure = errors.New("private failure")
		response = serveRegister(service, validRegistration)
		Expect(response.Code).To(Equal(http.StatusInternalServerError))
		Expect(response.Body.String()).To(ContainSubstring("IDENTITY_REGISTRATION_FAILED"))
		Expect(response.Body.String()).NotTo(ContainSubstring("private failure"))
	})
})
