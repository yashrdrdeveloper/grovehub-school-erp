package create_test

import (
    "context"
    "errors"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"
    pipeline "grovehub/libs/core/pipeline/service"
    dto "grovehub/services/foundation/internal/app/dto/person"
    create "grovehub/services/foundation/internal/app/services/person/operations/create"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

var strategy *mocks.MockCreateStrategy

// @BeforeEach
func prepare() { strategy=mocks.NewMockCreateStrategy(gomock.NewController(GinkgoT())) }

// @Testing suite=PersonCreateOperation
func Specs() {
    Describe("CreateOperation",func(){
        It("runs all four strategies in order, then returns the created Person",func(){
            gomock.InOrder(
                strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]()),
                strategy.EXPECT().Idempotency(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]()),
                strategy.EXPECT().Contacts(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]()),
                strategy.EXPECT().Persist(gomock.Any(),gomock.Any()).Return(pipeline.Return(fixtures.Response(),true)),
            )
            p,created,err:=create.New(strategy).Create(context.Background(),fixtures.CreateRequest())
            Expect(err).NotTo(HaveOccurred())
            Expect(created).To(BeTrue())
            Expect(p).To(Equal(fixtures.Response()))
        })
        It("short circuits an idempotent return before contact/persist operations",func(){
            gomock.InOrder(
                strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]()),
                strategy.EXPECT().Idempotency(gomock.Any(),gomock.Any()).Return(pipeline.Return(fixtures.Response(),false)),
            )
            _,created,err:=create.New(strategy).Create(context.Background(),fixtures.CreateRequest())
            Expect(err).NotTo(HaveOccurred())
            Expect(created).To(BeFalse())
        })
        It("short circuits the pipeline after a strategy error",func(){
            failure:=errors.New("invalid request")
            strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(pipeline.Fail[dto.Person](failure))
            _,_,err:=create.New(strategy).Create(context.Background(),fixtures.CreateRequest())
            Expect(err).To(MatchError(failure))
        })
        It("rejects a pipeline with no final result",func(){
            strategy.EXPECT().Prepare(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]())
            strategy.EXPECT().Idempotency(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]())
            strategy.EXPECT().Contacts(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]())
            strategy.EXPECT().Persist(gomock.Any(),gomock.Any()).Return(pipeline.Next[dto.Person]())
            _,_,err:=create.New(strategy).Create(context.Background(),fixtures.CreateRequest())
            Expect(err).To(MatchError(pipeline.ErrUnresolvedChain))
        })
    })
}
