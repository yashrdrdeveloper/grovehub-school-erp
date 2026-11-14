package helpers_test

import (
    "context"
    "encoding/json"

    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"

    ops "grovehub/services/foundation/internal/app/constants/events"
    changed "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
    flows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
    helpers "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/helpers"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
)

// @Testing suite=FoundationPersonEventHelpers
func MessageSpecs() {
    Describe("BuildMessage",func(){
        It("serializes a typed changed event for the transactional outbox",func(){
            next:=fixtures.Response()
            message,err:=helpers.BuildMessage(context.Background(),flows.Change{Operation:ops.Created,New:&next})
            Expect(err).NotTo(HaveOccurred())
            Expect(message.EventID).NotTo(BeEmpty())
            Expect(message.Topic).To(Equal(changed.ChangedTopic))
            Expect(message.Metadata["event-id"]).To(Equal(message.EventID))
            var payload map[string]any
            Expect(json.Unmarshal(message.Payload,&payload)).To(Succeed())
            Expect(payload["eventType"]).To(Equal(changed.ChangedType))
            Expect(payload["aggregateId"]).To(Equal(next.GID))
        })
    })
}
