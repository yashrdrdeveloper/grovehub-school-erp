package search_test

import (
    "context"
    "errors"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"
    coresearch "grovehub/libs/core/search"
    pipeline "grovehub/libs/core/pipeline/service"
    dto "grovehub/services/foundation/internal/app/dto/person"
    search "grovehub/services/foundation/internal/app/services/person/operations/search"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
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
                strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(pipeline.Next[Result]()),
                strategy.EXPECT().Query(gomock.Any(),gomock.Any()).Return(pipeline.Next[Result]()),
                strategy.EXPECT().Hydrate(gomock.Any(),gomock.Any()).Return(pipeline.Return(response,false)),
            )
            actual,err:=search.New(strategy).Search(context.Background(),dto.SearchPersonsRequest{})
            Expect(err).NotTo(HaveOccurred())
            Expect(actual).To(Equal(response))
        })
        It("stops on a repository query failure",func(){
            failure:=errors.New("unavailable")
            strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(pipeline.Next[Result]())
            strategy.EXPECT().Query(gomock.Any(),gomock.Any()).Return(pipeline.Fail[Result](failure))
            _,err:=search.New(strategy).Search(context.Background(),dto.SearchPersonsRequest{})
            Expect(err).To(MatchError(failure))
        })
    })
}
