package impl_test

import (
    "context"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"
    ops "grovehub/services/foundation/internal/app/constants/events"
    changed "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
    strategies "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/strategies/impl"
    producer "grovehub/services/foundation/internal/app/events/producers/person/impl"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

// @Testing suite=FoundationPersonProducer
func Specs(){
    Describe("PersonProducer",func(){
        It("forwards a created snapshot to the outbox pipeline",func(){
            outbox:=mocks.NewMockOutboxStore(gomock.NewController(GinkgoT()))
            next:=fixtures.Response()
            outbox.EXPECT().Enqueue(gomock.Any(),gomock.Any()).Return(nil)
            p:=producer.New(changed.NewChangedPipeline(strategies.New(),outbox))
            Expect(p.Enqueue(context.Background(),ops.Created,nil,&next)).To(Succeed())
        })
        It("suppresses unchanged updates",func(){
            outbox:=mocks.NewMockOutboxStore(gomock.NewController(GinkgoT()))
            old:=fixtures.Response()
            p:=producer.New(changed.NewChangedPipeline(strategies.New(),outbox))
            Expect(p.Enqueue(context.Background(),ops.Updated,&old,&old)).To(Succeed())
        })
    })
}
