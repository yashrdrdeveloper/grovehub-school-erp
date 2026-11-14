package routes_test

import (
    "net/http"
    "net/http/httptest"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"

    routes "grovehub/services/foundation/internal/app/routes"
    ids "grovehub/services/foundation/internal/tests/constants/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

// @Testing suite=FoundationRoutes
func Specs(){
    Describe("RegisterHTTP",func(){
        It("registers the annotated Person controller under /foundation/protected",func(){
            controller:=mocks.NewMockPersonController(gomock.NewController(GinkgoT()))
            controller.EXPECT().GetByGID(gomock.Any(),gomock.Any()).DoAndReturn(func(w http.ResponseWriter,r *http.Request){
                Expect(r.PathValue("gid")).To(Equal(ids.GID))
                w.WriteHeader(http.StatusNoContent)
            })
            mux:=http.NewServeMux()
            routes.RegisterHTTP(mux,controller)
            request:=httptest.NewRequest(http.MethodGet,"/foundation/protected/persons/"+ids.GID,nil)
            rec:=httptest.NewRecorder()
            mux.ServeHTTP(rec,request)
            Expect(rec.Code).To(Equal(http.StatusNoContent))
        })
        It("returns 404 for unregistered routes without invoking a controller",func(){
            controller:=mocks.NewMockPersonController(gomock.NewController(GinkgoT()))
            mux:=http.NewServeMux()
            routes.RegisterHTTP(mux,controller)
            request:=httptest.NewRequest(http.MethodGet,"/foundation/unrecognized",nil)
            rec:=httptest.NewRecorder()
            mux.ServeHTTP(rec,request)
            Expect(rec.Code).To(Equal(http.StatusNotFound))
        })
    })
}
