package search_test

import (
    "context"
    "errors"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"
    coresearch "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/search"
    corepipeline "github.com/yashrdrdeveloper/grovehub-school-erp/libs/backend/core/pipeline"
    dto "grovehub/services/foundation/internal/app/dto/person"
    search "grovehub/services/foundation/internal/app/services/person/operations/search"
    fixtures "grovehub/services/foundation/internal/resources/tests/helpers/person"
    mocks "grovehub/services/foundation/internal/resources/tests/mocks/person"
)

type Result = coresearch.Result[dto.Person]
var strategy *mocks.MockSearchStrategy

// @BeforeEach
func prepare(){ strategy=mocks.NewMockSearchStrategy(gomock.NewController(GinkgoT())) }

// @Testing suite=PersonSearchOperation
func Specs(){
    Describe("SearchOperation",func(){
        It("executes all search phases and returns hydrated people",func(){
            response:=Result{Items:[]dto.Person{fixtures.Response()}}
            gomock.InOrder(
                strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(corepipeline.Next[Result]()),
                strategy.EXPECT().Query(gomock.Any(),gomock.Any()).Return(corepipeline.Next[Result]()),
                strategy.EXPECT().Hydrate(gomock.Any(),gomock.Any()).Return(corepipeline.Return(response,false)),
            )
            actual,err:=search.New(strategy).Search(context.Background(),dto.SearchPersonsRequest{})
            Expect(err).NotTo(HaveOccurred())
            Expect(actual).To(Equal(response))
        })
        It("stops on a repository query failure",func(){
            failure:=errors.New("unavailable")
            strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(corepipeline.Next[Result]())
            strategy.EXPECT().Query(gomock.Any(),gomock.Any()).Return(corepipeline.Fail[Result](failure))
            _,err:=search.New(strategy).Search(context.Background(),dto.SearchPersonsRequest{})
            Expect(err).To(MatchError(failure))
        })
    })
}
