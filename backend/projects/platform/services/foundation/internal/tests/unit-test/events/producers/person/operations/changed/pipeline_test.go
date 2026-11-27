package changed_test

import (
    "context"
    "errors"
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    "go.uber.org/mock/gomock"

    ops "grovehub/services/foundation/internal/app/constants/events"
    changed "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
    flows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
    strategies "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/strategies/impl"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    mocks "grovehub/services/foundation/internal/tests/mocks/person"
)

var outbox *mocks.MockOutboxStore
var pipeline *changed.ChangedPipeline

// @BeforeEach
func prepare(){
    outbox=mocks.NewMockOutboxStore(gomock.NewController(GinkgoT()))
    pipeline=changed.NewChangedPipeline(strategies.New(),outbox)
}

// @Testing suite=FoundationPersonChangedPipeline
func Specs(){
    Describe("ChangedPipeline",func(){
        It("exposes all three required pipeline rules",func(){
            rules:=pipeline.Rules()
            Expect(rules.ShouldEmit).NotTo(BeNil())
            Expect(rules.Build).NotTo(BeNil())
            Expect(rules.Enqueue).NotTo(BeNil())
        })
        It("enqueues created person events into a mocked outbox",func(){
            next:=fixtures.Response()
            outbox.EXPECT().Enqueue(gomock.Any(),gomock.Any()).Return(nil)
            Expect(pipeline.Enqueue(context.Background(),flows.Change{Operation:ops.Created,New:&next})).To(Succeed())
        })
        It("suppresses no-op events and does not enqueue",func(){
            current:=fixtures.Response()
            change:=flows.Change{Operation:ops.Updated,Old:&current,New:&current}
            Expect(pipeline.Enqueue(context.Background(),change)).To(Succeed())
        })
        It("propagates an outbox enqueue failure",func(){
            next:=fixtures.Response()
            failure:=errors.New("simulated outbox failure")
            outbox.EXPECT().Enqueue(gomock.Any(),gomock.Any()).Return(failure)
            err:=pipeline.Enqueue(context.Background(),flows.Change{Operation:ops.Created,New:&next})
            Expect(err).To(MatchError(failure))
        })
    })
}
