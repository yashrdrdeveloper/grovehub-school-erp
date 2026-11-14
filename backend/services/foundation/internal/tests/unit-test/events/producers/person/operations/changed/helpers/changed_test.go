package helpers_test

import (
    . "github.com/onsi/ginkgo/v2"
    . "github.com/onsi/gomega"
    eventops "grovehub/services/foundation/internal/app/constants/events"
    flows "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/flows"
    helpers "grovehub/services/foundation/internal/app/events/producers/person/operations/changed/helpers"
    fixtures "grovehub/services/foundation/internal/tests/helpers/person"
)

// @Testing suite=FoundationPersonEventHelpers
func ChangeSpecs() {
    Describe("HasChanged",func(){
        It("permits only a first creation with a nonnil next snapshot",func(){
            next:=fixtures.Response()
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Created,New:&next})).To(BeTrue())
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Created})).To(BeFalse())
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Created,Old:&next,New:&next})).To(BeFalse())
        })
        It("suppresses idempotent or incomplete updates",func(){
            previous:=fixtures.Response()
            current:=previous
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Updated,Old:&previous,New:&current})).To(BeFalse())
            current.Version++
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Updated,Old:&previous,New:&current})).To(BeTrue())
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Updated,New:&current})).To(BeFalse())
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Updated,Old:&previous})).To(BeFalse())
        })
        It("permits deactivation only when aggregate version changes",func(){
            previous:=fixtures.Response()
            current:=previous
            current.Version++
            current.Status="DEACTIVATED"
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Deactivated,Old:&previous,New:&current})).To(BeTrue())
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Deactivated,Old:&previous,New:&previous})).To(BeFalse())
            Expect(helpers.HasChanged(flows.Change{Operation:eventops.Operation("READ"),New:&current})).To(BeFalse())
        })
    })
}
