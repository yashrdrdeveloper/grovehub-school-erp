package helpers_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    changed "grovehub/services/foundation/internal/app/events/producers/person/operations/changed"
    helpers "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/helpers"
)

// @Testing suite=FoundationPersonEventHelpers
func MetadataSpecs() {
    Describe("ChangedMetadata",func(){
        It("builds stable event transport metadata without PII",func(){
            metadata:=helpers.ChangedMetadata("event-42")
            Expect(metadata).To(Equal(map[string]string{
                "event-id":"event-42","event-type":changed.ChangedType,"source":changed.Source,
            }))
        })
    })
}
