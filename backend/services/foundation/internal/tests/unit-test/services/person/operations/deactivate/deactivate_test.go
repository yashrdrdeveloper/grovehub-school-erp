package deactivate_test

import (
    "context"
    "errors"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"
    pipeline "grovehub/libs/core/pipeline/service"
    dto "grovehub/services/foundation/internal/app/dto/person"
    deactivate "grovehub/services/foundation/internal/app/services/person/operations/deactivate"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

var strategy *mocks.MockDeactivateStrategy

// @BeforeEach
func prepare(){strategy=mocks.NewMockDeactivateStrategy(gomock.NewController(GinkgoT()))}

// @Testing suite=PersonDeactivateOperation
func Specs(){
    Describe("DeactivateOperation",func(){
        It("loads, prepares and persists deactivation",func(){
            gomock.InOrder(
                strategy.EXPECT().Load(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]()),
                strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]()),
                strategy.EXPECT().Persist(gomock.Any(),gomock.Any()).Return(pipeline.Return(fixtures.Response(),false)),
            )
            p,err:=deactivate.New(strategy).Deactivate(context.Background(),fixtures.DeactivateRequest())
            Expect(err).NotTo(HaveOccurred())
            Expect(p).To(Equal(fixtures.Response()))
        })
        It("short circuits after an unsuccessful prepare",func(){
            failure:=errors.New("not allowed")
            strategy.EXPECT().Load(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]())
            strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(pipeline.Fail[dto.Person](failure))
            _,err:=deactivate.New(strategy).Deactivate(context.Background(),fixtures.DeactivateRequest())
            Expect(err).To(MatchError(failure))
        })
    })
}
