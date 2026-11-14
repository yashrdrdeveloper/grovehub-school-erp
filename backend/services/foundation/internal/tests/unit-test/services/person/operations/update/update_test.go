package update_test

import (
    "context"
    "errors"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"
    pipeline "grovehub/libs/core/pipeline/service"
    dto "grovehub/services/foundation/internal/app/dto/person"
    update "grovehub/services/foundation/internal/app/services/person/operations/update"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

var strategy *mocks.MockUpdateStrategy

// @BeforeEach
func prepare(){ strategy=mocks.NewMockUpdateStrategy(gomock.NewController(GinkgoT())) }

// @Testing suite=PersonUpdateOperation
func Specs(){
    Describe("UpdateOperation",func(){
        It("loads, applies and persists updates in order",func(){
            gomock.InOrder(
                strategy.EXPECT().Load(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]()),
                strategy.EXPECT().Apply(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]()),
                strategy.EXPECT().Persist(gomock.Any(),gomock.Any()).Return(pipeline.Return(fixtures.Response(),false)),
            )
            p,err:=update.New(strategy).Update(context.Background(),fixtures.UpdateRequest())
            Expect(err).NotTo(HaveOccurred())
            Expect(p).To(Equal(fixtures.Response()))
        })
        It("short circuits a failed load",func(){
            failure:=errors.New("load failed")
            strategy.EXPECT().Load(gomock.Any(),gomock.Any()).Return(pipeline.Fail[dto.Person](failure))
            _,err:=update.New(strategy).Update(context.Background(),fixtures.UpdateRequest())
            Expect(err).To(MatchError(failure))
        })
    })
}
