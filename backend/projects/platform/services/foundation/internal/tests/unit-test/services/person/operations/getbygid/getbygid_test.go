package getbygid_test

import (
    "context"
    "errors"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"
    read "grovehub/services/foundation/internal/app/services/person/operations/getbygid"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

// @Testing suite=PersonReadOperation
func Specs() {
    Describe("GetByGIDOperation",func(){
        It("delegates a typed read request",func(){
            s:=mocks.NewMockGetByGIDStrategy(gomock.NewController(GinkgoT()))
            s.EXPECT().Execute(gomock.Any(),fixtures.GetRequest()).Return(fixtures.Response(),nil)
            result,err:=read.New(s).GetByGID(context.Background(),fixtures.GetRequest())
            Expect(err).NotTo(HaveOccurred())
            Expect(result).To(Equal(fixtures.Response()))
        })
        It("passes a read failure to the caller",func(){
            s:=mocks.NewMockGetByGIDStrategy(gomock.NewController(GinkgoT()))
            failure:=errors.New("no such person")
            s.EXPECT().Execute(gomock.Any(),fixtures.GetRequest()).Return(fixtures.Response(),failure)
            _,err:=read.New(s).GetByGID(context.Background(),fixtures.GetRequest())
            Expect(err).To(MatchError(failure))
        })
    })
}
