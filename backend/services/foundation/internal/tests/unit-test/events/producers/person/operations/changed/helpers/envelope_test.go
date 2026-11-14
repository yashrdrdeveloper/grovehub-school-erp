package helpers_test

import (
    "context"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    coreevents "grovehub/libs/core/events"
    ops "grovehub/services/foundation/internal/app/constants/events"
    changed "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
    helpers "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/helpers"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
    ids "grovehub/services/foundation/internal/tests/constants/person"
)

// @Testing suite=FoundationPersonEventHelpers
func EnvelopeSpecs() {
    Describe("BuildChanged",func(){
        It("builds a typed v1 event envelope with context IDs",func(){
            next:=fixtures.Response()
            ctx:=coreevents.WithCorrelationID(context.Background(),"correlation-1")
            ctx=coreevents.WithCausationID(ctx,"command-1")
            e:=helpers.BuildChanged(ctx,"event-1",ops.Created,nil,&next)
            Expect(e.EventID).To(Equal("event-1"))
            Expect(e.EventType).To(Equal(changed.ChangedType))
            Expect(e.Version).To(Equal(changed.ChangedVersion))
            Expect(e.AggregateID).To(Equal(ids.GID))
            Expect(e.CorrelationID).To(Equal("correlation-1"))
            Expect(e.CausationID).To(Equal("command-1"))
            Expect(e.Data.Operation).To(Equal(ops.Created))
            Expect(e.Data.Old).To(BeNil())
            Expect(e.Data.New).To(Equal(&next))
            Expect(e.OccurredAt.IsZero()).To(BeFalse())
        })
        It("uses the old snapshot GID when the new one is nil",func(){
            previous:=fixtures.Response()
            e:=helpers.BuildChanged(context.Background(),"event-old",ops.Updated,&previous,nil)
            Expect(e.AggregateID).To(Equal(ids.GID))
        })
        It("handles an absent snapshot without inventing a person ID",func(){
            e:=helpers.BuildChanged(context.Background(),"event-none",ops.Created,nil,nil)
            Expect(e.AggregateID).To(BeEmpty())
        })
    })
}
