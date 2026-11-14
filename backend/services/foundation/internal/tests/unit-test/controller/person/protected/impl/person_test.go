package impl_test

import (
    "net/http"
    "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"

    httphelpers "grovehub/libs/core/tests/helpers/http"
    coresearch "grovehub/libs/core/search"
    controller "grovehub/services/foundation/internal/app/controller/person/protected"
    controllerimpl "grovehub/services/foundation/internal/app/controller/person/protected/impl"
    dto "grovehub/services/foundation/internal/app/dto/person"
    validatorimpl "grovehub/services/foundation/internal/app/validators/person/impl"
    ids "grovehub/services/foundation/internal/tests/constants/person"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    personmocks "grovehub/services/foundation/internal/tests/mocks/person"
)

var service *personmocks.MockPersonService
var handler controller.PersonController
var fixedGID string

// @BeforeAll
func setupSuite() { fixedGID=ids.GID }

// @BeforeEach
func setupSpec() {
    service=personmocks.NewMockPersonService(gomock.NewController(ginkgo.GinkgoT()))
    handler=controllerimpl.NewPersonController(service,validatorimpl.NewPersonValidator())
}

// @Testing suite=FoundationPersonController
func ControllerSpecs() {
    Describe("Person protected HTTP controller",func(){
        It("creates a person using a mocked domain service",func(){
            service.EXPECT().Create(gomock.Any(),gomock.Any()).Return(fixtures.Response(),true,nil)
            body:=[]byte(`{"firstName":"Asha","lastName":"Shah","contacts":{"emails":[{"value":"person@example.com","isPrimary":true}]}}`)
            request,rec:=httphelpers.Request(http.MethodPost,"/foundation/protected/persons",body)
            handler.Create(rec,request)
            Expect(rec.Code).To(Equal(http.StatusCreated))
            Expect(rec.Body.String()).To(ContainSubstring(fixedGID))
        })
        It("returns a mocked person by valid GID",func(){
            service.EXPECT().GetByGID(gomock.Any(),gomock.Any()).Return(fixtures.Response(),nil)
            request,rec:=httphelpers.Request(http.MethodGet,"/foundation/protected/persons/"+ids.GID,nil)
            request.SetPathValue("gid",ids.GID)
            handler.GetByGID(rec,request)
            Expect(rec.Code).To(Equal(http.StatusOK))
            Expect(rec.Body.String()).To(ContainSubstring(ids.GID))
        })
        It("updates a person through the same request pipeline",func(){
            service.EXPECT().Update(gomock.Any(),gomock.Any()).Return(fixtures.Response(),nil)
            request,rec:=httphelpers.Request(http.MethodPatch,"/foundation/protected/persons/"+ids.GID,[]byte(`{"firstName":"Asha"}`))
            request.SetPathValue("gid",ids.GID)
            handler.Update(rec,request)
            Expect(rec.Code).To(Equal(http.StatusOK))
        })
        It("searches people without any external database",func(){
            service.EXPECT().Search(gomock.Any(),gomock.Any()).Return(coresearch.Result[dto.Person]{Items:[]dto.Person{fixtures.Response()}},nil)
            request,rec:=httphelpers.Request(http.MethodPost,"/foundation/protected/persons/search",[]byte(`{}`))
            handler.Search(rec,request)
            Expect(rec.Code).To(Equal(http.StatusOK))
            Expect(rec.Body.String()).To(ContainSubstring(ids.GID))
        })
        It("deactivates a person with a mocked service",func(){
            service.EXPECT().Deactivate(gomock.Any(),gomock.Any()).Return(fixtures.Response(),nil)
            request,rec:=httphelpers.Request(http.MethodPatch,"/foundation/protected/persons/"+ids.GID+"/deactivate",nil)
            request.SetPathValue("gid",ids.GID)
            handler.Deactivate(rec,request)
            Expect(rec.Code).To(Equal(http.StatusOK))
        })
        It("rejects invalid route identifiers without invoking the service",func(){
            request,rec:=httphelpers.Request(http.MethodGet,"/foundation/protected/persons/bad",nil)
            request.SetPathValue("gid","bad")
            handler.GetByGID(rec,request)
            Expect(rec.Code).To(Equal(http.StatusBadRequest))
        })
        It("rejects malformed create JSON without invoking the service",func(){
            request,rec:=httphelpers.Request(http.MethodPost,"/foundation/protected/persons",[]byte(`{"firstName":`))
            handler.Create(rec,request)
            Expect(rec.Code).To(Equal(http.StatusBadRequest))
        })
    })
}
