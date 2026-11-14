package impl_test

import (
    "context"
    "net/http"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"

    httpresponse "grovehub/libs/core/http/response"
    foundationclient "grovehub/services/foundation/client/person"
    clientimpl "grovehub/services/foundation/internal/app/client/person/impl"
    routes "grovehub/services/foundation/internal/app/routes"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    ids "grovehub/services/foundation/internal/tests/constants/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

// @Testing suite=FoundationInProcessClient
func Specs() {
    Describe("Person in-process client",func(){
        It("calls the generated GET route and decodes the typed response without TCP",func(){
            controller:=mocks.NewMockPersonController(gomock.NewController(GinkgoT()))
            controller.EXPECT().GetByGID(gomock.Any(),gomock.Any()).DoAndReturn(func(w http.ResponseWriter,r *http.Request){
                Expect(r.Method).To(Equal(http.MethodGet))
                Expect(r.PathValue("gid")).To(Equal(ids.GID))
                Expect(httpresponse.WriteSuccess(w,http.StatusOK,fixtures.Response())).To(Succeed())
            })
            mux:=http.NewServeMux()
            routes.RegisterHTTP(mux,controller)
            client:=clientimpl.NewPersonClient(mux)
            result,err:=client.GetPersonByGID(context.Background(),foundationclient.GetPersonByGIDRequest{
                Params:foundationclient.PersonByGID{GID:ids.GID},
            })
            Expect(err).NotTo(HaveOccurred())
            Expect(result.GID).To(Equal(ids.GID))
        })
        It("returns a typed response and a created flag from the generated POST route",func(){
            controller:=mocks.NewMockPersonController(gomock.NewController(GinkgoT()))
            controller.EXPECT().Create(gomock.Any(),gomock.Any()).DoAndReturn(func(w http.ResponseWriter,r *http.Request){
                Expect(httpresponse.WriteSuccess(w,http.StatusCreated,fixtures.Response())).To(Succeed())
            })
            mux:=http.NewServeMux()
            routes.RegisterHTTP(mux,controller)
            client:=clientimpl.NewPersonClient(mux)
            result,created,err:=client.CreatePerson(context.Background(),fixtures.CreateRequest())
            Expect(err).NotTo(HaveOccurred())
            Expect(created).To(BeTrue())
            Expect(result.GID).To(Equal(ids.GID))
        })
    })
}
