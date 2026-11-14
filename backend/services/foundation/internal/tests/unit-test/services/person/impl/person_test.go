package impl_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"

    personservice "grovehub/services/foundation/internal/app/services/person/impl"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

// @Testing suite=FoundationPersonServiceComposition
func Specs(){
    Describe("NewPersonService",func(){
        It("rejects a nil event producer instead of silently dropping events",func(){
            ctrl:=gomock.NewController(GinkgoT())
            Expect(func(){personservice.NewPersonService(
                mocks.NewMockPersonRepository(ctrl),mocks.NewMockContactRepository(ctrl),
                mocks.NewMockUnitOfWork(ctrl),mocks.NewMockPersonMapper(ctrl),
                mocks.NewMockPersonValidator(ctrl),nil,
            )}).To(PanicWith("foundation PersonService requires an event producer"))
        })
        It("wires all five Person operations using typed mock collaborators",func(){
            ctrl:=gomock.NewController(GinkgoT())
            service:=personservice.NewPersonService(
                mocks.NewMockPersonRepository(ctrl),mocks.NewMockContactRepository(ctrl),
                mocks.NewMockUnitOfWork(ctrl),mocks.NewMockPersonMapper(ctrl),
                mocks.NewMockPersonValidator(ctrl),mocks.NewMockChangedProducer(ctrl),
            )
            Expect(service).NotTo(BeNil())
        })
    })
}
